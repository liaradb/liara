package eventlog

import (
	"context"
	"errors"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/fixed"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/storage/node"
)

// TODO: Create latching
type EventLog struct {
	fc      *fixed.FixedCollection
	storage *storage.Storage
	cursor  *btree.Cursor
}

func New(s *storage.Storage, c *btree.Cursor) *EventLog {
	return &EventLog{
		fc:      fixed.New(s, c),
		storage: s,
		cursor:  c,
	}
}

func (l *EventLog) Append(ctx context.Context, tn tablename.TableName, pid value.PartitionID, e *entity.Event) error {
	b := buffer.New(l.storage.BufferSize())
	if err := e.Write(b); err != nil {
		return err
	}

	k := key.NewKey2(e.AggregateID.Bytes(), e.Version.Value())
	return l.AppendEvent(ctx, tn, pid, k, b.Bytes()[:b.Cursor()])
}

func (l *EventLog) AppendEvent(ctx context.Context, tn tablename.TableName, pid value.PartitionID, k key.Key, data []byte) error {
	return l.fc.Set(ctx, tn.EventLog(pid), tn.Index(0, pid), k, data)
}

func (l *EventLog) CanAppend(ctx context.Context, tn tablename.TableName, pid value.PartitionID, k key.Key) error {
	fn := tn.Index(0, pid)
	_, err := l.cursor.Search(ctx, fn, k)
	if err == nil {
		return btree.ErrExists
	}

	if errors.Is(err, btree.ErrNotFound) {
		return nil
	}

	return err
}

func (l *EventLog) Find(ctx context.Context, tn tablename.TableName, pid value.PartitionID, id value.EventID) (*entity.Event, error) {
	for e, err := range l.Events(ctx, tn, pid) {
		if err != nil {
			return nil, err
		}

		if e.ID == id {
			return e, nil
		}
	}

	return nil, page.ErrNotFound
}

func (l *EventLog) GetAggregate(ctx context.Context, tn tablename.TableName, pid value.PartitionID, id value.AggregateID) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		fn := tn.EventLog(pid)
		for rl, err := range l.cursor.SearchRange(ctx, tn.Index(0, pid), key.NewKey(id.Bytes()), 0, 0) {
			if err != nil {
				yield(nil, err)
				return
			}

			e, err := l.getEventByRecordLocator(ctx, fn, rl)
			if err != nil {
				yield(nil, err)
				return
			}

			if e.AggregateID != id || !yield(e, nil) {
				return
			}
		}
	}
}

func (l *EventLog) getEventByRecordLocator(ctx context.Context, fn link.FileName, rl link.RecordLocator) (*entity.Event, error) {
	b, err := l.storage.Request(ctx, link.NewBlockID(fn, rl.Block()))
	if err != nil {
		return nil, err
	}

	defer b.Release()

	n := node.New(b)

	if !n.IsPage() {
		return nil, page.ErrNotPage
	}

	data, ok := n.Child(int16(rl.Position()))
	if !ok {
		return nil, btree.ErrNotFound
	}

	var buf buffer.Buffer
	buf.Reset(data)

	var e entity.Event
	if err := e.Read(&buf); err != nil {
		return nil, err
	}

	return &e, nil
}

func (l *EventLog) Events(ctx context.Context, tn tablename.TableName, pid value.PartitionID) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		buf := buffer.NewFromSlice(nil)

		for i, err := range l.items(ctx, tn, pid) {
			if err != nil {
				yield(nil, err)
				return
			}

			buf.Reset(i)

			var e entity.Event
			if err := e.Read(buf); err != nil {
				yield(nil, err)
				return
			}

			if !yield(&e, nil) {
				return
			}
		}
	}
}

func (l *EventLog) EventsAfterGlobalVersion(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	version value.GlobalVersion,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		buf := buffer.NewFromSlice(nil)

		for i, err := range l.items(ctx, tn, pid) {
			if err != nil {
				yield(nil, err)
				return
			}

			buf.Reset(i)

			var e entity.Event
			if err := e.Read(buf); err != nil {
				yield(nil, err)
				return
			}

			// TODO: Use Index to skip
			if e.GlobalVersion.Value() < version.Value() {
				continue
			}

			if !yield(&e, nil) {
				return
			}
		}
	}
}

func (l *EventLog) items(ctx context.Context, tn tablename.TableName, pid value.PartitionID) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		for n, err := range l.Iterate(ctx, tn, pid) {
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(n, nil) {
				return
			}
		}
	}
}

func (l *EventLog) Iterate(ctx context.Context, tn tablename.TableName, pid value.PartitionID) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		fn := tn.EventLog(pid)
		highBid, err := l.storage.Highwater(ctx, fn)
		if err != nil {
			yield(nil, err)
			return
		}

		bid := fn.BlockID(0)
		for bid.Position() <= highBid.Position() {
			children, p, err := l.handleIteration(ctx, bid)
			if err != nil {
				yield(nil, err)
				return
			}

			for c := range children {
				if !yield(c, nil) {
					return
				}
			}

			bid.SetPosition(p)
		}
	}
}

func (l *EventLog) handleIteration(ctx context.Context, bid link.BlockID) (iter.Seq[[]byte], link.FilePosition, error) {
	b, err := l.storage.Request(ctx, bid)
	if err != nil {
		return nil, 0, err
	}

	n := node.New(b)
	if !n.IsPage() {
		return nil, 0, page.ErrNotPage
	}

	return func(yield func([]byte) bool) {
		defer n.Release()
		for c := range n.Children() {
			if !yield(c) {
				return
			}
		}
	}, bid.Position() + 1, nil
}
