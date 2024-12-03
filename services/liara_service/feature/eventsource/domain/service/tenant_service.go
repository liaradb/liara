package service

import (
	"context"
	"iter"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
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

func (ts *TenantService) Insert(ctx context.Context, tenant *entity.Tenant) error {
	return ts.transactionRepository.Run(ctx, func(tx Transaction) error {
		tenantID := tenant.ID()

		if err := ts.eventRepository.CreateTable(ctx, tenantID); err != nil {
			return err
		}

		if err := ts.eventRepository.CreateIndex(ctx, tenantID); err != nil {
			return err
		}

		if err := ts.outboxRepository.CreateTable(ctx, tenantID); err != nil {
			return err
		}

		if err := ts.requestRepository.CreateTable(ctx, tenantID); err != nil {
			return err
		}

		return ts.tenantRepository.Insert(ctx, tenant)
	})
}

func (ts *TenantService) List(ctx context.Context, limit int, offset int) iter.Seq2[*entity.Tenant, error] {
	return ts.tenantRepository.List(ctx, limit, offset)
}
