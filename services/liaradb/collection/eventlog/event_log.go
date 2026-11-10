package eventlog

import (
	"context"
	"errors"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/fixed"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/scan"
)

type EventLog struct {
	fc *fixed.FixedCollection
	c  *btree.Cursor
}

func New(fc *fixed.FixedCollection, c *btree.Cursor) *EventLog {
	return &EventLog{
		fc: fc,
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

	return l.c.Insert(ctx, tx, tn.Index(2), key.NewKey2(e.GlobalVersion.Bytes(), int64(e.SubVersion.Value())), rid)
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

func (l *EventLog) GetAggregate(
	ctx context.Context,
	tx fixed.Transaction,
	tn tablename.TableName,
	id value.AggregateID,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		if err := tx.SLock(ctx, tn, id.String()); err != nil {
			yield(nil, err)
			return
		}

		for sr, err := range l.c.SearchRange(ctx, tn.Index(0), key.NewKey(id.Bytes()), 0, 0) {
			if err != nil {
				yield(nil, err)
				return
			}

			if err := tx.SLock(ctx, tn, sr.Key.String()); err != nil {
				yield(nil, err)
				return
			}

			s, err := l.fc.GetByRecordID(ctx, tn, sr.RecordID)
			if err != nil {
				yield(nil, err)
				return
			}

			e, err := l.readEvent(s)
			if !yield(e, err) || err != nil {
				return
			}
		}
	}
}

func (l *EventLog) EventsAfterGlobalVersion(
	ctx context.Context,
	tx fixed.Transaction,
	tn tablename.TableName,
	version value.GlobalVersion,
	sub value.SubVersion,
	lw value.GlobalVersion,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		for sr, err := range l.c.SearchRange(ctx, tn.Index(2), key.NewKey2(version.Bytes(), int64(sub.Value())), 0, 0) {
			if err != nil {
				yield(nil, err)
				return
			}

			if err := tx.SLock(ctx, tn, sr.Key.String()); err != nil {
				yield(nil, err)
				return
			}

			s, err := l.fc.GetByRecordID(ctx, tn, sr.RecordID)
			if err != nil {
				yield(nil, err)
				return
			}

			e, err := l.readEvent(s)
			if e.GlobalVersion.Value() > lw.Value() {
				return
			}

			if !yield(e, err) || err != nil {
				return
			}

			if e.GlobalVersion.Value() == lw.Value() {
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
