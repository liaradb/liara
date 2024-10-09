package liara

import (
	"context"
	"iter"

	"github.com/cardboardrobots/eventsource/entity"
	"github.com/cardboardrobots/eventsource/value"
)

type EventRepository interface {
	Get(ctx context.Context, id value.AggregateID) iter.Seq2[entity.Event, error]
	GetAfterGlobalVersion(context.Context, value.GlobalVersion, value.Limit) iter.Seq2[entity.Event, error]
	GetByAggregateIDAndName(ctx context.Context, id value.AggregateID, name value.AggregateName) iter.Seq2[entity.Event, error]
	Append(ctx context.Context, e ...entity.AppendEvent) error
}
