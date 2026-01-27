package service

import (
	"bytes"
	"context"
	"iter"
	"time"

	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
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

type AppendOptions struct {
	requestID     *value.RequestID    // The ID of the Request, for idempotency
	correlationID value.CorrelationID // The ID of the entire Command and Event chain
	userID        value.UserID        // The ID of the User issuing the Command
	time          time.Time           // The Time this Event was created
}

func NewAppendOptions(
	requestID *value.RequestID, // The ID of the Request, for idempotency
	correlationID value.CorrelationID, // The ID of the entire Command and Event chain
	userID value.UserID, // The ID of the User issuing the Command
	time time.Time, // The Time this Event was created
) AppendOptions {
	return AppendOptions{
		requestID:     requestID,
		correlationID: correlationID,
		userID:        userID,
		time:          time,
	}
}

func (ao *AppendOptions) RequestID() (value.RequestID, bool) {
	if ao.requestID == nil {
		return value.NewRequestID(), false
	}

	return *ao.requestID, true
}

func (ao *AppendOptions) toMetadata() entity.Metadata {
	return entity.Metadata{
		UserID:        ao.userID,
		CorrelationID: ao.correlationID,
		Time:          value.NewTime(ao.time),
	}
}

type AppendEvent struct {
	ID            string              // The ID of the Event, used for de-duplication
	AggregateName value.AggregateName // The Name of the Aggregate
	AggregateID   value.AggregateID   // The ID of the Aggregate to which this Event applies
	Version       value.Version       // The Version of the Aggregate
	Name          value.EventName     // The Name of the Event
	Schema        value.Schema        // The Schema for the internal data
	Data          []byte              // The internal data of the Event
}

func (ae *AppendEvent) Valid() error {
	if ae.Version.Value() < 1 {
		return value.ErrAggregateVersionInvalid
	}

	return nil
}

func (ae *AppendEvent) toEvent(pid value.PartitionID, options AppendOptions) (entity.Event, error) {
	var id value.EventID
	if ae.ID == "" {
		id = value.NewEventID()
	} else {
		var err error
		id, err = value.NewEventIDFromString(ae.ID)
		if err != nil {
			return entity.Event{}, err
		}
	}

	return entity.Event{
		GlobalVersion: value.NewGlobalVersion(0),
		ID:            id,
		AggregateName: ae.AggregateName,
		AggregateID:   ae.AggregateID,
		Version:       ae.Version,
		PartitionID:   pid,
		Name:          ae.Name,
		Schema:        ae.Schema,
		Metadata:      options.toMetadata(),
		Data:          value.NewData(ae.Data),
	}, nil
}

func (es *EventService) Append(
	ctx context.Context,
	tenantID value.TenantID,
	options AppendOptions,
	pid value.PartitionID,
	e ...AppendEvent,
) error {
	if len(e) == 0 {
		return nil
	}

	for _, em := range e {
		if err := em.Valid(); err != nil {
			return err
		}
	}

	return es.append(ctx, tenantID, options, pid, e...)
}

func (es *EventService) append(
	ctx context.Context,
	tenantID value.TenantID,
	options AppendOptions,
	pid value.PartitionID,
	evs ...AppendEvent,
) error {
	if options.time.IsZero() {
		options.time = time.Now()
	}

	tx := es.txManager.Next()
	tn := tablename.New(tenantID)
	now := time.Now()
	// TODO: PartitionID should be on the transaction, not just the Event
	return tx.Run(ctx, tn, pid, now, func() error {
		if rqid, ok := options.RequestID(); ok {
			// Verify idempotency
			// TODO: What should this return if requestID is present?
			if ok, err := tx.TestIdempotency(ctx, tn, rqid); err != nil || !ok {
				return err
			}
		}

		buf := bytes.NewBuffer(nil)

		for _, em := range evs {
			e, err := em.toEvent(pid, options)
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

		if rqid, ok := options.RequestID(); ok {
			// TODO: Do we want to store this if the transaction doesn't complete?
			return tx.InsertRequestID(ctx, tn, rqid, now)
		}

		return nil
	})
}

func (es *EventService) TestIdempotency(
	ctx context.Context,
	tenantID value.TenantID,
	id value.RequestID,
) (result bool, err error) {
	tx := es.txManager.Next()
	tn := tablename.New(tenantID)
	now := time.Now()
	err = tx.Run(ctx, tn, value.PartitionID{}, now, func() error {
		result, err = tx.TestIdempotency(ctx, tn, id)
		return err
	})
	return
}

func (es *EventService) Get(
	ctx context.Context,
	tenantID value.TenantID,
	partitionID value.PartitionID,
	id value.AggregateID,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		tn := tablename.New(tenantID)
		tx := es.txManager.Next()
		for e, err := range tx.GetAggregate(ctx, tn, partitionID, id) {
			if err != nil {
				yield(nil, err)
				return
			}
			if !yield(e, nil) {
				return
			}
		}
	}
}

func (es *EventService) GetByAggregateIDAndName(
	ctx context.Context,
	tenantID value.TenantID,
	partitionID value.PartitionID,
	id value.AggregateID,
	name value.AggregateName,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		tn := tablename.New(tenantID)
		tx := es.txManager.Next()
		for e, err := range tx.GetAggregate(ctx, tn, partitionID, id) {
			if err != nil {
				yield(nil, err)
				return
			}

			// TODO: Move this to another layer
			if e.AggregateName == name && !yield(e, nil) {
				return
			}
		}
	}
}

func (es *EventService) GetAfterGlobalVersion(
	ctx context.Context,
	tenantID value.TenantID,
	version value.GlobalVersion,
	partitionRange value.PartitionRange,
	limit value.Limit,
) iter.Seq2[*entity.Event, error] {
	if limit == 0 {
		return func(yield func(*entity.Event, error) bool) {}
	}

	return func(yield func(*entity.Event, error) bool) {
		tn := tablename.New(tenantID)
		tx := es.txManager.Next()
		now := time.Now()
		// TODO: How do we handle a range?
		count := 0
		err := tx.Run(ctx, tn, partitionRange.Low(), now, func() error {
			for e, err := range tx.Events(ctx, tn, partitionRange.Low()) {
				if err != nil {
					yield(nil, err)
					return err
				}

				// TODO: Use Index to skip
				if e.GlobalVersion.Value() < version.Value() {
					continue
				}

				count++
				if !yield(e, nil) || count >= limit.Value() {
					return nil
				}
			}

			return nil
		})
		if err != nil {
			yield(nil, err)
		}
	}
}

func (es *EventService) GetByOutbox(
	ctx context.Context,
	tenantID value.TenantID,
	outboxID value.OutboxID,
	limit value.Limit,
) iter.Seq2[*entity.Event, error] {
	tn := tablename.New(tenantID)
	tx := es.txManager.Next()
	o, err := tx.GetOutbox(ctx, tn, outboxID)
	if err != nil {
		return iterator.Error[*entity.Event](err)
	}

	return es.GetAfterGlobalVersion(ctx, tenantID, o.GlobalVersion(), o.PartitionRange(), limit)
}

func (es *EventService) CreateOutbox(
	ctx context.Context,
	tenantID value.TenantID,
	outboxID value.OutboxID,
	partitionRange value.PartitionRange,
) (value.OutboxID, error) {
	tn := tablename.New(tenantID)
	tx := es.txManager.Next()
	now := time.Now()
	// TODO: How do we handle a range?
	err := tx.Run(ctx, tn, partitionRange.Low(), now, func() error {
		outbox := entity.NewOutbox(outboxID, partitionRange)
		return tx.InsertOutbox(ctx, tn, now, outboxID, outbox)
	})
	return outboxID, err
}

func (es *EventService) GetOutbox(
	ctx context.Context,
	tenantID value.TenantID,
	partitionID value.PartitionID,
	outboxID value.OutboxID,
) (*entity.Outbox, error) {
	tn := tablename.New(tenantID)
	tx := es.txManager.Next()
	var e *entity.Outbox
	err := tx.Run(ctx, tn, partitionID, time.Now(), func() error {
		var err error
		e, err = tx.GetOutbox(ctx, tn, outboxID)
		return err
	})
	return e, err
}

func (es *EventService) UpdateOutboxPosition(
	ctx context.Context,
	tenantID value.TenantID,
	partitionID value.PartitionID,
	outboxID value.OutboxID,
	globalVersion value.GlobalVersion,
) error {
	tn := tablename.New(tenantID)
	tx := es.txManager.Next()
	now := time.Now()
	return tx.Run(ctx, tn, partitionID, now, func() error {
		return tx.UpdateOutbox(ctx, tn, now, outboxID, globalVersion)
	})
}

func (es *EventService) ListOutboxes(
	ctx context.Context,
	tenantID value.TenantID,
) iter.Seq2[*entity.Outbox, error] {
	return func(yield func(*entity.Outbox, error) bool) {
		tn := tablename.New(tenantID)
		tx := es.txManager.Next()
		err := tx.Run(ctx, tn, value.NewPartitionID(0), time.Now(), func() error {
			for e, err := range tx.ListOutboxes(ctx, tn) {
				for !yield(e, err) {
					return err
				}
			}
			return nil
		})
		if err != nil {
			yield(nil, err)
		}
	}
}
