package service

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/outbox"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction"
)

type OutboxService struct {
	txManager *transaction.Manager
	outbox    *outbox.Outbox
}

func NewOutboxService(
	txManager *transaction.Manager,
	outbox *outbox.Outbox,
) *OutboxService {
	return &OutboxService{
		txManager: txManager,
		outbox:    outbox,
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

	tn := tablename.New(tid)
	oid := value.NewOutboxID()
	if err := tx.XLock(ctx, tn, oid.String()); err != nil {
		return value.OutboxID{}, err
	}

	// TODO: How do we handle a range?
	return tx.RunResult(ctx, func() (value.OutboxID, error) {
		outbox := entity.NewOutbox(oid, partitionRange)
		data := make([]byte, entity.OutboxSize)
		_ = outbox.WriteData(data)

		if err := os.outbox.Set(ctx, tx, tn, partitionRange.Low(), oid, outbox); err != nil {
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

	tn := tablename.New(tid)
	if err := tx.SLock(ctx, tn, oid.String()); err != nil {
		return nil, err
	}

	// TODO: What should PartitionID be?
	return os.outbox.Get(ctx, tn, value.NewPartitionID(0), oid)
}

func (os *OutboxService) UpdateOutbox(
	ctx context.Context,
	l span.Log,
	tn tablename.TableName,
	pid value.PartitionID,
	oid value.OutboxID,
	v value.GlobalVersion,
) error {
	tx, err := os.txManager.Next(ctx, tn.TenantID())
	if err != nil {
		return err
	}

	return tx.Run(ctx, func() error {
		if err := tx.XLock(ctx, tn, oid.String()); err != nil {
			return err
		}

		o, err := os.outbox.Get(ctx, tn, pid, oid)
		if err != nil {
			return err
		}

		prev := make([]byte, entity.OutboxSize)
		_ = o.WriteData(prev)

		o.UpdateGlobalVersion(v)

		data := make([]byte, entity.OutboxSize)
		_ = o.WriteData(data)

		_, err = tx.Update(ctx, link.RecordID{}, data, prev)
		if err != nil {
			return err
		}

		return os.outbox.Replace(ctx, l, tn, pid, oid, o)
	})
}

func (os *OutboxService) ListOutboxes(
	ctx context.Context,
	tid value.TenantID,
) iter.Seq2[*entity.Outbox, error] {
	// tx, err := es.txManager.Next(ctx, tid)
	// if err != nil {
	// 	return iterator.Error[*entity.Outbox](err)
	// }

	// TODO: How do we lock this query?
	// if err := t.concurrencyMgr.SLock(ctx, ItemID(oid.String())); err != nil {
	// 	return nil, err
	// }

	tn := tablename.New(tid)
	// TODO: What should PartitionID be?
	return os.outbox.List(ctx, tn, value.NewPartitionID(0))
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
		tn := tablename.New(tid)
		return os.UpdateOutbox(ctx, tx, tn, value.NewPartitionID(0), outboxID, globalVersion)
	})
}
