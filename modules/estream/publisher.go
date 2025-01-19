package estream

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cardboardrobots/liara"
	"github.com/nats-io/nats.go/jetstream"
)

type (
	StreamEventPublisher struct {
		js jetstream.JetStream
	}
)

func NewStreamEventPublisher(
	js jetstream.JetStream,
) *StreamEventPublisher {
	return &StreamEventPublisher{
		js: js,
	}
}

func (ses *StreamEventPublisher) Handle(ctx context.Context, event liara.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	streamName := fmt.Sprintf("%v.%v",
		event.AggregateName.String(),
		event.PartitionID.Value())

	return ses.handleStream(ctx, streamName, data, event.ID.String())
}

func (ses *StreamEventPublisher) handleStream(ctx context.Context, streamName string, data []byte, id string) error {
	_, err := ses.js.Publish(ctx, streamName, data, jetstream.WithMsgID(id))
	return err
}
