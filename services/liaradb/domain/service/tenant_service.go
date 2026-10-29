package service

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/collection/tenant"
	"github.com/liaradb/liaradb/domain/command"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/transaction"
)

type TenantService struct {
	txManager *transaction.Manager
	tc        *tenant.Tenant
}

func NewTenantService(
	txManager *transaction.Manager,
	tc *tenant.Tenant,
) *TenantService {
	return &TenantService{
		txManager: txManager,
		tc:        tc,
	}
}

func (ts *TenantService) Create(ctx context.Context, cmd command.CreateTenant) (value.TenantID, error) {
	tid := value.NewTenantID()
	tnt := entity.NewTenant(tid, cmd.TenantName)

	tx, err := ts.txManager.Next(ctx, tid)
	if err != nil {
		return value.TenantID{}, err
	}

	return tx.RunResult(ctx, func() (value.TenantID, error) {
		if err := ts.tc.Set(ctx, tx, tablename.Tenant, value.NewPartitionID(0), tid, tnt); err != nil {
			return value.TenantID{}, err
		}

		return tid, nil
	})
}

func (ts *TenantService) Delete(ctx context.Context, cmd command.DeleteTenant) error {
	panic("unimplemented")
}

func (ts *TenantService) Rename(ctx context.Context, cmd command.RenameTenant) error {
	tnt, err := ts.tc.Get(ctx, tablename.Tenant, value.NewPartitionID(0), cmd.TenantID)
	if err != nil {
		return err
	}

	if err := tnt.Rename(cmd.TenantName); err != nil {
		return err
	}

	tx, err := ts.txManager.Next(ctx, cmd.TenantID)
	if err != nil {
		return err
	}

	return tx.Run(ctx, func() error {
		return ts.tc.Replace(ctx, tx, tablename.Tenant, value.NewPartitionID(0), cmd.TenantID, tnt)
	})
}

func (ts *TenantService) Get(ctx context.Context, tenantID value.TenantID) (*entity.Tenant, error) {
	tx, err := ts.txManager.Next(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	defer tx.Release()

	if err := tx.SLock(ctx, tablename.New(tenantID), tenantID.String()); err != nil {
		return nil, err
	}

	return ts.tc.Get(ctx, tablename.Tenant, value.NewPartitionID(0), tenantID)
}

// TODO: Create transaction
func (ts *TenantService) List(ctx context.Context, limit int, offset int) iter.Seq2[*entity.Tenant, error] {
	return ts.tc.List(ctx, tablename.Tenant, value.NewPartitionID(0))
}
