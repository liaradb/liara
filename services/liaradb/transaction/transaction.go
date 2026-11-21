package transaction

import (
	"context"
	"errors"

	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/locktable"
	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/transaction/record"
)

type Transaction struct {
	id             record.TransactionID
	log            *log.Log
	bufferList     *BufferList
	concurrencyMgr *locktable.ConcurrencyMgr[ItemID]
	forceRollback  bool
	manager        *Manager
}

func newTransaction(
	id record.TransactionID,
	log *log.Log,
	bufferList *BufferList,
	concurrencyMgr *locktable.ConcurrencyMgr[ItemID],
	manager *Manager,
) *Transaction {
	return &Transaction{
		id:             id,
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

func (t *Transaction) Append(
	ctx context.Context,
	frid link.FileRecordID,
	data []byte,
) (logpage.LogSequenceNumber, error) {
	return t.log.Append(ctx,
		t.id,
		frid,
		data)
}

func (t *Transaction) Insert(
	ctx context.Context,
	frid link.FileRecordID,
	data []byte,
) (logpage.LogSequenceNumber, error) {
	return t.log.Insert(ctx,
		t.id,
		frid,
		data)
}

func (t *Transaction) Update(
	ctx context.Context,
	frid link.FileRecordID,
	data []byte,
	prev []byte,
) (logpage.LogSequenceNumber, error) {
	return t.log.Update(ctx,
		t.id,
		frid,
		data,
		prev)
}

func (t *Transaction) Clear(
	ctx context.Context,
	frid link.FileRecordID,
	prev []byte,
) (logpage.LogSequenceNumber, error) {
	return t.log.Clear(ctx,
		t.id,
		frid,
		prev)
}

func (t *Transaction) UpdateHeader(
	ctx context.Context,
	frid link.FileRecordID,
	data []byte,
	prev []byte,
) (logpage.LogSequenceNumber, error) {
	return t.log.UpdateHeader(ctx,
		t.id,
		frid,
		data,
		prev)
}

func (t *Transaction) Pin(ctx context.Context, fbid link.FileBlockID) (*storage.Buffer, error) {
	return t.bufferList.Pin(ctx, fbid)
}

func (t *Transaction) PinCurrent(ctx context.Context, fn link.FileName) (*storage.Buffer, error) {
	return t.bufferList.PinCurrent(ctx, fn)
}

func (t *Transaction) PinNext(ctx context.Context, fn link.FileName) (*storage.Buffer, error) {
	return t.bufferList.PinNext(ctx, fn)
}

func (t *Transaction) Unpin(b *storage.Buffer) {
	t.bufferList.ReleaseBuffer(b)
}
