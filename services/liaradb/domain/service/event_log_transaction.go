package service

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/transaction"
	"github.com/liaradb/liaradb/util/set"
)

type eventLogTransaction struct {
	el     *eventlog.EventLog
	tn     tablename.TableName
	events []*entity.Event
	keys   set.Set[key.Key]
}

func newEventLogTransaction(el *eventlog.EventLog, tn tablename.TableName) eventLogTransaction {
	return eventLogTransaction{
		el:   el,
		tn:   tn,
		keys: set.Set[key.Key]{},
	}
}

func (t *eventLogTransaction) append(ctx context.Context, e *entity.Event) error {
	k := key.NewKey2(e.AggregateID.Bytes(), e.Version.Value())

	// Verify this AggregateID and Version is unique in this transaction
	if t.keys.Includes(k) {
		return btree.ErrExists
	}

	t.keys.Add(k)

	// Verify this AggregateID and Version is unique in Index
	if err := t.el.CanAppend(ctx, t.tn, e.PartitionID, k); err != nil {
		return err
	}

	t.events = append(t.events, e)
	return nil
}

func (t *eventLogTransaction) commit(ctx context.Context, tx *transaction.Transaction) error {
	for _, e := range t.events {
		if err := t.el.Append(ctx, tx, t.tn, e.PartitionID, e); err != nil {
			return err
		}
	}

	return nil
}
