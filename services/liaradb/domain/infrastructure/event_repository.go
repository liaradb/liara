package infrastructure

import (
	"bytes"
	"context"
	"iter"
	"time"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction"
)

type EventRepository struct {
	txManager *transaction.Manager
	kv        *keyvalue.KeyValue
	eventLog  *eventlog.EventLog
	btree     *btree.Cursor
	fileName  link.FileName // TODO: Remove this
}

func NewEventRepository(
	txManager *transaction.Manager,
	kv *keyvalue.KeyValue,
	eventLog *eventlog.EventLog,
	btree *btree.Cursor,
	fn link.FileName,
) *EventRepository {
	return &EventRepository{
		txManager: txManager,
		kv:        kv,
		eventLog:  eventLog,
		btree:     btree,
		fileName:  fn,
	}
}

var _ service.EventRepository = (*EventRepository)(nil)

func (r *EventRepository) Append(
	ctx context.Context,
	tenantID value.TenantID,
	e entity.Event, // TODO: Should this be a pointer?
) error {
	tx := r.txManager.Next()

	buf := bytes.NewBuffer(nil)
	if err := e.Write(buf); err != nil {
		return err
	}

	if err := tx.Insert(ctx,
		action.ItemID(e.ID.String()),
		time.Now(),
		buf.Bytes(),
	); err != nil {
		return err
	}

	return tx.Commit(ctx, r.fileName, time.Now())
}

func (r *EventRepository) CreateIndex(context.Context, value.TenantID) error {
	panic("unimplemented")
}

func (r *EventRepository) CreateTable(context.Context, value.TenantID) error {
	panic("unimplemented")
}

func (r *EventRepository) DropTable(context.Context, value.TenantID) error {
	panic("unimplemented")
}

func (r *EventRepository) Get(
	ctx context.Context,
	tenantID value.TenantID,
	id value.AggregateID,
) iter.Seq2[entity.Event, error] { // TODO: Should this be a pointer?
	return func(yield func(entity.Event, error) bool) {
		for e, err := range r.eventLog.GetAggregate(ctx, r.fileName, id) {
			if err != nil {
				yield(entity.Event{}, err)
				return
			}

			if !yield(*e, nil) {
				return
			}
		}
	}
}

func (r *EventRepository) GetAfterGlobalVersion(context.Context, value.TenantID, value.GlobalVersion, value.PartitionRange, value.Limit) iter.Seq2[entity.Event, error] {
	panic("unimplemented")
}

func (r *EventRepository) GetByAggregateIDAndName(context.Context, value.TenantID, value.AggregateID, value.AggregateName) iter.Seq2[entity.Event, error] {
	panic("unimplemented")
}
