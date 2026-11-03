package service

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/idempotency"
	"github.com/liaradb/liaradb/collection/outbox"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/storage/link"
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
	idempotency *idempotency.Idempotency,
	outbox *outbox.Outbox,
) *OutboxService {
	return &OutboxService{
		baseService: newBaseService(txManager, idempotency),
		txManager:   txManager,
		outbox:      outbox,
	}
}

func (os *OutboxService) CreateOutbox(
	ctx context.Context,
	tid value.TenantID,
	partitionRange value.PartitionRange,
) (value.OutboxID, error) {
	tx, err := os.txManager.Next(ctx, tid)
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

		if err := os.outbox.Insert(ctx, tx, tn, oid, outbox); err != nil {
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
	tx, err := os.txManager.Next(ctx, tid)
	if err != nil {
		return nil, err
	}

	defer tx.Release()

	tn := tablename.NewTenantFolder(tid).Outbox()
	return os.outbox.Get(ctx, tx, tn, oid)
}

func (os *OutboxService) UpdateOutbox(
	ctx context.Context,
	l span.Log,
	tid value.TenantID,
	oid value.OutboxID,
	v value.GlobalVersion,
) error {
	tx, err := os.txManager.Next(ctx, tid)
	if err != nil {
		return err
	}

	tn := tablename.NewTenantFolder(tid).Outbox()
	return tx.Run(ctx, func() error {
		o, err := os.outbox.Get(ctx, tx, tn, oid)
		if err != nil {
			return err
		}

		prev := make([]byte, entity.OutboxSize)
		_ = o.WriteData(prev)

		o.UpdateGlobalVersion(v)

		data := make([]byte, entity.OutboxSize)
		_ = o.WriteData(data)

		_, err = tx.Update(ctx, link.FileRecordID{}, data, prev)
		if err != nil {
			return err
		}

		return os.outbox.Replace(ctx, l, tn, oid, o)
	})
}

func (os *OutboxService) ListOutboxes(
	ctx context.Context,
	tid value.TenantID,
) iter.Seq2[*entity.Outbox, error] {
	tx, err := os.txManager.Next(ctx, tid)
	if err != nil {
		return iterator.Error[*entity.Outbox](err)
	}

	// TODO: How do we lock this query?
	// if err := t.concurrencyMgr.SLock(ctx, ItemID(oid.String())); err != nil {
	// 	return nil, err
	// }

	tn := tablename.NewTenantFolder(tid).Outbox()
	return os.outbox.List(ctx, tx, tn)
}

func (os *OutboxService) UpdateOutboxPosition(
	ctx context.Context,
	tid value.TenantID,
	outboxID value.OutboxID,
	globalVersion value.GlobalVersion,
) error {
	tx, err := os.txManager.Next(ctx, tid)
	if err != nil {
		return err
	}

	return tx.Run(ctx, func() error {
		return os.UpdateOutbox(ctx, tx, tid, outboxID, globalVersion)
	})
}
