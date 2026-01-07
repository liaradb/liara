package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/liaradb/liaradb/collection/btree"
	key "github.com/liaradb/liaradb/collection/btree/value" // TODO: Fix this name
	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/manager"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

type Transaction struct {
	id             record.TransactionID
	lsn            record.LogSequenceNumber
	log            *recovery.Log
	bufferList     *BufferList
	concurrencyMgr *locktable.ConcurrencyMgr[action.ItemID]
	manager        *manager.Manager
	cursor         *btree.Cursor
	eventLog       *eventlog.EventLog
	keyValue       *keyvalue.KeyValue
	items          [][]byte
	forceRollback  bool
}

func newTransaction(
	id record.TransactionID,
	log *recovery.Log,
	bufferList *BufferList,
	concurrencyMgr *locktable.ConcurrencyMgr[action.ItemID],
	manager *manager.Manager,
	cursor *btree.Cursor,
	eventLog *eventlog.EventLog,
	keyValue *keyvalue.KeyValue,
) *Transaction {
	return &Transaction{
		id:             id,
		log:            log,
		bufferList:     bufferList,
		concurrencyMgr: concurrencyMgr,
		manager:        manager,
		cursor:         cursor,
		eventLog:       eventLog,
		keyValue:       keyValue,
	}
}

func (t *Transaction) ID() record.TransactionID                    { return t.id }
func (t *Transaction) LogSequenceNumber() record.LogSequenceNumber { return t.lsn }

func (t *Transaction) Insert(ctx context.Context, itemID action.ItemID, now time.Time, data []byte) error {
	if err := t.concurrencyMgr.XLock(ctx, itemID); err != nil {
		return err
	}

	lsn, err := t.log.Append(ctx, t.id, now, record.ActionInsert, data, nil)
	if err != nil {
		return err
	}

	t.items = append(t.items, data)

	t.lsn = lsn
	return nil
}

func (t *Transaction) Run(
	ctx context.Context,
	tn tablename.TableName, // TODO: How should this be specified?
	pid value.PartitionID,
	now time.Time, // TODO: How should this be specified?
	f func() error,
) error {
	if err := t.run(ctx, tn, pid, now, f); err != nil {
		return errTransactionFailed(t.id, err)
	}

	return nil
}

func (t *Transaction) run(
	ctx context.Context,
	tn tablename.TableName, // TODO: How should this be specified?
	pid value.PartitionID,
	now time.Time, // TODO: How should this be specified?
	f func() error,
) error {
	defer t.release()

	if err := f(); err != nil {
		return errors.Join(err, t.rollback(ctx, now))
	}

	if t.forceRollback {
		return t.rollback(ctx, now)
	}

	return t.commit(ctx, tn, pid, now)
}

func (t *Transaction) release() {
	t.concurrencyMgr.Release()
	t.bufferList.Release()
}

func (t *Transaction) commit(
	ctx context.Context,
	tn tablename.TableName, // TODO: How should this be specified?
	pid value.PartitionID,
	now time.Time,
) error {
	lsn, err := t.log.Append(ctx, t.id, now, record.ActionCommit, nil, nil)
	if err != nil {
		return err
	}

	if err := t.flush(ctx, lsn); err != nil {
		return err
	}

	if err := t.appendToEventLog(ctx, tn, pid); err != nil {
		return err
	}

	return nil
}

func (t *Transaction) appendToEventLog(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
) error {
	fn := tn.EventLog(pid)
	idxFn := tn.Index(0, pid)

	for _, item := range t.items {
		rid, err := t.eventLog.AppendEvent(ctx, fn, raw.NewBufferFromSlice(item))
		if err != nil {
			return err
		}

		// TODO: This needs an AggregateID
		if err := t.cursor.Insert(ctx, idxFn, key.NewKey(nil), rid); err != nil {
			return err
		}
	}

	return nil
}

func (t *Transaction) rollback(ctx context.Context, now time.Time) error {
	lsn, err := t.log.Append(ctx, t.id, now, record.ActionRollback, t.items[0], nil)
	if err != nil {
		return err
	}

	if err := t.flush(ctx, lsn); err != nil {
		return err
	}

	return nil
}

func (t *Transaction) flush(ctx context.Context, lsn record.LogSequenceNumber) error {
	t.lsn = lsn
	// TODO: Is this correct?
	return t.log.Flush(ctx, lsn)
}
