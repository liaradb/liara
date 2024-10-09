package estream

import (
	"context"
	"encoding/json"

	"github.com/cardboardrobots/eventsource/entity"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type (
	StreamEventPublisher struct {
		nc         *nats.Conn
		js         jetstream.JetStream
		streamName string
	}
)

func NewStreamEventPublisher(
	nc *nats.Conn,
	js jetstream.JetStream,
	streamName string,
) *StreamEventPublisher {
	return &StreamEventPublisher{
		nc:         nc,
		js:         js,
		streamName: streamName,
	}
}

const useStream = false

func (ses *StreamEventPublisher) Handle(ctx context.Context, event entity.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if useStream {
		return ses.handleStream(ctx, data, event.ID.String())
	} else {
		return ses.handleQueue(data)
	}
}

func (ses *StreamEventPublisher) handleStream(ctx context.Context, data []byte, id string) error {
	_, err := ses.js.Publish(ctx, ses.streamName, data, jetstream.WithMsgID(id))
	return err
}

func (ses *StreamEventPublisher) handleQueue(data []byte) error {
	return ses.nc.Publish(ses.streamName, data)
}
