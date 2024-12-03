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
		db       *sql.DB
		tenantID value.TenantID
	}

	requestModel struct {
		ID   value.RequestID
		Time time.Time
	}
)

var _ service.RequestRepository = &RequestRepository{}

func NewRequestRepository(db *sql.DB, tenantID value.TenantID) *RequestRepository {
	return &RequestRepository{
		db,
		tenantID,
	}
}

func (*RequestRepository) getName(tenantID value.TenantID) string {
	n := tenantID.String()
	if n == "" {
		n = "default"
	}
	return fmt.Sprintf("__%v__requests", n)
}

func (r *RequestRepository) Insert(ctx context.Context, requestID value.RequestID, time time.Time) error {
	return r.insertRow(ctx, r.db, entity.RequestLog{ID: requestID, Time: time})
}

func (r *RequestRepository) Purge(ctx context.Context, time time.Time) error {
	_, err := r.db.ExecContext(ctx, fmt.Sprintf(`
DELETE FROM %v WHERE
time <= $1
`,
		r.getName(r.tenantID)), time)
	return err
}

func (r *RequestRepository) Test(ctx context.Context, requestID value.RequestID) (bool, error) {
	row := r.db.QueryRowContext(ctx, fmt.Sprintf(`
SELECT * FROM %v WHERE
id = $1
`,
		r.getName(r.tenantID)), requestID)
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
		r.getName(r.tenantID)),
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

func (r *RequestRepository) CreateTable(ctx context.Context, tenantID value.TenantID) error {
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %v (
	id VARCHAR(50) PRIMARY KEY NOT NULL,
	time TIMESTAMP NOT NULL,
);
`,
		r.getName(tenantID))
	_, err := r.db.ExecContext(ctx, query)
	return err
}
