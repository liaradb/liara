package transaction

import (
	"context"
	"time"

	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/log"
	"github.com/liaradb/liaradb/log/action"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/storage"
)

type Transaction struct {
	id             record.TransactionID
	lsn            record.LogSequenceNumber
	log            *log.Log
	bufferList     *storage.BufferList
	concurrencyMgr *locktable.ConcurrencyMgr[action.ItemID]
}

func newTransaction(
	id record.TransactionID,
	log *log.Log,
	bufferList *storage.BufferList,
	concurrencyMgr *locktable.ConcurrencyMgr[action.ItemID],
) *Transaction {
	return &Transaction{
		id:             id,
		log:            log,
		bufferList:     bufferList,
		concurrencyMgr: concurrencyMgr,
	}
}

func (t Transaction) ID() record.TransactionID                     { return t.id }
func (t *Transaction) LogSequenceNumber() record.LogSequenceNumber { return t.lsn }

func (t *Transaction) Insert(ctx context.Context, itemID action.ItemID, now time.Time, data []byte) error {
	if err := t.concurrencyMgr.XLock(ctx, itemID); err != nil {
		return errTransactionFailed(t.id, err)
	}

	lsn, err := t.log.Append(ctx, t.id, now, record.ActionInsert, data, nil)
	if err != nil {
		return errTransactionFailed(t.id, err)
	}

	t.lsn = lsn
	return nil
}

func (t *Transaction) Commit(ctx context.Context, now time.Time) error {
	lsn, err := t.log.Append(ctx, t.id, now, record.ActionCommit, nil, nil)
	if err != nil {
		return errTransactionFailed(t.id, err)
	}

	t.lsn = lsn
	// TODO: Is this correct?
	if err := t.log.Flush(ctx, lsn); err != nil {
		return errTransactionFailed(t.id, err)
	}

	t.concurrencyMgr.Release()
	t.bufferList.Release()

	return nil
}
