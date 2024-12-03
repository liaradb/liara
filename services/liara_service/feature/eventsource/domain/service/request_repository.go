package service

import (
	"context"
	"time"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type RequestRepository interface {
	CreateTable(context.Context, value.TenantID) error
	Test(context.Context, value.RequestID) (bool, error)
	Insert(context.Context, value.RequestID, time.Time) error
	Purge(context.Context, time.Time) error
}
