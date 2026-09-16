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
	"github.com/liaradb/liaradb/transaction/record"
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

// TODO: Create transaction
func (ts *TenantService) Create(ctx context.Context, cmd command.CreateTenant) (value.TenantID, error) {
	tid := value.NewTenantID()
	tnt := entity.NewTenant(tid, cmd.TenantName)

	tx, err := ts.txManager.Next(ctx, tid)
	if err != nil {
		return value.TenantID{}, err
	}

	// TODO: Should this log?
	lg := tx.Logger(ctx, record.CollectionOutbox)

	return transaction.RunResult(ctx, lg, tx, func() (value.TenantID, error) {
		// TODO: Use a logger
		if err := ts.tc.Set(ctx, lg, tablename.Tenant, value.NewPartitionID(0), tid, tnt); err != nil {
			return value.TenantID{}, err
		}
		return tid, nil

		// 	if err := ts.eventRepository.CreateTable(ctx, id); err != nil {
		// 		return err
		// 	}

		// 	if err := ts.eventRepository.CreateIndex(ctx, id); err != nil {
		// 		return err
		// 	}

		// 	if err := ts.outboxRepository.CreateTable(ctx, id); err != nil {
		// 		return err
		// 	}

		// 	if err := ts.requestRepository.CreateTable(ctx, id); err != nil {
		// 		return err
		// 	}

		// 	return ts.tenantRepository.Insert(ctx, tenant)
	})
}

func (ts *TenantService) Delete(ctx context.Context, cmd command.DeleteTenant) error {
	panic("unimplemented")
	// return ts.transactionContainer.Run(ctx, func() error {
	// 	if err := ts.eventRepository.DropTable(ctx, cmd.TenantID); err != nil {
	// 		return nil
	// 	}

	// 	if err := ts.outboxRepository.DropTable(ctx, cmd.TenantID); err != nil {
	// 		return nil
	// 	}

	// 	if err := ts.requestRepository.DropTable(ctx, cmd.TenantID); err != nil {
	// 		return nil
	// 	}

	// 	return ts.tenantRepository.Delete(ctx, cmd.TenantID)
	// })
}

// TODO: Create transaction
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

	// TODO: Should this log?
	lg := tx.Logger(ctx, record.CollectionOutbox)

	return transaction.Run(ctx, lg, tx, func() error {
		// TODO: Use a logger
		return ts.tc.Replace(ctx, lg, tablename.Tenant, value.NewPartitionID(0), cmd.TenantID, tnt)
	})
}

// TODO: Create transaction
func (ts *TenantService) Get(ctx context.Context, tenantID value.TenantID) (*entity.Tenant, error) {
	return ts.tc.Get(ctx, tablename.Tenant, value.NewPartitionID(0), tenantID)
}

// TODO: Create transaction
func (ts *TenantService) List(ctx context.Context, limit int, offset int) iter.Seq2[*entity.Tenant, error] {
	return ts.tc.List(ctx, tablename.Tenant, value.NewPartitionID(0))
}
