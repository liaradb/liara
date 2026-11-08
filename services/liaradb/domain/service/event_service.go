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
	"github.com/liaradb/liaradb/collection/tenant"
	"github.com/liaradb/liaradb/domain/command"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/query"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/filecache"
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
	tenant *tenant.Tenant,
	idempotency *idempotency.Idempotency,
	eventLog *eventlog.EventLog,
	keyValue *keyvalue.KeyValue,
	outbox *outbox.Outbox,
) *EventService {
	return &EventService{
		baseService: newBaseService(txManager, tenant, idempotency),
		txManager:   txManager,
		idempotency: idempotency,
		eventLog:    eventLog,
		keyValue:    keyValue,
		outbox:      outbox,
	}
}

func (es *EventService) InitGlobalVersions(ctx context.Context, fsys filecache.FileSystem, dir string, tIDs []value.TenantID) error {
	tIDMap := make(map[value.TenantID][]value.PartitionID)
	for _, tid := range tIDs {
		tf := tablename.NewTenantFolder(tid)
		pids, err := tf.Partitions(fsys, dir)
		if err != nil {
			return err
		}

		tIDMap[tid] = pids
	}

	return es.initGlobalVersions(ctx, tIDMap)
}

func (es *EventService) initGlobalVersions(ctx context.Context, tIDs map[value.TenantID][]value.PartitionID) error {
	for tid, pids := range tIDs {
		tf := tablename.NewTenantFolder(tid)

		var gv value.GlobalVersion
		for _, pid := range pids {
			gv1, err := es.eventLog.GetMaxGlobalVersion(ctx, tf.EventLog(pid))
			if err != nil {
				return err
			}

			if gv1.Value() > gv.Value() {
				gv = gv1
			}
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

	tx, err := es.initTransaction(ctx, cmd.TenantID)
	if err != nil {
		return err
	}

	now := time.Now()
	// TODO: PartitionID should be on the transaction, not just the Event
	return tx.Run(ctx, func() error {
		tn := tablename.NewTenantFolder(cmd.TenantID).EventLog(cmd.PartitionID)
		elt := newEventLogTransaction(es.eventLog, tn)
		if rqid, ok := cmd.Options.RequestID(); ok {
			// Verify idempotency
			if ok, err := es.testRequestID(ctx, tx, tn, rqid); err != nil {
				return err
			} else if !ok {
				return value.ErrRequestIDExists
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
			if err := es.insertRequestID(ctx, tx, tn, rqid, value.NewTime(now)); err != nil {
				return err
			}
		}

		return elt.commit(ctx, tx)
	})
}

func (es *EventService) TestIdempotency(
	ctx context.Context,
	qry query.TestIdempotency,
) (result bool, err error) {
	tx, err := es.initTransaction(ctx, qry.TenantID)
	if err != nil {
		return false, err
	}

	defer tx.Release()

	tn := tablename.NewTenantFolder(qry.TenantID).RequestLog()
	return es.testRequestID(ctx, tx, tn, qry.RequestID)
}

func (es *EventService) Get(
	ctx context.Context,
	tid value.TenantID,
	pid value.PartitionID,
	id value.AggregateID,
) iter.Seq2[*entity.Event, error] {
	tx, err := es.initTransaction(ctx, tid)
	if err != nil {
		return iterator.Error[*entity.Event](err)
	}

	defer tx.Release()

	return func(yield func(*entity.Event, error) bool) {
		tn := tablename.NewTenantFolder(tid).EventLog(pid)
		for e, err := range es.eventLog.GetAggregate(ctx, tx, tn, id) {
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
	tx, err := es.initTransaction(ctx, tid)
	if err != nil {
		return iterator.Error[*entity.Event](err)
	}

	defer tx.Release()

	return func(yield func(*entity.Event, error) bool) {
		tn := tablename.NewTenantFolder(tid).EventLog(pid)
		for e, err := range es.eventLog.GetAggregate(ctx, tx, tn, id) {
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
	pid value.PartitionID,
	limit value.Limit,
) iter.Seq2[*entity.Event, error] {
	if limit == 0 {
		return func(yield func(*entity.Event, error) bool) {}
	}

	return func(yield func(*entity.Event, error) bool) {
		tx, err := es.initTransaction(ctx, tid)
		if err != nil {
			yield(nil, err)
			return
		}

		count := 0
		if err := tx.Run(ctx, func() error {
			for e, err := range es.eventsAfterGlobalVersion(ctx, tx, tid, pid, version) {
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
	tid value.TenantID,
	pid value.PartitionID,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		tx, err := es.initTransaction(ctx, tid)
		if err != nil {
			yield(nil, err)
			return
		}

		defer tx.Release()

		tn := tablename.NewTenantFolder(tid).EventLog(pid)
		for e, err := range es.eventLog.EventsAfterGlobalVersion(ctx, tx, tn, value.NewGlobalVersion(0)) {
			if !yield(e, err) {
				return
			}
		}
	}
}

func (es *EventService) eventsAfterGlobalVersion(
	ctx context.Context,
	tx *transaction.Transaction,
	tid value.TenantID,
	pid value.PartitionID,
	version value.GlobalVersion,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		tn := tablename.NewTenantFolder(tid).EventLog(pid)
		for e, err := range es.eventLog.EventsAfterGlobalVersion(ctx, tx, tn, version) {
			if !yield(e, err) {
				return
			}
		}
	}
}

func (es *EventService) GetByOutbox(
	ctx context.Context,
	tid value.TenantID,
	oid value.OutboxID,
	limit value.Limit,
) iter.Seq2[*entity.Event, error] {
	tx, err := es.initTransaction(ctx, tid)
	if err != nil {
		return iterator.Error[*entity.Event](err)
	}

	o, err := es.outbox.Get(ctx, tx, tid, oid)
	if err != nil {
		return iterator.Error[*entity.Event](err)
	}

	return es.GetAfterGlobalVersion(ctx, tid, o.GlobalVersion(), o.PartitionID(), limit)
}

func (es *EventService) SetValue(
	ctx context.Context,
	tid value.TenantID,
	r *entity.Row,
	data []byte,
) error {
	tx, err := es.initTransaction(ctx, tid)
	if err != nil {
		return err
	}

	tn := tablename.NewTenantFolder(tid).GlobalTable("value")
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
