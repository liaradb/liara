package service

import (
	"context"

	"github.com/liaradb/liaradb/collection/fixed"
	"github.com/liaradb/liaradb/collection/idempotency"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/collection/tenant"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/transaction"
)

type baseService struct {
	txManager   *transaction.Manager
	tenant      *tenant.Tenant
	idempotency *idempotency.Idempotency
}

func newBaseService(
	txManager *transaction.Manager,
	tenant *tenant.Tenant,
	idempotency *idempotency.Idempotency,
) baseService {
	return baseService{
		txManager:   txManager,
		tenant:      tenant,
		idempotency: idempotency,
	}
}

func (s *baseService) insertRequestID(
	ctx context.Context,
	tx fixed.Transaction,
	tn tablename.TableName,
	rqid value.RequestID,
	now value.Time,
) error {
	return s.idempotency.Insert(ctx, tx, tn, rqid, entity.NewRequestLog(rqid, now))
}

func (s *baseService) testRequestID(
	ctx context.Context,
	tx fixed.Transaction,
	tn tablename.TableName,
	rqid value.RequestID,
) (bool, error) {
	return s.idempotency.Test(ctx, tx, tn, rqid)
}

func (s *baseService) initTransaction(ctx context.Context, tid value.TenantID) (*transaction.Transaction, error) {
	tx, err := s.txManager.Next(ctx)
	if err != nil {
		return nil, err
	}

	if tid == value.DefaultTenantID {
		return tx, nil
	}

	if _, err := s.tenant.Get(ctx, tx, tid); err != nil {
		tx.Release()
		return nil, err
	}

	return tx, nil
}
