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
		return es.appendNoRequestID(ctx, e...)
	}

	return es.appendRequestID(ctx, requestID, e...)
}

func (es *EventService) appendNoRequestID(
	ctx context.Context,
	e ...AppendEvent,
) error {
	if len(e) == 1 {
		return es.appendEvents(ctx, e...)
	}

	return es.transactionRepository.Run(ctx, func(tx Transaction) error {
		return es.appendEvents(ctx, e...)
	})
}

func (es *EventService) appendRequestID(
	ctx context.Context,
	requestID value.RequestID,
	e ...AppendEvent,
) error {
	t := time.Now()
	return es.transactionRepository.Run(ctx, func(tx Transaction) error {
		// TODO: What should this return if requestID is present?
		if ok, err := es.requestRepository.Test(ctx, requestID); err != nil || !ok {
			return err
		}

		if err := es.appendEvents(ctx, e...); err != nil {
			return err
		}

		return es.requestRepository.Insert(ctx, requestID, t)
	})
}

func (es *EventService) appendEvents(
	ctx context.Context,
	e ...AppendEvent,
) error {
	for _, em := range e {
		err := es.eventRepository.Append(ctx, em)
		if err != nil {
			return err
		}
	}
	return nil
}

func (es *EventService) Get(
	ctx context.Context,
	id value.AggregateID,
) iter.Seq2[entity.Event, error] {
	return es.eventRepository.Get(ctx, id)
}

func (es *EventService) GetByAggregateIDAndName(
	ctx context.Context,
	id value.AggregateID,
	name value.AggregateName,
) iter.Seq2[entity.Event, error] {
	return es.eventRepository.GetByAggregateIDAndName(ctx, id, name)
}

func (es *EventService) GetAfterGlobalVersion(
	ctx context.Context,
	version value.GlobalVersion,
	partitionRange value.PartitionRange,
	limit value.Limit,
) iter.Seq2[entity.Event, error] {
	return es.eventRepository.GetAfterGlobalVersion(ctx, version, partitionRange, limit)
}

func (es *EventService) GetByOutbox(
	ctx context.Context,
	outboxID value.OutboxID,
	limit value.Limit,
) iter.Seq2[entity.Event, error] {
	outbox, err := es.outboxRepository.GetOutbox(ctx, outboxID)
	if err != nil {
		return base.IterError[entity.Event](err)
	}

	return es.eventRepository.GetAfterGlobalVersion(ctx, outbox.GlobalVersion(), outbox.PartitionRange(), limit)
}

func (es *EventService) CreateOutbox(
	ctx context.Context,
	outboxID value.OutboxID,
	partitionRange value.PartitionRange,
) (value.OutboxID, error) {
	outbox := entity.NewOutbox(outboxID, partitionRange)
	err := es.outboxRepository.CreateOutbox(ctx, outbox)
	return outbox.ID(), err
}

func (es *EventService) GetOutbox(
	ctx context.Context,
	outboxID value.OutboxID,
) (*entity.Outbox, error) {
	return es.outboxRepository.GetOutbox(ctx, outboxID)
}

func (es *EventService) UpdateOutboxPosition(
	ctx context.Context,
	outboxID value.OutboxID,
	globalVersion value.GlobalVersion,
) error {
	return es.outboxRepository.UpdateOutboxPosition(ctx, outboxID, globalVersion)
}
