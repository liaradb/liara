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
	transactionRepository TransactionRepository
	eventRepository       EventRepository
	outboxRepository      OutboxRepository
	requestRepository     RequestRepository
}

func NewEventService(
	transactionRepository TransactionRepository,
	eventRepository EventRepository,
	outboxRepository OutboxRepository,
	requestRepository RequestRepository,
) *EventService {
	return &EventService{
		transactionRepository: transactionRepository,
		eventRepository:       eventRepository,
		outboxRepository:      outboxRepository,
		requestRepository:     requestRepository,
	}
}

func (es *EventService) Append(
	ctx context.Context,
	tenantID value.TenantID,
	requestID value.RequestID,
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

	if requestID == "" {
		return es.appendNoRequestID(ctx, tenantID, e...)
	}

	return es.appendRequestID(ctx, tenantID, requestID, e...)
}

func (es *EventService) appendNoRequestID(
	ctx context.Context,
	tenantID value.TenantID,
	e ...AppendEvent,
) error {
	if len(e) == 1 {
		return es.appendEvents(ctx, tenantID, e...)
	}

	return es.transactionRepository.Run(ctx, func(tx Transaction) error {
		return es.appendEvents(ctx, tenantID, e...)
	})
}

func (es *EventService) appendRequestID(
	ctx context.Context,
	tenantID value.TenantID,
	requestID value.RequestID,
	e ...AppendEvent,
) error {
	t := time.Now()
	return es.transactionRepository.Run(ctx, func(tx Transaction) error {
		// TODO: What should this return if requestID is present?
		if ok, err := es.requestRepository.Test(ctx, tenantID, requestID); err != nil || !ok {
			return err
		}

		if err := es.appendEvents(ctx, tenantID, e...); err != nil {
			return err
		}

		return es.requestRepository.Insert(ctx, tenantID, requestID, t)
	})
}

func (es *EventService) appendEvents(
	ctx context.Context,
	tenantID value.TenantID,
	e ...AppendEvent,
) error {
	for _, em := range e {
		err := es.eventRepository.Append(ctx, tenantID, em)
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
