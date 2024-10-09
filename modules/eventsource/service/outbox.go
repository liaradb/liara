package service

import (
	"context"
	"errors"
	"time"

	"github.com/cardboardrobots/eventsource/entity"
	"github.com/cardboardrobots/eventsource/value"
)

type (
	Outbox struct {
		outboxRepository OutboxRepository
		eventRepository  EventRepository
		subscriptions    []EventSubscriber
	}

	EventSubscriber interface {
		Handle(context.Context, entity.Event) error
	}

	// TODO: Where is this used?
	EventPublisher interface {
		Init(context.Context) (func() error, error)
		Subscribe(EventSubscriber) func()
	}
)

func NewOutbox(
	outboxRepository OutboxRepository,
	eventRepository EventRepository,
) Outbox {
	return Outbox{
		outboxRepository,
		eventRepository,
		nil,
	}
}

func (o *Outbox) Run(ctx context.Context, outboxID value.OutboxID, duration time.Duration, limit value.Limit) {
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			o.read(ctx, outboxID, limit)
		}
	}()
}

func (o *Outbox) read(ctx context.Context, outboxID value.OutboxID, limit value.Limit) error {
	globalVersion, err := o.outboxRepository.GetOrCreateOutbox(ctx, outboxID)
	if err != nil {
		return err
	}

	updatedGlobalVersion := globalVersion
	for em, err := range o.eventRepository.GetAfterGlobalVersion(ctx, globalVersion, limit) {
		if err != nil {
			return err
		}

		for _, s := range o.subscriptions {
			err := s.Handle(ctx, em)
			if err != nil {
				return err
			}
		}
		updatedGlobalVersion = em.GlobalVersion
	}
	if updatedGlobalVersion == globalVersion {
		return err
	}

	updateErr := o.outboxRepository.UpdateOutboxPosition(ctx, outboxID, updatedGlobalVersion)
	return errors.Join(err, updateErr)
}

func (o *Outbox) Subscribe(es EventSubscriber) func() {
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
