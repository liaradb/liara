package service

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/idempotency"
	"github.com/liaradb/liaradb/collection/outbox"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/collection/tenant"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/transaction"
	"github.com/liaradb/liaradb/util/iterator"
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
	partitionRange value.PartitionRange,
) (value.OutboxID, error) {
	tx, err := os.initTransaction(ctx, tid)
	if err != nil {
		return value.OutboxID{}, err
	}

	tn := tablename.NewTenantFolder(tid).Outbox()
	oid := value.NewOutboxID()
	if err := tx.XLock(ctx, tn, oid.String()); err != nil {
		return value.OutboxID{}, err
	}

	// TODO: How do we handle a range?
	return tx.RunResult(ctx, func() (value.OutboxID, error) {
		outbox := entity.NewOutbox(oid, partitionRange)
		data := make([]byte, entity.OutboxSize)
		_ = outbox.WriteData(data)

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

		o.UpdateGlobalVersion(v)

		return os.outbox.Replace(ctx, tx, tid, oid, o)
	})
}

func (os *OutboxService) ListOutboxes(
	ctx context.Context,
	tid value.TenantID,
) iter.Seq2[*entity.Outbox, error] {
	tx, err := os.initTransaction(ctx, tid)
	if err != nil {
		return iterator.Error[*entity.Outbox](err)
	}

	return os.outbox.List(ctx, tx, tid)
}
