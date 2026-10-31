package transaction

import (
	"context"
	"errors"

	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/locktable"
	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/transaction/record"
)

type Transaction struct {
	id             record.TransactionID
	tid            value.TenantID
	log            *log.Log
	bufferList     *BufferList
	concurrencyMgr *locktable.ConcurrencyMgr[ItemID]
	forceRollback  bool
	manager        *Manager
}

func newTransaction(
	id record.TransactionID,
	tid value.TenantID,
	log *log.Log,
	bufferList *BufferList,
	concurrencyMgr *locktable.ConcurrencyMgr[ItemID],
	manager *Manager,
) *Transaction {
	return &Transaction{
		id:             id,
		tid:            tid,
		log:            log,
		bufferList:     bufferList,
		concurrencyMgr: concurrencyMgr,
		manager:        manager,
	}
}

func (t *Transaction) ID() record.TransactionID { return t.id }

func (t *Transaction) XLock(
	ctx context.Context,
	tn tablename.TableName,
	id string,
) error {
	return t.concurrencyMgr.XLock(ctx, newItemID(tn, id))
}

func (t *Transaction) SLock(
	ctx context.Context,
	tn tablename.TableName,
	id string,
) error {
	return t.concurrencyMgr.SLock(ctx, newItemID(tn, id))
}

func (t *Transaction) Run(
	ctx context.Context,
	f func() error,
) error {
	_, err := t.run(ctx, func() (any, error) {
		return struct{}{}, f()
	})
	if err != nil {
		return errTransactionFailed(t.id, err)
	}

	return nil
}

func (t *Transaction) RunResult[R any](
	ctx context.Context,
	f func() (R, error),
) (R, error) {
	r, err := t.run(ctx, func() (R, error) {
		return f()
	})
	if err != nil {
		var r R
		return r, errTransactionFailed(t.id, err)
	}

	return r, nil
}

func (t *Transaction) run[R any](
	ctx context.Context,
	f func() (R, error),
) (R, error) {
	_, err := t.log.Start(ctx, t.id)
	if err != nil {
		var r R
		return r, err
	}

	defer t.Release()

	r, err := f()
	if err != nil {
		var r R
		return r, errors.Join(err, t.rollback(ctx))
	}

	if t.forceRollback {
		var r R
		return r, t.rollback(ctx)
	}

	return r, t.commit(ctx)
}

func (t *Transaction) Release() {
	t.concurrencyMgr.Release()
	t.bufferList.Release()
	t.manager.End(t.id)
}

func (t *Transaction) commit(ctx context.Context) error {
	_, err := t.log.Commit(ctx, t.id)
	return err
}

func (t *Transaction) rollback(ctx context.Context) error {
	_, err := t.log.Rollback(ctx, t.id)
	if err != nil {
		return err
	}

	return nil
}

func (l *Transaction) Append(
	ctx context.Context,
	rid link.RecordID,
	data []byte,
) (logpage.LogSequenceNumber, error) {
	return l.log.Append(ctx,
		l.id,
		rid,
		data)
}

func (l *Transaction) Insert(
	ctx context.Context,
	rid link.RecordID,
	data []byte,
) (logpage.LogSequenceNumber, error) {
	return l.log.Insert(ctx,
		l.id,
		rid,
		data)
}

func (l *Transaction) Update(
	ctx context.Context,
	rid link.RecordID,
	data []byte,
	prev []byte,
) (logpage.LogSequenceNumber, error) {
	return l.log.Update(ctx,
		l.id,
		rid,
		data,
		prev)
}

func (l *Transaction) Clear(
	ctx context.Context,
	rid link.RecordID,
	prev []byte,
) (logpage.LogSequenceNumber, error) {
	return l.log.Clear(ctx,
		l.id,
		rid,
		prev)
}

func (l *Transaction) UpdateHeader(
	ctx context.Context,
	rid link.RecordID,
	data []byte,
	prev []byte,
) (logpage.LogSequenceNumber, error) {
	return l.log.UpdateHeader(ctx,
		l.id,
		rid,
		data,
		prev)
}
