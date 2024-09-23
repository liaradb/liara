package estream

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cardboardrobots/eventsource"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// TODO: Move these out of library
const StreamEvents = "events"
const QueueGroup = "queue"

type (
	StreamEventSubscriber struct {
		nc            *nats.Conn
		js            jetstream.JetStream
		subscriptions []eventsource.EventSubscriber
		streamName    string
		queueName     string
	}
)

func NewStreamEventSubscriber(
	nc *nats.Conn,
	js jetstream.JetStream,
	streamName string,
	queueName string,
) *StreamEventSubscriber {
	return &StreamEventSubscriber{
		nc:            nc,
		js:            js,
		subscriptions: nil,
		streamName:    streamName,
		queueName:     queueName,
	}
}

func (ses *StreamEventSubscriber) Init(ctx context.Context) (func() error, error) {
	if useStream {
		return ses.streamInit(ctx)
	} else {
		return ses.queueInit(ctx)
	}
}

func (ses *StreamEventSubscriber) streamInit(ctx context.Context) (func() error, error) {
	s, err := ses.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "Events",
		Subjects: []string{StreamEvents},
	})
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
		em := eventsource.Event{}
		err := json.Unmarshal(msg.Data(), &em)
		if err != nil {
			log.Println(err)
			msg.Nak()
			return
		}

		for _, s := range ses.subscriptions {
			err := s.Handle(ctx, em)
			if err != nil {
				log.Println(err)
				msg.Nak()
				return
			}
		}

		msg.Ack()
	})
	if err != nil {
		return nil, err
	}

	return func() error { cons.Stop(); return nil }, nil
}

func (ses *StreamEventSubscriber) queueInit(ctx context.Context) (func() error, error) {
	sub, err := ses.nc.QueueSubscribe(ses.streamName, ses.queueName, func(msg *nats.Msg) {
		em := eventsource.Event{}
		err := json.Unmarshal(msg.Data, &em)
		if err != nil {
			log.Println(err)
			msg.Nak()
			return
		}

		for _, s := range ses.subscriptions {
			err := s.Handle(ctx, em)
			if err != nil {
				log.Println(err)
				msg.Nak()
				return
			}
		}

		msg.Ack()
	})
	if err != nil {
		return nil, err
	}

	return sub.Unsubscribe, nil
}

func (ses *StreamEventSubscriber) Subscribe(es eventsource.EventSubscriber) func() {
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

func ConnectStream(ctx context.Context, url string) (*nats.Conn, error) {
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
		Subjects: []string{StreamEvents},
	})
	if err != nil {
		return nil, err
	}

	_, err = js.Publish(ctx, StreamEvents, []byte("test"))
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
