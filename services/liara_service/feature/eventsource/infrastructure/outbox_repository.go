package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type (
	OutboxRepository struct {
		db   *sql.DB
		name string
	}

	outboxModel struct {
		ID       value.OutboxID
		Position value.GlobalVersion
	}
)

func NewOutboxRepository(db *sql.DB, name string) *OutboxRepository {
	return &OutboxRepository{
		db,
		name,
	}
}

func (s *OutboxRepository) GetOrCreateOutbox(
	ctx context.Context,
	outboxID value.OutboxID,
) (value.GlobalVersion, error) {
	outbox, err := s.getOutbox(ctx, outboxID)
	if errors.Is(err, sql.ErrNoRows) {
		err = s.createOutbox(ctx, outboxID)
	}
	return outbox.Position, err
}

func (s *OutboxRepository) getOutbox(
	ctx context.Context,
	outboxID value.OutboxID,
) (outboxModel, error) {
	row := s.db.QueryRowContext(ctx, fmt.Sprintf(`
SELECT * FROM %v
WHERE id = $1
`,
		s.name), outboxID)
	return s.scanRow(row)
}

func (s *OutboxRepository) createOutbox(
	ctx context.Context,
	outboxID value.OutboxID,
) error {
	_, err := s.db.ExecContext(ctx, fmt.Sprintf(`
INSERT INTO %v
VALUES( $1, $2 )
`,
		s.name), outboxID, 0)
	return err
}

func (s *OutboxRepository) UpdateOutboxPosition(
	ctx context.Context,
	outboxID value.OutboxID,
	position value.GlobalVersion,
) error {
	_, err := s.db.ExecContext(ctx, fmt.Sprintf(`
UPDATE %v
SET position = $2
WHERE id = $1
`,
		s.name), outboxID, position)
	return err
}

func (s *OutboxRepository) scanRow(row Row) (outboxModel, error) {
	outbox := outboxModel{}
	err := row.Scan(
		&outbox.ID,
		&outbox.Position,
	)
	return outbox, err
}

func (s *OutboxRepository) CreateTable(ctx context.Context) error {
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %v (
	id VARCHAR(50) PRIMARY KEY NOT NULL,
	position BIGINT NOT NULL
);
`,
		s.name)
	_, err := s.db.ExecContext(ctx, query)
	return err
}
