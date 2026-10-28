package transaction

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/util/set"
)

type eventLogTransaction struct {
	events []*entity.Event
	keys   set.Set[key.Key]
	el     *eventlog.EventLog
}

func (t *eventLogTransaction) append(ctx context.Context, tn tablename.TableName, e *entity.Event) error {
	k := key.NewKey2(e.AggregateID.Bytes(), e.Version.Value())

	// Verify this AggregateID and Version is unique in this transaction
	if t.keys.Includes(k) {
		return btree.ErrExists
	}

	t.keys.Add(k)

	// Verify this AggregateID and Version is unique in Index
	if err := t.el.CanAppend(ctx, tn, e.PartitionID, k); err != nil {
		return err
	}

	t.events = append(t.events, e)
	return nil
}

func (t *eventLogTransaction) commit(ctx context.Context, tx *Transaction) error {
	tn := tablename.New(tx.tid)
	for _, e := range t.events {
		if err := t.el.Append(ctx, tx, tn, e.PartitionID, e); err != nil {
			return err
		}
	}

	return nil
}
