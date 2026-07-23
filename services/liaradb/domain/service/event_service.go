package service

import (
	"bytes"
	"context"
	"iter"
	"time"

	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/command"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/query"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/transaction"
	"github.com/liaradb/liaradb/util/iterator"
)

type EventService struct {
	txManager *transaction.Manager
}

func NewEventService(
	txManager *transaction.Manager,
) *EventService {
	return &EventService{
		txManager: txManager,
	}
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
	return transaction.Run(ctx, tx, now, func() error {
		tn := tablename.New(cmd.TenantID)
		if rqid, ok := cmd.Options.RequestID(); ok {
			// Verify idempotency
			// TODO: What should this return if requestID is present?
			if ok, err := tx.TestRequestID(ctx, tn, rqid); err != nil || !ok {
				return err
			}
		}

		buf := bytes.NewBuffer(nil)

		for _, em := range cmd.Events {
			e, err := em.ToEvent(cmd.PartitionID, cmd.Options)
			if err != nil {
				return err
			}

			if err := e.Write(buf); err != nil {
				return err
			}

			if err := tx.Insert(ctx, tn, now, &e, buf.Bytes()); err != nil {
				return err
			}
		}

		if rqid, ok := cmd.Options.RequestID(); ok {
			// TODO: Do we want to store this if the transaction doesn't complete?
			return tx.InsertRequestID(ctx, tn, rqid, value.NewTime(now))
		}

		return nil
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

	now := time.Now()
	return transaction.RunResult(ctx, tx, now, func() (bool, error) {
		tn := tablename.New(qry.TenantID)
		return tx.TestRequestID(ctx, tn, qry.RequestID)
	})
}

func (es *EventService) Get(
	ctx context.Context,
	tid value.TenantID,
	partitionID value.PartitionID,
	id value.AggregateID,
) iter.Seq2[*entity.Event, error] {
	tn := tablename.New(tid)
	tx, err := es.txManager.Next(ctx, tid)
	if err != nil {
		return iterator.Error[*entity.Event](err)
	}

	return tx.GetAggregate(ctx, tn, partitionID, id)
}

func (es *EventService) GetByAggregateIDAndName(
	ctx context.Context,
	tid value.TenantID,
	partitionID value.PartitionID,
	id value.AggregateID,
	name value.AggregateName,
) iter.Seq2[*entity.Event, error] {
	tn := tablename.New(tid)
	tx, err := es.txManager.Next(ctx, tid)
	if err != nil {
		return iterator.Error[*entity.Event](err)
	}

	return tx.GetAggregateByIDAndName(ctx, tn, partitionID, id, name)
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
		now := time.Now()
		// TODO: How do we handle a range?
		count := 0
		if err := transaction.Run(ctx, tx, now, func() error {
			tn := tablename.New(tid)
			for e, err := range tx.EventsAfterGlobalVersion(ctx, tn, partitionRange.Low(), version) {
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

func (es *EventService) GetByOutbox(
	ctx context.Context,
	tid value.TenantID,
	outboxID value.OutboxID,
	limit value.Limit,
) iter.Seq2[*entity.Event, error] {
	tx, err := es.txManager.Next(ctx, tid)
	if err != nil {
		return iterator.Error[*entity.Event](err)
	}

	tn := tablename.New(tid)
	o, err := tx.GetOutbox(ctx, tn, value.NewPartitionID(0), outboxID)
	if err != nil {
		return iterator.Error[*entity.Event](err)
	}

	return es.GetAfterGlobalVersion(ctx, tid, o.GlobalVersion(), o.PartitionRange(), limit)
}

func (es *EventService) CreateOutbox(
	ctx context.Context,
	tid value.TenantID,
	partitionRange value.PartitionRange,
) (value.OutboxID, error) {
	tx, err := es.txManager.Next(ctx, tid)
	if err != nil {
		return value.OutboxID{}, err
	}

	now := time.Now()
	// TODO: How do we handle a range?
	return transaction.RunResult(ctx, tx, now, func() (value.OutboxID, error) {
		tn := tablename.New(tid)
		oid := value.NewOutboxID()
		outbox := entity.NewOutbox(oid, partitionRange)
		if err := tx.InsertOutbox(ctx, tn, partitionRange.Low(), now, oid, outbox); err != nil {
			return value.OutboxID{}, err
		}
		return oid, nil
	})
}

func (es *EventService) GetOutbox(
	ctx context.Context,
	tid value.TenantID,
	outboxID value.OutboxID,
) (*entity.Outbox, error) {
	tx, err := es.txManager.Next(ctx, tid)
	if err != nil {
		return nil, err
	}

	return transaction.RunResult(ctx, tx, time.Now(), func() (*entity.Outbox, error) {
		tn := tablename.New(tid)
		return tx.GetOutbox(ctx, tn, value.NewPartitionID(0), outboxID)
	})
}

func (es *EventService) UpdateOutboxPosition(
	ctx context.Context,
	tid value.TenantID,
	outboxID value.OutboxID,
	globalVersion value.GlobalVersion,
) error {
	tx, err := es.txManager.Next(ctx, tid)
	if err != nil {
		return err
	}

	now := time.Now()
	return transaction.Run(ctx, tx, now, func() error {
		tn := tablename.New(tid)
		return tx.UpdateOutbox(ctx, tn, value.NewPartitionID(0), now, outboxID, globalVersion)
	})
}

func (es *EventService) ListOutboxes(
	ctx context.Context,
	tid value.TenantID,
) iter.Seq2[*entity.Outbox, error] {
	tx, err := es.txManager.Next(ctx, tid)
	if err != nil {
		return iterator.Error[*entity.Outbox](err)
	}

	tn := tablename.New(tid)
	return tx.ListOutboxes(ctx, tn, value.NewPartitionID(0))
}
