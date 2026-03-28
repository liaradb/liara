package service

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/collection/tenant"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type TenantService struct {
	tc *tenant.Tenant
}

func NewTenantService(
	tc *tenant.Tenant,
) *TenantService {
	return &TenantService{
		tc: tc,
	}
}

type CreateTenantCommand struct {
	TenantName value.TenantName
}

// TODO: Create transaction
func (ts *TenantService) Create(ctx context.Context, cmd CreateTenantCommand) (value.TenantID, error) {
	tid := value.NewTenantID()
	tnt := entity.NewTenant(tid, cmd.TenantName)

	if err := ts.tc.Set(ctx, tablename.Tenant, value.NewPartitionID(0), tid, tnt); err != nil {
		return value.TenantID{}, err
	}

	return tid, nil
	// id := cmd.TenantID.NewIfEmpty()
	// tenant := entity.NewTenant(id, cmd.TenantName)

	// if err := ts.transactionContainer.Run(ctx, func() error {
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
	// }); err != nil {
	// 	return "", err
	// }

	// return id, nil
}

type DeleteTenantCommand struct {
	TenantID value.TenantID
}

func (ts *TenantService) Delete(ctx context.Context, cmd DeleteTenantCommand) error {
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

type RenameTenantCommand struct {
	TenantID   value.TenantID
	TenantName value.TenantName
}

// TODO: Create transaction
func (ts *TenantService) Rename(ctx context.Context, cmd RenameTenantCommand) error {
	tnt, err := ts.tc.Get(ctx, tablename.Tenant, value.NewPartitionID(0), cmd.TenantID)
	if err != nil {
		return err
	}

	if err := tnt.Rename(cmd.TenantName); err != nil {
		return err
	}

	return ts.tc.Replace(ctx, tablename.Tenant, value.NewPartitionID(0), cmd.TenantID, tnt)
	// t, err := ts.tenantRepository.Get(ctx, cmd.TenantID)
	// if err != nil {
	// 	return err
	// }

	// if err := t.Rename(cmd.TenantName); err != nil {
	// 	return err
	// }

	// return ts.tenantRepository.Replace(ctx, t)
}

// TODO: Create transaction
func (ts *TenantService) Get(ctx context.Context, tenantID value.TenantID) (*entity.Tenant, error) {
	return ts.tc.Get(ctx, tablename.Tenant, value.NewPartitionID(0), tenantID)
}

// TODO: Create transaction
func (ts *TenantService) List(ctx context.Context, limit int, offset int) iter.Seq2[*entity.Tenant, error] {
	return ts.tc.List(ctx, tablename.Tenant, value.NewPartitionID(0))
}
