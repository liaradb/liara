package eventlog

import (
	"context"
	"errors"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	fixed "github.com/liaradb/liaradb/collection/fixedv2"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/transaction/log"
)

// TODO: Create latching
type EventLog struct {
	fc      *fixed.FixedCollection
	storage *storage.Storage
	cursor  *btree.Cursor
	l       *log.Log
}

func New(s *storage.Storage, c *btree.Cursor, l *log.Log) *EventLog {
	return &EventLog{
		fc:      fixed.New(s, c, l),
		storage: s,
		cursor:  c,
		l:       l,
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
	return l.fc.Insert(ctx, tn.EventLog(pid), tn.Index(0, pid), k, data)
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

			d, err := l.fc.GetItemByRecordLocator(ctx, fn, rl)
			if err != nil {
				yield(nil, err)
				return
			}

			var buf buffer.Buffer
			buf.Reset(d)

			var e entity.Event
			if err := e.Read(&buf); err != nil {
				yield(nil, err)
				return
			}

			if e.AggregateID != id || !yield(&e, nil) {
				return
			}
		}
	}
}

func (l *EventLog) Events(ctx context.Context, tn tablename.TableName, pid value.PartitionID) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		buf := buffer.NewFromSlice(nil)

		for i, err := range l.fc.List(ctx, tn.EventLog(pid), tn.Index(0, pid), pid) {
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

		for i, err := range l.fc.List(ctx, tn.EventLog(pid), tn.Index(0, pid), pid) {
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

// TODO: Iterate until highwater
func (l *EventLog) Iterate(ctx context.Context, tn tablename.TableName, pid value.PartitionID) iter.Seq2[[]byte, error] {
	return l.fc.List(ctx, tn.EventLog(pid), tn.Index(0, pid), pid)
}
