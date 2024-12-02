package service

import (
	"context"
	"iter"
	"time"

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
	return es.transactionRepository.Run(ctx, func(tx Transaction) error {
		err := es.requestRepository.Test(ctx, requestID)
		if err != nil {
			return err
		}

		t := time.Now()
		err = es.eventRepository.Append(ctx, e...)
		if err != nil {
			return err
		}

		err = es.requestRepository.Insert(ctx, requestID, t)
		return err
	})
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
	partitionID value.PartitionID,
	limit value.Limit,
) iter.Seq2[entity.Event, error] {
	return es.eventRepository.GetAfterGlobalVersion(ctx, version, partitionID, limit)
}

func (es *EventService) GetOrCreateOutbox(
	ctx context.Context,
	outboxID value.OutboxID,
) (value.GlobalVersion, error) {
	return es.outboxRepository.GetOrCreateOutbox(ctx, outboxID)
}

func (es *EventService) UpdateOutboxPosition(
	ctx context.Context,
	outboxID value.OutboxID,
	globalVersion value.GlobalVersion,
) error {
	return es.outboxRepository.UpdateOutboxPosition(ctx, outboxID, globalVersion)
}
