package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/cardboardrobots/baseerror"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type (
	OutboxRepository struct {
		db *sql.DB
	}

	outboxModel struct {
		ID             value.OutboxID
		PartitionRange value.PartitionRange
		Position       value.GlobalVersion
	}
)

func NewOutboxRepository(
	db *sql.DB,
) *OutboxRepository {
	return &OutboxRepository{
		db,
	}
}

func (*OutboxRepository) getName(tenantID value.TenantID) string {
	n := tenantID.String()
	if n == "" {
		n = "default"
	}
	n = strings.ReplaceAll(n, "-", "_")
	return fmt.Sprintf("__%v__outboxes", n)
}

func (s *OutboxRepository) GetOutbox(
	ctx context.Context,
	tenantID value.TenantID,
	outboxID value.OutboxID,
) (*entity.Outbox, error) {
	row := s.db.QueryRowContext(ctx, fmt.Sprintf(`
SELECT * FROM %v
WHERE id = $1
`,
		s.getName(tenantID)), outboxID)
	m, err := s.scanRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			err = baseerror.ErrNotFound
		}
		return nil, err
	}

	return entity.RestoreOutbox(m.ID, m.PartitionRange, m.Position), nil
}

func (s *OutboxRepository) CreateOutbox(
	ctx context.Context,
	tenantID value.TenantID,
	outbox *entity.Outbox,
) error {
	low, high := outbox.PartitionRange().All()
	_, err := s.db.ExecContext(ctx, fmt.Sprintf(`
INSERT INTO %v
VALUES( $1, $2, $3, $4 )
`,
		s.getName(tenantID)), outbox.ID(), low, high, 0)
	return err
}

func (s *OutboxRepository) UpdateOutboxPosition(
	ctx context.Context,
	tenantID value.TenantID,
	outboxID value.OutboxID,
	position value.GlobalVersion,
) error {
	_, err := s.db.ExecContext(ctx, fmt.Sprintf(`
UPDATE %v
SET position = $2
WHERE id = $1
`,
		s.getName(tenantID)), outboxID, position)
	return err
}

func (s *OutboxRepository) scanRow(row Row) (outboxModel, error) {
	outbox := outboxModel{}
	var low, high value.PartitionID
	err := row.Scan(
		&outbox.ID,
		&low,
		&high,
		&outbox.Position,
	)
	outbox.PartitionRange = value.NewPartitionRange(low, high)
	return outbox, err
}

func (s *OutboxRepository) CreateTable(ctx context.Context, tenantID value.TenantID) error {
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %v (
	id VARCHAR(50) PRIMARY KEY NOT NULL,
	partition_low INT NOT NULL,
	partition_high INT NOT NULL,
	position BIGINT NOT NULL
);
`,
		s.getName(tenantID))
	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (or *OutboxRepository) DropTable(ctx context.Context, tenantID value.TenantID) error {
	query := fmt.Sprintf(`
DROP TABLE %v;
`,
		or.getName(tenantID))
	_, err := or.db.ExecContext(ctx, query)
	return err
}
