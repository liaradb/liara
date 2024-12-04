package service

import (
	"context"
	"iter"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type TenantService struct {
	transactionRepository TransactionRepository
	eventRepository       EventRepository
	outboxRepository      OutboxRepository
	requestRepository     RequestRepository
	tenantRepository      TenantRepository
}

type TenantRepository interface {
	Insert(context.Context, *entity.Tenant) error
	Replace(context.Context, *entity.Tenant) error
	Delete(context.Context, value.TenantID) error
	Get(context.Context, value.TenantID) (*entity.Tenant, error)
	List(context.Context, int, int) iter.Seq2[*entity.Tenant, error]
}

func NewTenantService(
	transactionRepository TransactionRepository,
	eventRepository EventRepository,
	outboxRepository OutboxRepository,
	requestRepository RequestRepository,
	tenantRepository TenantRepository,
) *TenantService {
	return &TenantService{
		transactionRepository: transactionRepository,
		eventRepository:       eventRepository,
		outboxRepository:      outboxRepository,
		requestRepository:     requestRepository,
		tenantRepository:      tenantRepository,
	}
}

type CreateTenantCommand struct {
	TenantName value.TenantName
}

func (ts *TenantService) Create(ctx context.Context, cmd CreateTenantCommand) (value.TenantID, error) {
	id := value.NewTenantID()
	tenant := entity.NewTenant(id, cmd.TenantName)

	if err := ts.transactionRepository.Run(ctx, func(tx Transaction) error {
		if err := ts.eventRepository.CreateTable(ctx, id); err != nil {
			return err
		}

		if err := ts.eventRepository.CreateIndex(ctx, id); err != nil {
			return err
		}

		if err := ts.outboxRepository.CreateTable(ctx, id); err != nil {
			return err
		}

		if err := ts.requestRepository.CreateTable(ctx, id); err != nil {
			return err
		}

		return ts.tenantRepository.Insert(ctx, tenant)
	}); err != nil {
		return "", err
	}

	return id, nil
}

type DeleteTenantCommand struct {
	TenantID value.TenantID
}

func (ts *TenantService) Delete(ctx context.Context, cmd DeleteTenantCommand) error {
	return ts.transactionRepository.Run(ctx, func(tx Transaction) error {
		if err := ts.eventRepository.DropTable(ctx, cmd.TenantID); err != nil {
			return nil
		}

		if err := ts.outboxRepository.DropTable(ctx, cmd.TenantID); err != nil {
			return nil
		}

		if err := ts.requestRepository.DropTable(ctx, cmd.TenantID); err != nil {
			return nil
		}

		return ts.tenantRepository.Delete(ctx, cmd.TenantID)
	})
}

type RenameTenantCommand struct {
	TenantID   value.TenantID
	TenantName value.TenantName
}

func (ts *TenantService) Rename(ctx context.Context, cmd RenameTenantCommand) error {
	t, err := ts.tenantRepository.Get(ctx, cmd.TenantID)
	if err != nil {
		return err
	}

	if err := t.Rename(cmd.TenantName); err != nil {
		return err
	}

	return ts.tenantRepository.Replace(ctx, t)
}

func (ts *TenantService) Get(ctx context.Context, tenantID value.TenantID) (*entity.Tenant, error) {
	return ts.tenantRepository.Get(ctx, tenantID)
}

func (ts *TenantService) List(ctx context.Context, limit int, offset int) iter.Seq2[*entity.Tenant, error] {
	return ts.tenantRepository.List(ctx, limit, offset)
}
