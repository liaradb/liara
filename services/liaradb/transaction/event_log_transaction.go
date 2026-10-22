package transaction

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type EventLogTransaction struct {
	events []*entity.Event
	el     *eventlog.EventLog
}

func (t *EventLogTransaction) CanAppend(ctx context.Context, tid value.TenantID, e *entity.Event) error {
	tn := tablename.New(tid)
	k := key.NewKey2(e.AggregateID.Bytes(), e.Version.Value())
	return t.el.CanAppend(ctx, tn, e.PartitionID, k)
}

func (t *EventLogTransaction) Append(e *entity.Event) {
	t.events = append(t.events, e)
}

func (t *EventLogTransaction) Commit(ctx context.Context, l span.Log, tx *Transaction) error {
	tn := tablename.New(tx.tid)
	for _, e := range t.events {
		if err := t.el.Append(ctx, l, tn, e.PartitionID, e); err != nil {
			return err
		}
	}

	return nil
}
