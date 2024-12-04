package infrastructure

import (
	"context"
	"database/sql"
	"iter"

	"github.com/cardboardrobots/baseerror"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/service"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type TenantRepository struct {
	db *sql.DB
}

var _ service.TenantRepository = (*TenantRepository)(nil)

func NewTenantRepository(
	db *sql.DB,
) *TenantRepository {
	return &TenantRepository{
		db: db,
	}
}

func (tr *TenantRepository) Insert(ctx context.Context, t *entity.Tenant) error {
	m := tenantToModel(t)
	_, err := tr.db.ExecContext(ctx, `
INSERT INTO tenants
VALUES( $1, $2, $3 )
`,
		m.ID,
		m.Version,
		m.Name)
	return err
}

func (tr *TenantRepository) Replace(ctx context.Context, t *entity.Tenant) error {
	v := t.Version()
	m := tenantToModel(t)
	m.increment()
	_, err := tr.db.ExecContext(ctx, `
UPDATE tenants
SET
  version = $3,
  name = $4
WHERE id = $1
AND version = $2
`,
		m.ID,
		v,
		m.Version,
		m.Name)
	return err
}

func (tr *TenantRepository) Delete(ctx context.Context, tenantID value.TenantID) error {
	_, err := tr.db.ExecContext(ctx, `
DELETE FROM tenants
WHERE id = $1
`,
		tenantID)
	return err
}

func (tr *TenantRepository) Get(ctx context.Context, tenantID value.TenantID) (*entity.Tenant, error) {
	row := tr.db.QueryRowContext(ctx, `
SELECT * FROM %v
WHERE id = $1
`,
		tenantID)
	t, err := tr.scanRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			err = baseerror.ErrNotFound
		}
		return nil, err
	}

	return t, nil
}

func (tr *TenantRepository) List(ctx context.Context, limit int, offset int) iter.Seq2[*entity.Tenant, error] {
	return func(yield func(*entity.Tenant, error) bool) {
		rows, err := tr.db.QueryContext(ctx, `
SELECT * FROM tenants
`,
			limit, offset)
		if err != nil {
			yield(nil, err)
			return
		}

		defer func() { err = rows.Close() }()
		for rows.Next() {
			t, err := tr.scanRow(rows)
			if !yield(t, err) {
				return
			}
		}
	}
}

func (tr *TenantRepository) scanRow(row Row) (*entity.Tenant, error) {
	m := tenantModel{}
	err := row.Scan(
		&m.ID,
		&m.Version,
		&m.Name,
	)
	return tenantFromModel(m), err
}

func (tr *TenantRepository) CreateTable(ctx context.Context) error {
	query := `
CREATE TABLE IF NOT EXISTS tenants (
	id VARCHAR(50) PRIMARY KEY NOT NULL,
	version BIGINT NOT NULL,
	name VARCHAR(128) NOT NULL
);
`
	_, err := tr.db.ExecContext(ctx, query)
	return err
}

type tenantModel struct {
	ID      value.TenantID
	Version value.Version
	Name    value.TenantName
}

func tenantToModel(t *entity.Tenant) tenantModel {
	return tenantModel{
		ID:      t.ID(),
		Version: t.Version(),
		Name:    t.Name(),
	}
}

func tenantFromModel(m tenantModel) *entity.Tenant {
	return entity.RestoreTenant(
		m.ID,
		m.Version,
		m.Name,
	)
}

func (m *tenantModel) increment() {
	m.Version++
}
