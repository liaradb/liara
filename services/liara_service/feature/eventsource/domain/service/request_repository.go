package service

import (
	"context"
	"time"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type RequestRepository interface {
	CreateTable(context.Context, value.TenantID) error
	DropTable(context.Context, value.TenantID) error
	Test(context.Context, value.TenantID, value.RequestID) (bool, error)
	Insert(context.Context, value.TenantID, value.RequestID, time.Time) error
	Purge(context.Context, value.TenantID, time.Time) error
}
