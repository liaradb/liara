package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"iter"
	"strings"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/service"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type (
	EventRepository struct {
		db *sql.DB
	}

	queryRunner interface {
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}
)

var _ service.EventRepository = EventRepository{}

// TODO: Change to pointer
func NewEventRepository(
	db *sql.DB,
) EventRepository {
	return EventRepository{
		db,
	}
}

func (*EventRepository) getName(tenantID value.TenantID) string {
	n := tenantID.String()
	if n == "" {
		n = "default"
	}
	return fmt.Sprintf("__%v__events", n)
}

func (er EventRepository) Append(
	ctx context.Context,
	tenantID value.TenantID,
	em service.AppendEvent,
) error {
	_, err := er.db.ExecContext(ctx, fmt.Sprintf(`
INSERT INTO %v VALUES( null, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11 )
`,
		er.getName(tenantID)),
		em.ID,
		em.AggregateName,
		em.AggregateID,
		em.Version,
		em.PartitionID,
		em.Name,
		em.Schema,
		em.Metadata.CorrelationID,
		em.Metadata.UserID,
		em.Metadata.Time,
		em.Data,
	)
	return err
}

func (er EventRepository) GetAfterGlobalVersion(
	ctx context.Context,
	tenantID value.TenantID,
	globalVersion value.GlobalVersion,
	partitionRange value.PartitionRange,
	limit value.Limit,
) iter.Seq2[entity.Event, error] {
	return func(yield func(entity.Event, error) bool) {
		var rows *sql.Rows
		var err error

		b := strings.Builder{}
		b.WriteString(fmt.Sprintf("SELECT * FROM %v", er.getName(tenantID)))
		b.WriteString(" WHERE global_version > $1")
		partition0, partition1 := partitionRange.All()
		if partition0 > 0 {
			b.WriteString(" AND partition_id >= $2")
		}
		if partition1 > partition0 {
			b.WriteString(" AND partition_id <= $3")
		}
		b.WriteString(" ORDER BY global_version")
		if limit > 0 {
			b.WriteString(" LIMIT $4")
		}

		rows, err = er.db.QueryContext(ctx, b.String(),
			globalVersion,
			partition0,
			partition1,
			limit)
		if err != nil {
			yield(entity.Event{}, err)
			return
		}

		defer func() { err = rows.Close() }()
		for rows.Next() {
			event, err := er.scanRow(rows)
			if !yield(event, err) {
				return
			}
		}
	}
}

func (er EventRepository) Get(
	ctx context.Context,
	tenantID value.TenantID,
	aggregateID value.AggregateID,
) iter.Seq2[entity.Event, error] {
	return func(yield func(entity.Event, error) bool) {
		rows, err := er.db.QueryContext(ctx, fmt.Sprintf(`
SELECT * FROM %v WHERE
aggregate_id = $1
ORDER BY global_version
`,
			er.getName(tenantID)), aggregateID)
		if err != nil {
			yield(entity.Event{}, err)
			return
		}

		defer func() { err = rows.Close() }()
		for rows.Next() {
			event, err := er.scanRow(rows)
			if !yield(event, err) {
				return
			}
		}
	}
}

func (er EventRepository) GetByAggregateIDAndName(
	ctx context.Context,
	tenantID value.TenantID,
	aggregateID value.AggregateID,
	name value.AggregateName,
) iter.Seq2[entity.Event, error] {
	return func(yield func(entity.Event, error) bool) {
		rows, err := er.db.QueryContext(ctx, fmt.Sprintf(`
SELECT * FROM %v 
WHERE aggregate_id = $1
AND aggregate_name = $2
ORDER BY global_version
`,
			er.getName(tenantID)), aggregateID, name)
		if err != nil {
			yield(entity.Event{}, err)
			return
		}

		defer func() { err = rows.Close() }()
		for rows.Next() {
			event, err := er.scanRow(rows)
			if !yield(event, err) {
				return
			}
		}
	}
}

func (er EventRepository) Rollback(
	ctx context.Context,
	tenantID value.TenantID,
	gv value.GlobalVersion,
) error {
	q := fmt.Sprintf(`
DELETE FROM %v
WHERE global_version > $1
`,
		er.getName(tenantID))
	_, err := er.db.ExecContext(ctx, q, gv)
	return err
}

func (er EventRepository) scanRow(row Row) (entity.Event, error) {
	event := entity.Event{}
	err := row.Scan(
		&event.GlobalVersion,
		&event.ID,
		&event.AggregateName,
		&event.AggregateID,
		&event.Version,
		&event.PartitionID,
		&event.Name,
		&event.Schema,
		&event.Metadata.CorrelationID,
		&event.Metadata.UserID,
		&event.Metadata.Time,
		&event.Data,
	)
	return event, err
}

func (er EventRepository) CreateTable(ctx context.Context, tenantID value.TenantID) error {
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %v (
	global_version INTEGER PRIMARY KEY AUTOINCREMENT,
	id VARCHAR(50) UNIQUE NOT NULL,
	aggregate_name VARCHAR(50) NOT NULL,
	aggregate_id VARCHAR(50) NOT NULL,
	version BIGINT NOT NULL,
	partition_id INT NOT NULL,
	name VARCHAR(50) NOT NULL,
	schema VARCHAR(50) NOT NULL,
	correlation_id VARCHAR(50) NOT NULL,
	user_id VARCHAR(50) NOT NULL,
	time TIMESTAMP NOT NULL,
	data BLOB,
	CONSTRAINT event_version UNIQUE (aggregate_name, aggregate_id, version)
);
`,
		er.getName(tenantID))
	_, err := er.db.ExecContext(ctx, query)
	return err
}

func (er EventRepository) CreateIndex(ctx context.Context, tenantID value.TenantID) error {
	query := fmt.Sprintf(`
CREATE INDEX IF NOT EXISTS index_aggregate_id_aggregate_name
ON %v (aggregate_id, aggregate_name);
`,
		er.getName(tenantID))
	_, err := er.db.ExecContext(ctx, query)
	return err
}
