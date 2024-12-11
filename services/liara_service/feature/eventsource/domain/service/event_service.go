package service

import (
	"context"
	"iter"
	"time"

	"github.com/cardboardrobots/liara_service/feature/base"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type EventService struct {
	transactionContainer TransactionContainer
	eventRepository      EventRepository
	outboxRepository     OutboxRepository
	requestRepository    RequestRepository
}

func NewEventService(
	transactionRepository TransactionContainer,
	eventRepository EventRepository,
	outboxRepository OutboxRepository,
	requestRepository RequestRepository,
) *EventService {
	return &EventService{
		transactionContainer: transactionRepository,
		eventRepository:      eventRepository,
		outboxRepository:     outboxRepository,
		requestRepository:    requestRepository,
	}
}

type AppendOptions struct {
	RequestID     value.RequestID     // The ID of the Request, for idempotency
	CorrelationID value.CorrelationID // The ID of the entire Command and Event chain
	UserID        value.UserID        // The ID of the User issuing the Command
	Time          time.Time           // The Time this Event was created
}

func (ao *AppendOptions) toEventMetadata() entity.EventMetadata {
	return entity.EventMetadata{
		UserID:        ao.UserID,
		CorrelationID: ao.CorrelationID,
		Time:          ao.Time,
	}
}

type AppendEvent struct {
	ID            value.EventID       // The ID of the Event, used for de-duplication
	AggregateName value.AggregateName // The Name of the Aggregate
	AggregateID   value.AggregateID   // The ID of the Aggregate to which this Event applies
	Version       value.Version       // The Version of the Aggregate
	PartitionID   value.PartitionID   // The ID to partition Events
	Name          value.EventName     // The Name of the Event
	Schema        value.Schema        // The Schema for the internal data
	Data          []byte              // The internal data of the Event
}

func (ae *AppendEvent) Valid() error {
	if ae.Version < 1 {
		return value.ErrAggregateVersionInvalid
	}

	return nil
}

func (ae *AppendEvent) toEvent(options AppendOptions) entity.Event {
	id := ae.ID
	if id == "" {
		id = value.NewEventID()
	}

	return entity.Event{
		GlobalVersion: 0,
		ID:            id,
		AggregateName: ae.AggregateName,
		AggregateID:   ae.AggregateID,
		Version:       ae.Version,
		PartitionID:   ae.PartitionID,
		Name:          ae.Name,
		Schema:        ae.Schema,
		Metadata:      options.toEventMetadata(),
		Data:          ae.Data,
	}
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
	if len(e) == 1 {
		return es.appendEvents(ctx, tenantID, options, e...)
	}

	return es.transactionContainer.Run(ctx, func() error {
		return es.appendEvents(ctx, tenantID, options, e...)
	})
}

func (es *EventService) appendRequestID(
	ctx context.Context,
	tenantID value.TenantID,
	options AppendOptions,
	e ...AppendEvent,
) error {
	t := time.Now()
	return es.transactionContainer.Run(ctx, func() error {
		// TODO: What should this return if requestID is present?
		if ok, err := es.requestRepository.Test(ctx, tenantID, options.RequestID); err != nil || !ok {
			return err
		}

		if err := es.appendEvents(ctx, tenantID, options, e...); err != nil {
			return err
		}

		return es.requestRepository.Insert(ctx, tenantID, options.RequestID, t)
	})
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
		err := es.eventRepository.Append(ctx, tenantID, em.toEvent(options))
		if err != nil {
			return err
		}
	}
	return nil
}

func (es *EventService) Get(
	ctx context.Context,
	tenantID value.TenantID,
	id value.AggregateID,
) iter.Seq2[entity.Event, error] {
	return es.eventRepository.Get(ctx, tenantID, id)
}

func (es *EventService) GetByAggregateIDAndName(
	ctx context.Context,
	tenantID value.TenantID,
	id value.AggregateID,
	name value.AggregateName,
) iter.Seq2[entity.Event, error] {
	return es.eventRepository.GetByAggregateIDAndName(ctx, tenantID, id, name)
}

func (es *EventService) GetAfterGlobalVersion(
	ctx context.Context,
	tenantID value.TenantID,
	version value.GlobalVersion,
	partitionRange value.PartitionRange,
	limit value.Limit,
) iter.Seq2[entity.Event, error] {
	return es.eventRepository.GetAfterGlobalVersion(ctx, tenantID, version, partitionRange, limit)
}

func (es *EventService) GetByOutbox(
	ctx context.Context,
	tenantID value.TenantID,
	outboxID value.OutboxID,
	limit value.Limit,
) iter.Seq2[entity.Event, error] {
	outbox, err := es.outboxRepository.GetOutbox(ctx, tenantID, outboxID)
	if err != nil {
		return base.IterError[entity.Event](err)
	}

	return es.eventRepository.GetAfterGlobalVersion(ctx, tenantID, outbox.GlobalVersion(), outbox.PartitionRange(), limit)
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
