package eventlog

import (
	"context"
	"errors"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/heap/record"
	"github.com/liaradb/liaradb/collection/indexed"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/scan"
)

type EventLog struct {
	ic *indexed.IndexedCollection
}

func New(ic *indexed.IndexedCollection) *EventLog {
	return &EventLog{
		ic: ic,
	}
}

func (l *EventLog) Append(
	ctx context.Context,
	tx indexed.Transaction,
	tn tablename.TableName,
	e *entity.Event,
) error {
	rid, err := l.ic.Insert(ctx, tx, tn, key.NewKey2(e.AggregateID().Bytes(), e.Version().Value()), e)
	if err != nil {
		return err
	}

	c := btree.NewCursor(tx)
	if err := c.Insert(ctx, tx, tn.Index(1), key.NewKey(e.ID().Bytes()), rid); err != nil {
		return err
	}

	return c.Insert(ctx, tx, tn.Index(2), key.NewKey2(e.GlobalVersion().Bytes(), int64(e.SubVersion().Value())), rid)
}

func (l *EventLog) CanAppend(
	ctx context.Context,
	tx indexed.Transaction,
	tn tablename.TableName,
	k key.Key,
) error {
	c := btree.NewCursor(tx)
	fn := tn.Index(0)
	_, err := c.Search(ctx, fn, k)
	if err == nil {
		return btree.ErrExists
	}

	if errors.Is(err, btree.ErrNotFound) {
		return nil
	}

	return err
}

func (l *EventLog) Find(ctx context.Context, tx indexed.Transaction, tn tablename.TableName, id value.EventID) (*entity.Event, error) {
	c := btree.NewCursor(tx)
	rid, err := c.Search(ctx, tn.Index(1), key.NewKey(id.Bytes()))
	if err != nil {
		return nil, err
	}

	r, err := l.ic.GetByRecordID(ctx, tx, tn, rid)

	return l.readEvent(r)
}

func (l *EventLog) GetAggregate(
	ctx context.Context,
	tx indexed.Transaction,
	tn tablename.TableName,
	id value.AggregateID,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		if err := tx.SLock(ctx, tn, id.String()); err != nil {
			yield(nil, err)
			return
		}

		c := btree.NewCursor(tx)
		for sr, err := range c.SearchRange(ctx, tn.Index(0), key.NewKey(id.Bytes()), 0, 0) {
			if err != nil {
				yield(nil, err)
				return
			}

			if err := tx.SLock(ctx, tn, sr.Key.String()); err != nil {
				yield(nil, err)
				return
			}

			r, err := l.ic.GetByRecordID(ctx, tx, tn, sr.RecordID)
			if err != nil {
				yield(nil, err)
				return
			}

			if e, err := l.readEvent(r); !yield(e, err) || err != nil {
				return
			}
		}
	}
}

func (l *EventLog) EventsAfterGlobalVersion(
	ctx context.Context,
	tx indexed.Transaction,
	tn tablename.TableName,
	version value.GlobalVersion,
	sub value.SubVersion,
	lw value.GlobalVersion,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		c := btree.NewCursor(tx)
		for sr, err := range c.SearchRange(ctx, tn.Index(2), key.NewKey2(version.Bytes(), int64(sub.Value())), 0, 0) {
			if err != nil {
				yield(nil, err)
				return
			}

			if err := tx.SLock(ctx, tn, sr.Key.String()); err != nil {
				yield(nil, err)
				return
			}

			r, err := l.ic.GetByRecordID(ctx, tx, tn, sr.RecordID)
			if err != nil {
				yield(nil, err)
				return
			}

			e, err := l.readEvent(r)
			if e.GlobalVersion().Value() > lw.Value() {
				return
			}

			if !yield(e, err) || err != nil {
				return
			}

			if e.GlobalVersion().Value() == lw.Value() {
				return
			}
		}
	}
}

func (*EventLog) readEvent(r *record.Record) (*entity.Event, error) {
	defer r.Release()

	e := entity.Event{}
	return &e, e.Read(r)
}

func (l *EventLog) GetMaxGlobalVersion(
	ctx context.Context,
	tx indexed.Transaction,
	tn tablename.TableName,
) (value.GlobalVersion, error) {
	c := btree.NewCursor(tx)
	k, err := c.Maximum(ctx, tn.Index(2))
	if err != nil {
		if errors.Is(err, btree.ErrNotFound) {
			return value.NewGlobalVersion(0), nil
		}

		return value.GlobalVersion{}, err
	}

	v, _ := scan.Uint64(k.A.Value())
	return value.NewGlobalVersion(v), nil
}
