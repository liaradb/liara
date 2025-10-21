package service

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type MockEventRepository struct {
}

func (mer *MockEventRepository) Get(
	ctx context.Context,
	aggregateID value.AggregateID,
) iter.Seq2[entity.Event, error] {
	return func(yield func(entity.Event, error) bool) {}
}

func (mer *MockEventRepository) GetAfterGlobalVersion(
	ctx context.Context,
	globalVersion value.GlobalVersion,
	partitionRange value.PartitionRange,
	limit value.Limit,
) iter.Seq2[entity.Event, error] {
	return func(yield func(entity.Event, error) bool) {}
}

func (mer *MockEventRepository) GetByAggregateIDAndName(
	ctx context.Context,
	aggregateID value.AggregateID,
	aggregateName value.AggregateName,
) iter.Seq2[entity.Event, error] {
	return func(yield func(entity.Event, error) bool) {}
}

func (mer *MockEventRepository) Append(
	ctx context.Context,
	e AppendEvent,
) error {
	return nil
}
