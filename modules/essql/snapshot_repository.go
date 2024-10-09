package essql

import (
	"context"
	"database/sql"

	"github.com/cardboardrobots/eventsource/value"
)

type (
	SnapshotModel[U ~string] struct {
		AggregateID U
		EventID     value.EventID
		Version     value.Version
		Schema      value.Schema
		Data        string
	}

	SnapshotRepository[U ~string] struct {
		db *sql.DB
	}
)

func NewSnapshotRepository[U ~string](db *sql.DB) SnapshotRepository[U] {
	return SnapshotRepository[U]{
		db,
	}
}

func (sr SnapshotRepository[U]) Insert(ctx context.Context, sm SnapshotModel[U]) error {
	_, err := sr.db.ExecContext(ctx, "INSERT INTO snapshots VALUES( $1, $2, $3, $4, $5 )",
		sm.AggregateID,
		sm.EventID,
		sm.Version,
		sm.Schema,
		sm.Data,
	)
	return err
}

func (sr SnapshotRepository[U]) GetByAggregateIDAndVersion(ctx context.Context, aggregateID U, version int) (SnapshotModel[U], error) {
	row := sr.db.QueryRowContext(ctx, "SELECT * FROM snapshots WHERE id = $1 AND version = $2", aggregateID, version)
	return sr.scanRow(row)
}

func (er SnapshotRepository[U]) scanRow(row Row) (SnapshotModel[U], error) {
	snapshot := SnapshotModel[U]{}
	err := row.Scan(
		&snapshot.AggregateID,
		&snapshot.EventID,
		&snapshot.Version,
		&snapshot.Schema,
		&snapshot.Data,
	)
	return snapshot, err
}
