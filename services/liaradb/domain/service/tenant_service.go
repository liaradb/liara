package service

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type TenantService struct {
	kv *keyvalue.KeyValue
}

func NewTenantService(
	kv *keyvalue.KeyValue,
) *TenantService {
	return &TenantService{
		kv: kv,
	}
}

type CreateTenantCommand struct {
	TenantName value.TenantName
}

// TODO: Create transaction
func (ts *TenantService) Create(ctx context.Context, cmd CreateTenantCommand) (value.TenantID, error) {
	tid := value.NewTenantID()
	tnt := entity.NewTenant(tid, cmd.TenantName)

	data := make([]byte, entity.TenantSize)
	_ = tnt.Write(data)

	if err := ts.kv.Set(ctx, tablename.Tenant, key.NewKey(tid.Bytes()), data); err != nil {
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

func (ts *TenantService) Rename(ctx context.Context, cmd RenameTenantCommand) error {
	panic("unimplemented")
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
	k := key.NewKey(tenantID.Bytes())
	data, err := ts.kv.Get(ctx, tablename.Tenant, k)
	if err != nil {
		return nil, err
	}

	tnt := &entity.Tenant{}
	_ = tnt.Read(data)

	return tnt, nil
}

// TODO: Create transaction
func (ts *TenantService) List(ctx context.Context, limit int, offset int) iter.Seq2[*entity.Tenant, error] {
	return func(yield func(*entity.Tenant, error) bool) {
		for data, err := range ts.kv.List(ctx, tablename.Tenant) {
			if err != nil {
				yield(nil, err)
				return
			}

			tnt := &entity.Tenant{}
			_ = tnt.Read(data)
			if !yield(tnt, nil) {
				return
			}
		}
	}
}
