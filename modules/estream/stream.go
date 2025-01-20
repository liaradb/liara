package estream

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/cardboardrobots/liara"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type (
	StreamEventSubscriber struct {
		cons          jetstream.ConsumeContext
		subscriptions []liara.EventHandler
	}
)

func NewStreamEventSubscriber(
	ctx context.Context,
	js jetstream.JetStream,
	streamName string,
	consumerName string,
	subjects ...string,
) (*StreamEventSubscriber, error) {
	s, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: subjects,
	})
	if err != nil {
		return nil, err
	}

	c, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   consumerName,
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return nil, err
	}

	ses := &StreamEventSubscriber{
		subscriptions: nil,
	}

	cons, err := c.Consume(func(msg jetstream.Msg) {
		// TODO: Should we handle this error?
		if err := ses.consume(ctx, msg); err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		return nil, err
	}

	ses.cons = cons

	return ses, nil
}

func (ses *StreamEventSubscriber) Close() error {
	ses.cons.Stop()
	return nil
}

func (ses *StreamEventSubscriber) consume(ctx context.Context, msg jetstream.Msg) error {
	em, err := ses.unmarshalEvent(msg)
	if err != nil {
		return errors.Join(err, msg.Nak())
	}

	for _, s := range ses.subscriptions {
		err := s.Handle(ctx, em)
		if err != nil {
			return errors.Join(err, msg.Nak())
		}
	}

	return msg.Ack()
}

func (*StreamEventSubscriber) unmarshalEvent(msg jetstream.Msg) (liara.Event, error) {
	em := liara.Event{}
	err := json.Unmarshal(msg.Data(), &em)
	return em, err
}

func (ses *StreamEventSubscriber) Subscribe(es liara.EventHandler) func() {
	ses.subscriptions = append(ses.subscriptions, es)

	return func() {
		for i, s := range ses.subscriptions {
			if s == es {
				ses.subscriptions = append(ses.subscriptions[:i], ses.subscriptions[i+1:]...)
				break
			}
		}
	}
}

func Connect(url string) (*nats.Conn, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return nc, nil
}

func connectStream(ctx context.Context, url string, subject string) (*nats.Conn, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	s, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "Events",
		Subjects: []string{subject},
	})
	if err != nil {
		return nil, err
	}

	_, err = js.Publish(ctx, subject, []byte("test"))
	if err != nil {
		return nil, err
	}

	c, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "Consumer",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return nil, err
	}

	cons, err := c.Consume(func(msg jetstream.Msg) {
		msg.Ack()
	})
	if err != nil {
		return nil, err
	}
	defer cons.Stop()

	return nc, nil
}
