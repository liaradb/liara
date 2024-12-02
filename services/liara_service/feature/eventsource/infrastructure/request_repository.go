package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/service"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type (
	RequestRepository struct {
		db   *sql.DB
		name string
	}

	requestModel struct {
		ID   value.RequestID
		Time time.Time
	}
)

var _ service.RequestRepository = &RequestRepository{}

func NewRequestRepository(db *sql.DB, name string) *RequestRepository {
	return &RequestRepository{
		db,
		name,
	}
}

func (r *RequestRepository) Insert(context.Context, value.RequestID, time.Time) error {
	return nil
}

func (r *RequestRepository) Purge(context.Context, time.Time) error {
	return nil
}

func (r *RequestRepository) Test(context.Context, value.RequestID) error {
	return nil
}
