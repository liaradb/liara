package service

import (
	"context"

	"github.com/liaradb/liaradb/collection/idempotency"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/transaction"
)

type baseService struct {
	txManager   *transaction.Manager
	idempotency *idempotency.Idempotency
}

func newBaseService(
	txManager *transaction.Manager,
	idempotency *idempotency.Idempotency,
) baseService {
	return baseService{
		txManager:   txManager,
		idempotency: idempotency,
	}
}

func (s *baseService) insertRequestID(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
	rqid value.RequestID,
	now value.Time,
) error {
	return s.idempotency.Insert(ctx, tx, tn, rqid, entity.NewRequestLog(rqid, now))
}

func (s *baseService) testRequestID(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
	rqid value.RequestID,
) (bool, error) {
	return s.idempotency.Test(ctx, tx, tn, rqid)
}
