package service

import (
	"context"
	"iter"
	"time"

	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/idempotency"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/outbox"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/command"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/query"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/transaction"
	"github.com/liaradb/liaradb/util/iterator"
)

type EventService struct {
	baseService
	txManager   *transaction.Manager
	idempotency *idempotency.Idempotency
	eventLog    *eventlog.EventLog
	keyValue    *keyvalue.KeyValue
	outbox      *outbox.Outbox
	gvg         GlobalVersionGenerator
}

func NewEventService(
	txManager *transaction.Manager,
	idempotency *idempotency.Idempotency,
	eventLog *eventlog.EventLog,
	keyValue *keyvalue.KeyValue,
	outbox *outbox.Outbox,
) *EventService {
	return &EventService{
		baseService: newBaseService(txManager, idempotency),
		txManager:   txManager,
		idempotency: idempotency,
		eventLog:    eventLog,
		keyValue:    keyValue,
		outbox:      outbox,
	}
}

func (es *EventService) InitGlobalVersions(ctx context.Context, tIDs []value.TenantID) error {
	for _, tid := range tIDs {
		// TODO: Restore all partitions from disk
		gv, err := es.eventLog.GetMaxGlobalVersion(ctx, tablename.New(tid), value.NewPartitionID(0))
		if err != nil {
			return err
		}

		es.gvg.Init(tid, gv)
	}

	return nil
}

func (es *EventService) Append(
	ctx context.Context,
	cmd command.AppendEvent,
) error {
	if err := cmd.Valid(); err != nil {
		return err
	}

	tx, err := es.txManager.Next(ctx, cmd.TenantID)
	if err != nil {
		return err
	}

	now := time.Now()
	// TODO: PartitionID should be on the transaction, not just the Event
	return tx.Run(ctx, func() error {
		tn := tablename.New(cmd.TenantID)
		elt := newEventLogTransaction(es.eventLog, tn)
		if rqid, ok := cmd.Options.RequestID(); ok {
			// Verify idempotency
			// TODO: What should this return if requestID is present?
			if ok, err := es.testRequestID(ctx, tn, rqid); err != nil || !ok {
				return err
			}
		}

		for _, em := range cmd.Events {
			e, err := em.ToEvent(cmd.PartitionID, cmd.Options, es.gvg.Next(cmd.TenantID))
			if err != nil {
				return err
			}

			// TODO: What do we lock?
			if err := tx.XLock(ctx, tn, e.AggregateID.String()); err != nil {
				return err
			}

			if err := elt.append(ctx, &e); err != nil {
				return err
			}
		}

		if rqid, ok := cmd.Options.RequestID(); ok {
			// TODO: Do we want to store this if the transaction doesn't complete?
			return es.insertRequestID(ctx, tx, tn, rqid, value.NewTime(now))
		}

		return elt.commit(ctx, tx)
	})
}

func (es *EventService) TestIdempotency(
	ctx context.Context,
	qry query.TestIdempotency,
) (result bool, err error) {
	tx, err := es.txManager.Next(ctx, qry.TenantID)
	if err != nil {
		return false, err
	}

	// TODO: Should this log?
	return tx.RunResult(ctx, func() (bool, error) {
		tn := tablename.New(qry.TenantID)
		return es.testRequestID(ctx, tn, qry.RequestID)
	})
}

func (es *EventService) Get(
	ctx context.Context,
	tid value.TenantID,
	pid value.PartitionID,
	id value.AggregateID,
) iter.Seq2[*entity.Event, error] {
	tn := tablename.New(tid)
	tx, err := es.txManager.Next(ctx, tid)
	if err != nil {
		return iterator.Error[*entity.Event](err)
	}

	defer tx.Release()

	if err := tx.SLock(ctx, tn, id.String()); err != nil {
		return iterator.Error[*entity.Event](err)
	}

	return func(yield func(*entity.Event, error) bool) {
		for e, err := range es.eventLog.GetAggregate(ctx, tn, pid, id) {
			if !yield(e, err) {
				return
			}
		}
	}
}

func (es *EventService) GetByAggregateIDAndName(
	ctx context.Context,
	tid value.TenantID,
	pid value.PartitionID,
	id value.AggregateID,
	name value.AggregateName,
) iter.Seq2[*entity.Event, error] {
	tn := tablename.New(tid)
	tx, err := es.txManager.Next(ctx, tid)
	if err != nil {
		return iterator.Error[*entity.Event](err)
	}

	defer tx.Release()

	if err := tx.SLock(ctx, tn, id.String()); err != nil {
		return iterator.Error[*entity.Event](err)
	}

	return func(yield func(*entity.Event, error) bool) {
		for e, err := range es.eventLog.GetAggregate(ctx, tn, pid, id) {
			if err != nil {
				yield(nil, err)
				return
			}

			if e.AggregateName == name && !yield(e, nil) {
				return
			}
		}
	}
}

func (es *EventService) GetAfterGlobalVersion(
	ctx context.Context,
	tid value.TenantID,
	version value.GlobalVersion,
	partitionRange value.PartitionRange,
	limit value.Limit,
) iter.Seq2[*entity.Event, error] {
	if limit == 0 {
		return func(yield func(*entity.Event, error) bool) {}
	}

	return func(yield func(*entity.Event, error) bool) {
		tx, err := es.txManager.Next(ctx, tid)
		if err != nil {
			yield(nil, err)
			return
		}

		// TODO: Should this log?
		// TODO: How do we handle a range?
		count := 0
		if err := tx.Run(ctx, func() error {
			tn := tablename.New(tid)
			for e, err := range es.EventsAfterGlobalVersion(ctx, tn, partitionRange.Low(), version) {
				if err != nil {
					yield(nil, err)
					return err
				}

				count++
				if !yield(e, nil) || count >= limit.Value() {
					return nil
				}
			}

			return nil
		}); err != nil {
			yield(nil, err)
		}
	}
}

func (es *EventService) Events(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		// Should we lock?
		// if err := t.SLock(ctx, tn, id.String()); err != nil {
		// 	yield(nil, err)
		// 	return
		// }

		// defer t.release()

		for e, err := range es.eventLog.Events(ctx, tn, pid) {
			if !yield(e, err) {
				return
			}
		}
	}
}

func (es *EventService) EventsAfterGlobalVersion(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	version value.GlobalVersion,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		// Should we lock?
		// if err := t.SLock(ctx, tn, id.String()); err != nil {
		// 	yield(nil, err)
		// 	return
		// }

		// defer t.release()

		for e, err := range es.eventLog.EventsAfterGlobalVersion(ctx, tn, pid, version) {
			if !yield(e, err) {
				return
			}
		}
	}
}

func (es *EventService) getOutbox(
	ctx context.Context,
	tid value.TenantID,
	oid value.OutboxID,
) (*entity.Outbox, error) {
	tx, err := es.txManager.Next(ctx, tid)
	if err != nil {
		return nil, err
	}

	defer tx.Release()

	tn := tablename.New(tid)
	if err := tx.SLock(ctx, tn, oid.String()); err != nil {
		return nil, err
	}

	// TODO: What should PartitionID be?
	return es.outbox.Get(ctx, tn, value.NewPartitionID(0), oid)
}

func (es *EventService) GetByOutbox(
	ctx context.Context,
	tid value.TenantID,
	outboxID value.OutboxID,
	limit value.Limit,
) iter.Seq2[*entity.Event, error] {
	// tx, err := es.txManager.Next(ctx, tid)
	// if err != nil {
	// 	return iterator.Error[*entity.Event](err)
	// }

	o, err := es.getOutbox(ctx, tid, outboxID)
	if err != nil {
		return iterator.Error[*entity.Event](err)
	}

	return es.GetAfterGlobalVersion(ctx, tid, o.GlobalVersion(), o.PartitionRange(), limit)
}

func (es *EventService) SetValue(
	ctx context.Context,
	tn tablename.TableName,
	r *entity.Row,
	data []byte,
) error {
	tx, err := es.txManager.Next(ctx, tn.TenantID())
	if err != nil {
		return err
	}

	if err := tx.XLock(ctx, tn, r.ID().String()); err != nil {
		return err
	}

	return tx.Run(ctx, func() error {
		kvt := newKeyValueTransaction(es.keyValue)

		if err := kvt.setValue(ctx, tn, r, data); err != nil {
			return err
		}

		return kvt.commit(ctx)
	})
}
