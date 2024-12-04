package service

import (
	"context"
	"errors"
	"testing"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

func TestEventService_Append(t *testing.T) {
	t.Run("should not append invalid version", func(t *testing.T) {
		es := NewEventService(nil, nil, nil, nil)

		aggregateID := value.AggregateID("aggregateID")
		want := AppendEvent{
			AggregateName: "example",
			ID:            "eventID",
			AggregateID:   aggregateID,
			Version:       0,
		}

		err := es.Append(context.Background(), "", "", want)
		if !errors.Is(err, value.ErrAggregateVersionInvalid) {
			t.Error("should return error")
		}
	})
}
