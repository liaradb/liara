package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
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

func (r *RequestRepository) Insert(ctx context.Context, requestID value.RequestID, time time.Time) error {
	return r.insertRow(ctx, r.db, entity.RequestLog{ID: requestID, Time: time})
}

func (r *RequestRepository) Purge(ctx context.Context, time time.Time) error {
	_, err := r.db.ExecContext(ctx, fmt.Sprintf(`
DELETE FROM %v WHERE
time <= $1
`,
		r.name), time)
	return err
}

func (r *RequestRepository) Test(ctx context.Context, requestID value.RequestID) (bool, error) {
	row := r.db.QueryRowContext(ctx, fmt.Sprintf(`
SELECT * FROM %v WHERE
id = $1
`,
		r.name), requestID)
	_, err := r.scanRow(row)
	if err == sql.ErrNoRows {
		return true, nil
	}

	return false, err
}

func (r RequestRepository) insertRow(
	ctx context.Context,
	qr queryRunner,
	request entity.RequestLog,
) error {
	_, err := qr.ExecContext(ctx, fmt.Sprintf(`
INSERT INTO %v VALUES( $1, $2 )
`,
		r.name),
		request.ID,
		request.Time,
	)
	return err
}

func (r RequestRepository) scanRow(row Row) (entity.RequestLog, error) {
	request := entity.RequestLog{}
	err := row.Scan(
		&request.ID,
		&request.Time,
	)
	return request, err
}

func (r *RequestRepository) CreateTable(ctx context.Context) error {
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %v (
	id VARCHAR(50) PRIMARY KEY NOT NULL,
	time TIMESTAMP NOT NULL,
);
`,
		r.name)
	_, err := r.db.ExecContext(ctx, query)
	return err
}
