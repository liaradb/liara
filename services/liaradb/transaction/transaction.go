package transaction

import (
	"context"
	"time"

	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/manager"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/storage/link"
)

type Transaction struct {
	id             record.TransactionID
	lsn            record.LogSequenceNumber
	log            *recovery.Log
	bufferList     *BufferList
	concurrencyMgr *locktable.ConcurrencyMgr[action.ItemID]
	manager        *manager.Manager
	eventLog       *eventlog.EventLog
	keyValue       *keyvalue.KeyValue
	items          [][]byte
}

func newTransaction(
	id record.TransactionID,
	log *recovery.Log,
	bufferList *BufferList,
	concurrencyMgr *locktable.ConcurrencyMgr[action.ItemID],
	manager *manager.Manager,
	eventLog *eventlog.EventLog,
	keyValue *keyvalue.KeyValue,
) *Transaction {
	return &Transaction{
		id:             id,
		log:            log,
		bufferList:     bufferList,
		concurrencyMgr: concurrencyMgr,
		manager:        manager,
		eventLog:       eventLog,
		keyValue:       keyValue,
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

	t.items = append(t.items, data)

	t.lsn = lsn
	return nil
}

func (t *Transaction) Commit(
	ctx context.Context,
	fn link.FileName, // TODO: How should this be specified?
	now time.Time,
) error {
	lsn, err := t.log.Append(ctx, t.id, now, record.ActionCommit, t.items[0], nil)
	if err != nil {
		return errTransactionFailed(t.id, err)
	}

	t.lsn = lsn
	// TODO: Is this correct?
	if err := t.log.Flush(ctx, lsn); err != nil {
		return errTransactionFailed(t.id, err)
	}

	if _, err := t.eventLog.AppendEvent(ctx, fn, raw.NewBufferFromSlice(t.items[0])); err != nil {
		return errTransactionFailed(t.id, err)
	}

	t.concurrencyMgr.Release()
	t.bufferList.Release()

	return nil
}

func (t *Transaction) Rollback(ctx context.Context, now time.Time) error {
	lsn, err := t.log.Append(ctx, t.id, now, record.ActionRollback, t.items[0], nil)
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
