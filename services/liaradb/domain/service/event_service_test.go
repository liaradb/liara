package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/liaradb/liaradb/domain/value"
)

func TestEventService_Append(t *testing.T) {
	t.Parallel()

	t.Run("should not append invalid version", func(t *testing.T) {
		es := NewEventService(nil, nil)

		aggregateID := value.NewAggregateID(uuid.NewString())
		want := AppendEvent{
			AggregateName: value.NewAggregateName("example"),
			// ID:            value.NewEventID(),
			AggregateID: aggregateID,
			Version:     value.NewVersion(0),
		}

		err := es.Append(context.Background(), "", AppendOptions{}, value.NewPartitionID(0), want)
		if !errors.Is(err, value.ErrAggregateVersionInvalid) {
			t.Error("should return error")
		}
	})
}
