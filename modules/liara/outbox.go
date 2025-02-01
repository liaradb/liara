package liara

import (
	"context"
	"log"
	"time"
)

type (
	Outbox struct {
		tenantID         TenantID
		outboxRepository OutboxRepository
		eventRepository  EventRepository
		subscriptions    []EventHandler
	}

	EventHandler interface {
		Handle(context.Context, Event) error
	}
)

func NewOutbox(
	outboxRepository OutboxRepository,
	eventRepository EventRepository,
) *Outbox {
	return &Outbox{
		tenantID:         "",
		outboxRepository: outboxRepository,
		eventRepository:  eventRepository,
		subscriptions:    nil,
	}
}

func (o *Outbox) Create(ctx context.Context, outboxID OutboxID, partitionIDs []PartitionID) (OutboxID, error) {
	return o.outboxRepository.CreateOutbox(ctx, o.tenantID, outboxID, partitionIDs)
}

func (o *Outbox) Run(ctx context.Context, outboxID OutboxID, duration time.Duration, limit Limit) {
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			// TODO: Should we handle this error?
			if err := o.read(ctx, outboxID, limit); err != nil {
				log.Println(err)
			}
		}
	}()
}

func (o *Outbox) read(ctx context.Context, outboxID OutboxID, limit Limit) error {
	update := false
	var globalVersion GlobalVersion
	for em, err := range o.eventRepository.GetByOutbox(ctx, o.tenantID, outboxID, limit) {
		if err != nil {
			return err
		}

		for _, s := range o.subscriptions {
			err := s.Handle(ctx, em)
			if err != nil {
				return err
			}
		}
		update = true
		globalVersion = em.GlobalVersion
	}

	if update {
		return o.outboxRepository.UpdateOutboxPosition(ctx, o.tenantID, outboxID, globalVersion)
	}
	return nil
}

func (o *Outbox) Subscribe(es EventHandler) func() {
	o.subscriptions = append(o.subscriptions, es)

	return func() {
		for i, s := range o.subscriptions {
			if s == es {
				o.subscriptions = append(o.subscriptions[:i], o.subscriptions[i+1:]...)
				break
			}
		}
	}
}
