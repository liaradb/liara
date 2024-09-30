package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"iter"

	"github.com/cardboardrobots/eventsource"
)

type (
	EventRepository struct {
		db   *sql.DB
		name string
	}

	queryRunner interface {
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}
)

var _ eventsource.EventSource = EventRepository{}

func NewEventRepository(db *sql.DB, name string) EventRepository {
	return EventRepository{
		db,
		name,
	}
}

func (er EventRepository) Append(ctx context.Context, ems ...eventsource.Event) error {
	for _, em := range ems {
		if err := em.Valid(); err != nil {
			return err
		}
	}

	switch len(ems) {
	case 0:
		return nil
	case 1:
		return er.appendEvent(ctx, er.db, ems[0])
	default:
		return er.appendEvents(ctx, ems)
	}
}

func (er EventRepository) appendEvents(ctx context.Context, ems []eventsource.Event) error {
	return runTx(ctx, er.db, &sql.TxOptions{Isolation: sql.LevelDefault}, func(tx *sql.Tx) error {
		for _, em := range ems {
			if err := er.appendEvent(ctx, tx, em); err != nil {
				return err
			}
		}
		return nil
	})
}

func (er EventRepository) appendEvent(ctx context.Context, qr queryRunner, em eventsource.Event) error {
	_, err := qr.ExecContext(ctx, fmt.Sprintf("INSERT INTO %v VALUES( null, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10 )", er.name),
		em.ID,
		em.AggregateID,
		em.AggregateName,
		em.Version,
		em.Name,
		em.CorrelationID,
		em.UserID,
		em.Time,
		em.Schema,
		em.Data,
	)
	return err
}

func (er EventRepository) GetAfterGlobalVersion(ctx context.Context, globalVersion eventsource.GlobalVersion, limit eventsource.Limit) iter.Seq2[eventsource.Event, error] {
	return func(yield func(eventsource.Event, error) bool) {
		var rows *sql.Rows
		var err error

		if limit > 0 {
			rows, err = er.db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %v WHERE global_version > $1 ORDER BY global_version LIMIT $2", er.name),
				globalVersion,
				limit)
		} else {
			rows, err = er.db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %v WHERE global_version > $1 ORDER BY global_version", er.name),
				globalVersion)
		}
		if err != nil {
			yield(eventsource.Event{}, err)
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

func (er EventRepository) Get(ctx context.Context, aggregateID eventsource.AggregateID) iter.Seq2[eventsource.Event, error] {
	return func(yield func(eventsource.Event, error) bool) {
		rows, err := er.db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %v WHERE aggregate_id = $1 ORDER BY global_version", er.name), aggregateID)
		if err != nil {
			yield(eventsource.Event{}, err)
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

func (er EventRepository) GetByAggregateIDAndName(ctx context.Context, aggregateID eventsource.AggregateID, name eventsource.AggregateName) iter.Seq2[eventsource.Event, error] {
	return func(yield func(eventsource.Event, error) bool) {
		rows, err := er.db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %v WHERE aggregate_id = $1 AND aggregate_name = $2 ORDER BY global_version", er.name), aggregateID, name)
		if err != nil {
			yield(eventsource.Event{}, err)
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

func (er EventRepository) Rollback(ctx context.Context, globalVersion eventsource.GlobalVersion) error {
	_, err := er.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %v WHERE global_version > $1", er.name), globalVersion)
	return err
}

func (er EventRepository) scanRow(row Row) (eventsource.Event, error) {
	event := eventsource.Event{}
	err := row.Scan(
		&event.GlobalVersion,
		&event.ID,
		&event.AggregateID,
		&event.AggregateName,
		&event.Version,
		&event.Name,
		&event.CorrelationID,
		&event.UserID,
		&event.Time,
		&event.Schema,
		&event.Data,
	)
	return event, err
}

func (er EventRepository) CreateTable(ctx context.Context) error {
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %v (
	global_version INTEGER PRIMARY KEY AUTOINCREMENT,
	id VARCHAR(50) UNIQUE NOT NULL,
	aggregate_id VARCHAR(50) NOT NULL,
	aggregate_name VARCHAR(50) NOT NULL,
	version BIGINT NOT NULL,
	name VARCHAR(50) NOT NULL,
	correlation_id VARCHAR(50) NOT NULL,
	user_id VARCHAR(50) NOT NULL,
	time TIMESTAMP NOT NULL,
	schema VARCHAR(50) NOT NULL,
	data BLOB,
	CONSTRAINT event_version UNIQUE (aggregate_id, aggregate_name, version)
);
`, er.name)
	_, err := er.db.ExecContext(ctx, query)
	return err
}

func (er EventRepository) CreateIndex(ctx context.Context) error {
	query := fmt.Sprintf(`
CREATE INDEX IF NOT EXISTS index_aggregate_id_aggregate_name ON %v (aggregate_id, aggregate_name);
`, er.name)
	_, err := er.db.ExecContext(ctx, query)
	return err
}
