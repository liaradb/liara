package service

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/infrastructure/idempotency"
	"github.com/liaradb/liaradb/domain/infrastructure/outbox"
	"github.com/liaradb/liaradb/domain/infrastructure/tenant"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/transaction"
)

type OutboxService struct {
	baseService
	txManager *transaction.Manager
	outbox    *outbox.Outbox
}

func NewOutboxService(
	txManager *transaction.Manager,
	tenant *tenant.Tenant,
	idempotency *idempotency.Idempotency,
	outbox *outbox.Outbox,
) *OutboxService {
	return &OutboxService{
		baseService: newBaseService(txManager, tenant, idempotency),
		txManager:   txManager,
		outbox:      outbox,
	}
}

func (os *OutboxService) CreateOutbox(
	ctx context.Context,
	tid value.TenantID,
	pid value.PartitionID,
) (value.OutboxID, error) {
	tx, err := os.initTransaction(ctx, tid)
	if err != nil {
		return value.OutboxID{}, err
	}

	return tx.RunResult(ctx, func() (value.OutboxID, error) {
		oid := value.NewOutboxID()
		outbox := entity.NewOutbox(oid, pid)

		if err := os.outbox.Insert(ctx, tx, tid, oid, outbox); err != nil {
			return value.OutboxID{}, err
		}

		return oid, nil
	})
}

func (os *OutboxService) GetOutbox(
	ctx context.Context,
	tid value.TenantID,
	oid value.OutboxID,
) (*entity.Outbox, error) {
	tx, err := os.initTransaction(ctx, tid)
	if err != nil {
		return nil, err
	}

	defer tx.Release()

	return os.outbox.Get(ctx, tx, tid, oid)
}

func (os *OutboxService) UpdateOutboxPosition(
	ctx context.Context,
	tid value.TenantID,
	oid value.OutboxID,
	v value.GlobalVersion,
	sv value.SubVersion,
) error {
	tx, err := os.initTransaction(ctx, tid)
	if err != nil {
		return err
	}

	return tx.Run(ctx, func() error {
		o, err := os.outbox.Get(ctx, tx, tid, oid)
		if err != nil {
			return err
		}

		o.UpdatePosition(v, sv)

		return os.outbox.Replace(ctx, tx, tid, oid, o)
	})
}

func (os *OutboxService) ListOutboxes(
	ctx context.Context,
	tid value.TenantID,
) iter.Seq2[*entity.Outbox, error] {
	return func(yield func(*entity.Outbox, error) bool) {
		tx, err := os.initTransaction(ctx, tid)
		if err != nil {
			yield(nil, err)
			return
		}

		defer tx.Release()

		for o, err := range os.outbox.List(ctx, tx, tid) {
			if !yield(o, err) {
				return
			}
		}
	}
}
