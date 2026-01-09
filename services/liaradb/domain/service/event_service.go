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
	outboxRepository  OutboxRepository
	requestRepository RequestRepository
	txManager         *transaction.Manager
}

func NewEventService(
	outboxRepository OutboxRepository,
	requestRepository RequestRepository,
	txManager *transaction.Manager,
) *EventService {
	return &EventService{
		outboxRepository:  outboxRepository,
		requestRepository: requestRepository,
		txManager:         txManager,
	}
}

type AppendOptions struct {
	RequestID     value.RequestID     // The ID of the Request, for idempotency
	CorrelationID value.CorrelationID // The ID of the entire Command and Event chain
	UserID        value.UserID        // The ID of the User issuing the Command
	Time          time.Time           // The Time this Event was created
}

func (ao *AppendOptions) toMetadata() entity.Metadata {
	return entity.Metadata{
		UserID:        ao.UserID,
		CorrelationID: ao.CorrelationID,
		Time:          value.NewTime(ao.Time),
	}
}

type AppendEvent struct {
	ID            string              // The ID of the Event, used for de-duplication
	AggregateName value.AggregateName // The Name of the Aggregate
	AggregateID   value.AggregateID   // The ID of the Aggregate to which this Event applies
	Version       value.Version       // The Version of the Aggregate
	PartitionID   value.PartitionID   // The ID to partition Events
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

func (ae *AppendEvent) toEvent(options AppendOptions) (entity.Event, error) {
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
		PartitionID:   ae.PartitionID,
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

	if options.RequestID == "" {
		return es.appendNoRequestID(ctx, tenantID, options, e...)
	}

	return es.appendRequestID(ctx, tenantID, options, e...)
}

func (es *EventService) appendNoRequestID(
	ctx context.Context,
	tenantID value.TenantID,
	options AppendOptions,
	e ...AppendEvent,
) error {
	// if len(e) == 1 {
	return es.appendEvents(ctx, tenantID, options, e...)
	// }

	// return es.transactionContainer.Run(ctx, func() error {
	// 	return es.appendEvents(ctx, tenantID, options, e...)
	// })
}

func (es *EventService) appendRequestID(
	ctx context.Context,
	tenantID value.TenantID,
	options AppendOptions,
	e ...AppendEvent,
) error {
	panic("unimplemented")
	// t := time.Now()
	// return es.transactionContainer.Run(ctx, func() error {
	// 	// TODO: What should this return if requestID is present?
	// 	if ok, err := es.requestRepository.Test(ctx, tenantID, options.RequestID); err != nil || !ok {
	// 		return err
	// 	}

	// 	if err := es.appendEvents(ctx, tenantID, options, e...); err != nil {
	// 		return err
	// 	}

	// 	return es.requestRepository.Insert(ctx, tenantID, options.RequestID, t)
	// })
}

func (es *EventService) appendEvents(
	ctx context.Context,
	tenantID value.TenantID,
	options AppendOptions,
	e ...AppendEvent,
) error {
	if options.Time.IsZero() {
		options.Time = time.Now()
	}

	for _, em := range e {
		event, err := em.toEvent(options)
		if err != nil {
			return err
		}

		err = es.append(ctx, tenantID, event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (es *EventService) append(
	ctx context.Context,
	tenantID value.TenantID,
	e entity.Event, // TODO: Should this be a pointer?
) error {
	tx := es.txManager.Next()
	tn := tablename.New(tenantID)
	// TODO: PartitionID should be on the transaction, not just the Event
	return tx.Run(ctx, tn, e.PartitionID, time.Now(), func() error {
		buf := bytes.NewBuffer(nil)
		if err := e.Write(buf); err != nil {
			return err
		}

		return tx.Insert(ctx,
			tn,
			time.Now(),
			&e,
			buf.Bytes(),
		)
	})
}

func (es *EventService) TestIdempotency(
	ctx context.Context,
	tenantID value.TenantID,
	id value.RequestID,
) (bool, error) {
	return es.requestRepository.Test(ctx, tenantID, id)
}

func (es *EventService) Get(
	ctx context.Context,
	tenantID value.TenantID,
	partitionID value.PartitionID,
	id value.AggregateID,
) iter.Seq2[entity.Event, error] { // TODO: Should this be a pointer?
	return func(yield func(entity.Event, error) bool) {
		tn := tablename.New(tenantID)
		tx := es.txManager.Next()
		for e, err := range tx.GetAggregate(ctx, tn, partitionID, id) {
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

func (es *EventService) GetByAggregateIDAndName(
	ctx context.Context,
	tenantID value.TenantID,
	id value.AggregateID,
	name value.AggregateName,
) iter.Seq2[entity.Event, error] {
	panic("unimplemented")
}

func (es *EventService) GetAfterGlobalVersion(
	ctx context.Context,
	tenantID value.TenantID,
	version value.GlobalVersion,
	partitionRange value.PartitionRange,
	limit value.Limit,
) iter.Seq2[entity.Event, error] {
	panic("unimplemented")
}

func (es *EventService) GetByOutbox(
	ctx context.Context,
	tenantID value.TenantID,
	outboxID value.OutboxID,
	limit value.Limit,
) iter.Seq2[entity.Event, error] {
	_, err := es.outboxRepository.GetOutbox(ctx, tenantID, outboxID)
	if err != nil {
		return iterator.Error[entity.Event](err)
	}

	// return es.eventRepository.GetAfterGlobalVersion(ctx, tenantID, outbox.GlobalVersion(), outbox.PartitionRange(), limit)
	panic("unimplemented")
}

func (es *EventService) CreateOutbox(
	ctx context.Context,
	tenantID value.TenantID,
	outboxID value.OutboxID,
	partitionRange value.PartitionRange,
) (value.OutboxID, error) {
	outbox := entity.NewOutbox(outboxID, partitionRange)
	err := es.outboxRepository.CreateOutbox(ctx, tenantID, outbox)
	return outbox.ID(), err
}

func (es *EventService) GetOutbox(
	ctx context.Context,
	tenantID value.TenantID,
	outboxID value.OutboxID,
) (*entity.Outbox, error) {
	return es.outboxRepository.GetOutbox(ctx, tenantID, outboxID)
}

func (es *EventService) UpdateOutboxPosition(
	ctx context.Context,
	tenantID value.TenantID,
	outboxID value.OutboxID,
	globalVersion value.GlobalVersion,
) error {
	return es.outboxRepository.UpdateOutboxPosition(ctx, tenantID, outboxID, globalVersion)
}
