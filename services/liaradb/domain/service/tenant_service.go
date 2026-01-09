package service

import (
	"context"
	"encoding/json"
	"iter"

	key "github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type TenantService struct {
	outboxRepository  OutboxRepository
	requestRepository RequestRepository
	kv                *keyvalue.KeyValue
}

func NewTenantService(
	outboxRepository OutboxRepository,
	requestRepository RequestRepository,
	kv *keyvalue.KeyValue,
) *TenantService {
	return &TenantService{
		outboxRepository:  outboxRepository,
		requestRepository: requestRepository,
		kv:                kv,
	}
}

type CreateTenantCommand struct {
	TenantID   value.TenantID
	TenantName value.TenantName
}

type TenantModel struct {
	ID      string `json:"id"`
	Version int64  `json:"version"`
	Name    string `json:"name"`
}

// TODO: Create transaction
func (ts *TenantService) Create(ctx context.Context, cmd CreateTenantCommand) (value.TenantID, error) {
	tn := tablename.New("tenants")

	tnt := entity.NewTenant(cmd.TenantID, cmd.TenantName)
	data, err := json.Marshal(TenantModel{
		ID:      tnt.ID().String(),
		Version: int64(tnt.Version().Value()),
		Name:    tnt.Name().String(),
	})
	if err != nil {
		return "", err
	}

	return cmd.TenantID, ts.kv.Set(ctx, tn, key.NewKey([]byte(cmd.TenantID)), data)
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
	tn := tablename.New("tenants")

	k := key.NewKey([]byte(tenantID))
	data, err := ts.kv.Get(ctx, tn, k)
	if err != nil {
		return nil, err
	}

	m := TenantModel{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return entity.RestoreTenant(
		value.TenantID(m.ID),
		value.NewVersion(uint64(m.Version)),
		value.TenantName(m.Name),
	), nil
}

// TODO: Create transaction
func (ts *TenantService) List(ctx context.Context, limit int, offset int) iter.Seq2[*entity.Tenant, error] {
	return func(yield func(*entity.Tenant, error) bool) {
		tn := tablename.New("tenants")

		for data, err := range ts.kv.List(ctx, tn) {
			if err != nil {
				yield(nil, err)
				return
			}

			m := TenantModel{}
			if err := json.Unmarshal(data, &m); err != nil {
				yield(nil, err)
				return
			}

			if !yield(entity.RestoreTenant(
				value.TenantID(m.ID),
				value.NewVersion(uint64(m.Version)),
				value.TenantName(m.Name),
			), nil) {
				return
			}
		}
	}
}
