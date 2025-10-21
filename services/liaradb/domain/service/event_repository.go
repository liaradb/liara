package service

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type EventRepository interface {
	CreateTable(context.Context, value.TenantID) error
	CreateIndex(context.Context, value.TenantID) error
	DropTable(context.Context, value.TenantID) error
	Get(context.Context, value.TenantID, value.AggregateID) iter.Seq2[entity.Event, error]
	GetAfterGlobalVersion(context.Context, value.TenantID, value.GlobalVersion, value.PartitionRange, value.Limit) iter.Seq2[entity.Event, error]
	GetByAggregateIDAndName(context.Context, value.TenantID, value.AggregateID, value.AggregateName) iter.Seq2[entity.Event, error]
	Append(context.Context, value.TenantID, entity.Event) error
}
