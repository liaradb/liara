package eventlog

import (
	"context"
	"errors"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/fixed"
	"github.com/liaradb/liaradb/collection/record"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/scan"
	"github.com/liaradb/liaradb/storage"
)

// TODO: Create latching
type EventLog struct {
	fc *fixed.FixedCollection
	c  *btree.Cursor
}

func New(s *storage.Storage, c *btree.Cursor) *EventLog {
	return &EventLog{
		fc: fixed.New(record.NewRecordCollection(s), c),
		c:  c,
	}
}

func (l *EventLog) Append(
	ctx context.Context,
	tx fixed.Transaction,
	tn tablename.TableName,
	e *entity.Event,
) error {
	rid, err := l.fc.Insert(ctx, tx, tn, key.NewKey2(e.AggregateID.Bytes(), e.Version.Value()), e)
	if err != nil {
		return err
	}

	if err := l.c.Insert(ctx, tx, tn.Index(1), key.NewKey(e.ID.Bytes()), rid); err != nil {
		return err
	}

	return l.c.Insert(ctx, tx, tn.Index(2), key.NewKey(e.GlobalVersion.Bytes()), rid)
}

func (l *EventLog) CanAppend(ctx context.Context, tn tablename.TableName, k key.Key) error {
	fn := tn.Index(0)
	_, err := l.c.Search(ctx, fn, k)
	if err == nil {
		return btree.ErrExists
	}

	if errors.Is(err, btree.ErrNotFound) {
		return nil
	}

	return err
}

func (l *EventLog) Find(ctx context.Context, tn tablename.TableName, id value.EventID) (*entity.Event, error) {
	rid, err := l.c.Search(ctx, tn.Index(1), key.NewKey(id.Bytes()))
	if err != nil {
		return nil, err
	}

	s, err := l.fc.GetByRecordID(ctx, tn, rid)

	return l.readEvent(s)
}

func (l *EventLog) GetAggregate(ctx context.Context, tn tablename.TableName, id value.AggregateID) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		for rid, err := range l.c.SearchRange(ctx, tn.Index(0), key.NewKey(id.Bytes()), 0, 0) {
			if err != nil {
				yield(nil, err)
				return
			}

			s, err := l.fc.GetByRecordID(ctx, tn, rid)
			if err != nil {
				yield(nil, err)
				return
			}

			e, err := l.readEvent(s)
			if err != nil {
				yield(nil, err)
				return
			}

			// TODO: Do we need to check this?
			if e.AggregateID != id || !yield(e, nil) {
				return
			}
		}
	}
}

func (*EventLog) readEvent(s *span.Span) (*entity.Event, error) {
	defer s.Release()

	t := entity.Event{}
	return &t, t.Read(s)
}

func (l *EventLog) Events(ctx context.Context, tx fixed.Transaction, tn tablename.TableName) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		for s, err := range l.fc.List(ctx, tx, tn) {
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(l.readEvent(s)) {
				return
			}
		}
	}
}

func (l *EventLog) EventsAfterGlobalVersion(
	ctx context.Context,
	tn tablename.TableName,
	version value.GlobalVersion,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		for rid, err := range l.c.SearchRange(ctx, tn.Index(2), key.NewKey(version.Bytes()), 0, 0) {
			if err != nil {
				yield(nil, err)
				return
			}

			s, err := l.fc.GetByRecordID(ctx, tn, rid)
			if err != nil {
				yield(nil, err)
				return
			}

			e, err := l.readEvent(s)
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(e, nil) {
				return
			}
		}
	}
}

// TODO: Iterate until highwater
func (l *EventLog) Iterate(ctx context.Context, tx fixed.Transaction, tn tablename.TableName) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		for s, err := range l.fc.List(ctx, tx, tn) {
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(s.BytesAndRelease()) {
				return
			}
		}
	}
}

func (l *EventLog) GetMaxGlobalVersion(ctx context.Context, tn tablename.TableName) (value.GlobalVersion, error) {
	k, err := l.c.Maximum(ctx, tn.Index(2))
	if err != nil {
		if errors.Is(err, btree.ErrNotFound) {
			return value.NewGlobalVersion(0), nil
		}

		return value.GlobalVersion{}, err
	}

	v, _ := scan.Uint64([]byte(k.A.String()))
	return value.NewGlobalVersion(v), nil
}
