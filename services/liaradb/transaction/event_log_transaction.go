package transaction

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type EventLogTransaction struct {
	events []eventItem
	el     *eventlog.EventLog
}

type eventItem struct {
	e    *entity.Event
	data []byte
}

func (t *EventLogTransaction) CanAppend(ctx context.Context, tid value.TenantID, e *entity.Event) error {
	tn := tablename.New(tid)
	k := key.NewKey2(e.AggregateID.Bytes(), e.Version.Value())
	return t.el.CanAppend(ctx, tn, e.PartitionID, k)
}

func (t *EventLogTransaction) Append(e *entity.Event, data []byte) {
	t.events = append(t.events, eventItem{
		e:    e,
		data: data,
	})
}

func (t *EventLogTransaction) Commit(ctx context.Context, tid value.TenantID) error {
	tn := tablename.New(tid)
	for _, item := range t.events {
		k := key.NewKey2(item.e.AggregateID.Bytes(), item.e.Version.Value())
		err := t.el.AppendEvent(ctx, tn, item.e.PartitionID, k, item.data)
		if err != nil {
			return err
		}
	}

	return nil
}
