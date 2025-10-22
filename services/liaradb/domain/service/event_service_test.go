package service

import (
	"context"
	"errors"
	"testing"

	"github.com/liaradb/liaradb/domain/value"
)

func TestEventService_Append(t *testing.T) {
	t.Run("should not append invalid version", func(t *testing.T) {
		es := NewEventService(nil, nil, nil, nil)

		aggregateID := value.NewAggregateID()
		want := AppendEvent{
			AggregateName: "example",
			ID:            value.NewEventID(),
			AggregateID:   aggregateID,
			Version:       0,
		}

		err := es.Append(context.Background(), "", AppendOptions{}, want)
		if !errors.Is(err, value.ErrAggregateVersionInvalid) {
			t.Error("should return error")
		}
	})
}
