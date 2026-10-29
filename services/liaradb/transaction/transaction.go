package transaction

import (
	"context"
	"errors"
	"iter"

	"github.com/liaradb/liaradb/collection"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
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
	collection     *collection.Collections
	forceRollback  bool
	manager        *Manager
}

func newTransaction(
	id record.TransactionID,
	tid value.TenantID,
	log *log.Log,
	bufferList *BufferList,
	concurrencyMgr *locktable.ConcurrencyMgr[ItemID],
	collection *collection.Collections,
	manager *Manager,
) *Transaction {
	return &Transaction{
		id:             id,
		tid:            tid,
		log:            log,
		bufferList:     bufferList,
		concurrencyMgr: concurrencyMgr,
		collection:     collection,
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

func (t *Transaction) GetAggregate(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	id value.AggregateID,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		if err := t.SLock(ctx, tn, id.String()); err != nil {
			yield(nil, err)
			return
		}

		defer t.release()

		for e, err := range t.collection.EventLog.GetAggregate(ctx, tn, pid, id) {
			if !yield(e, err) {
				return
			}
		}
	}
}

func (t *Transaction) GetAggregateByIDAndName(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	id value.AggregateID,
	name value.AggregateName,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		if err := t.SLock(ctx, tn, id.String()); err != nil {
			yield(nil, err)
			return
		}

		defer t.release()

		for e, err := range t.collection.EventLog.GetAggregate(ctx, tn, pid, id) {
			if err != nil {
				yield(nil, err)
				return
			}

			if e.AggregateName == name && !yield(e, nil) {
				return
			}
		}
	}
}

func (t *Transaction) Events(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		// Should we lock?
		// if err := t.SLock(ctx, tn, id.String()); err != nil {
		// 	yield(nil, err)
		// 	return
		// }

		// defer t.release()

		for e, err := range t.collection.EventLog.Events(ctx, tn, pid) {
			if !yield(e, err) {
				return
			}
		}
	}
}

func (t *Transaction) EventsAfterGlobalVersion(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	version value.GlobalVersion,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		// Should we lock?
		// if err := t.SLock(ctx, tn, id.String()); err != nil {
		// 	yield(nil, err)
		// 	return
		// }

		// defer t.release()

		for e, err := range t.collection.EventLog.EventsAfterGlobalVersion(ctx, tn, pid, version) {
			if !yield(e, err) {
				return
			}
		}
	}
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

	defer t.release()

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

func (t *Transaction) release() {
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

func (t *Transaction) GetOutbox(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	oid value.OutboxID,
) (*entity.Outbox, error) {
	if err := t.SLock(ctx, tn, oid.String()); err != nil {
		return nil, err
	}

	return t.collection.Outbox.Get(ctx, tn, pid, oid)
}

func (t *Transaction) InsertOutbox(
	ctx context.Context,
	l span.Log,
	tn tablename.TableName,
	pid value.PartitionID,
	oid value.OutboxID,
	e *entity.Outbox,
) error {
	if err := t.XLock(ctx, tn, oid.String()); err != nil {
		return err
	}

	data := make([]byte, entity.OutboxSize)
	_ = e.WriteData(data)
	return t.collection.Outbox.Set(ctx, l, tn, pid, oid, e)
}

func (t *Transaction) UpdateOutbox(
	ctx context.Context,
	l span.Log,
	tn tablename.TableName,
	pid value.PartitionID,
	oid value.OutboxID,
	v value.GlobalVersion,
) error {
	if err := t.XLock(ctx, tn, oid.String()); err != nil {
		return err
	}

	o, err := t.collection.Outbox.Get(ctx, tn, pid, oid)
	if err != nil {
		return err
	}

	prev := make([]byte, entity.OutboxSize)
	_ = o.WriteData(prev)

	o.UpdateGlobalVersion(v)

	data := make([]byte, entity.OutboxSize)
	_ = o.WriteData(data)

	_, err = t.log.Update(ctx, t.id, link.RecordID{}, data, prev)
	if err != nil {
		return err
	}

	return t.collection.Outbox.Replace(ctx, l, tn, pid, oid, o)
}

func (t *Transaction) ListOutboxes(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
) iter.Seq2[*entity.Outbox, error] {
	// TODO: How do we lock this query?
	// if err := t.concurrencyMgr.SLock(ctx, ItemID(oid.String())); err != nil {
	// 	return nil, err
	// }

	return t.collection.Outbox.List(ctx, tn, pid)
}

func (t *Transaction) InsertRequestID(
	ctx context.Context,
	l span.Log,
	tn tablename.TableName,
	rqid value.RequestID,
	now value.Time,
) error {
	return t.collection.Idempotency.Set(ctx, l, tn, value.NewPartitionID(0), rqid, entity.NewRequestLog(rqid, now))
}

func (t *Transaction) TestRequestID(
	ctx context.Context,
	tn tablename.TableName,
	rqid value.RequestID,
) (bool, error) {
	return t.collection.Idempotency.Test(ctx, tn, value.NewPartitionID(0), rqid)
}

func (l *Transaction) Append(
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
