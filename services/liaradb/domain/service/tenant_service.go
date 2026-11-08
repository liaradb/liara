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
	"github.com/liaradb/liaradb/util/iterator"
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

	tx, err := ts.txManager.Next(ctx)
	if err != nil {
		return value.TenantID{}, err
	}

	return tx.RunResult(ctx, func() (value.TenantID, error) {
		if err := ts.tc.Insert(ctx, tx, tnt); err != nil {
			return value.TenantID{}, err
		}

		return tid, nil
	})
}

func (ts *TenantService) Delete(ctx context.Context, cmd command.DeleteTenant) error {
	panic("unimplemented")
}

func (ts *TenantService) Rename(ctx context.Context, cmd command.RenameTenant) error {
	tx, err := ts.txManager.Next(ctx)
	if err != nil {
		return err
	}

	if err := tx.SLock(ctx, tablename.Tenant, cmd.TenantID.String()); err != nil {
		return err
	}

	tnt, err := ts.tc.Get(ctx, tx, cmd.TenantID)
	if err != nil {
		return err
	}

	if err := tnt.Rename(cmd.TenantName); err != nil {
		return err
	}

	return tx.Run(ctx, func() error {
		return ts.tc.Replace(ctx, tx, tnt)
	})
}

func (ts *TenantService) Get(ctx context.Context, tenantID value.TenantID) (*entity.Tenant, error) {
	tx, err := ts.txManager.Next(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Release()

	return ts.tc.Get(ctx, tx, tenantID)
}

func (ts *TenantService) List(ctx context.Context, limit int, offset int) iter.Seq2[*entity.Tenant, error] {
	tx, err := ts.txManager.Next(ctx)
	if err != nil {
		return iterator.Error[*entity.Tenant](err)
	}

	defer tx.Release()

	return ts.tc.List(ctx, tx)
}
