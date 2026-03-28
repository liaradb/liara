package transaction

import (
	"context"
	"errors"
	"io"
	"iter"
	"time"

	"github.com/liaradb/liaradb/collection"
	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/util/set"
)

type Transaction struct {
	id             record.TransactionID
	tid            value.TenantID
	log            *recovery.Log
	bufferList     *BufferList
	concurrencyMgr *locktable.ConcurrencyMgr[action.ItemID]
	collection     *collection.Collections
	events         []eventItem
	values         []valueItem
	keys           set.Set[key.Key]
	forceRollback  bool
	manager        *Manager
}

type eventItem struct {
	e    *entity.Event
	data []byte
}

type valueItem struct {
	r    *entity.Row
	data []byte
}

func newTransaction(
	id record.TransactionID,
	tid value.TenantID,
	log *recovery.Log,
	bufferList *BufferList,
	concurrencyMgr *locktable.ConcurrencyMgr[action.ItemID],
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
		keys:           set.Set[key.Key]{},
		manager:        manager,
	}
}

func (t *Transaction) ID() record.TransactionID { return t.id }

func (t *Transaction) GetAggregate(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	id value.AggregateID,
) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		if err := t.concurrencyMgr.SLock(ctx, action.ItemID(id.String())); err != nil {
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
		if err := t.concurrencyMgr.SLock(ctx, action.ItemID(id.String())); err != nil {
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
		// if err := t.concurrencyMgr.SLock(ctx, action.ItemID(id.String())); err != nil {
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
		// if err := t.concurrencyMgr.SLock(ctx, action.ItemID(id.String())); err != nil {
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

func (t *Transaction) Insert(
	ctx context.Context,
	tn tablename.TableName,
	now time.Time,
	e *entity.Event,
	data []byte,
) error {
	if err := t.concurrencyMgr.XLock(ctx, action.ItemID(e.AggregateID.String())); err != nil {
		return err
	}

	k := key.NewKey2(e.AggregateID.Bytes(), int64(e.Version.Value()))
	// Verify this AggregateID and Version is unique in this transaction
	if t.keys.Includes(k) {
		return btree.ErrExists
	}

	t.keys.Add(k)

	// Verify this AggregateID and Version is unique in Index
	if err := t.collection.EventLog.CanAppend(ctx, tn, e.PartitionID, k); err != nil {
		return err
	}

	_, err := t.log.Insert(ctx, t.tid, t.id, now, data)
	if err != nil {
		return err
	}

	t.events = append(t.events, eventItem{
		e:    e,
		data: data,
	})

	return nil
}

func (t *Transaction) SetValue(
	ctx context.Context,
	tn tablename.TableName,
	now time.Time,
	r *entity.Row,
	data []byte,
) error {
	if err := t.concurrencyMgr.XLock(ctx, action.ItemID(r.ID().String())); err != nil {
		return err
	}

	// k := key.NewKey2(r.ID().Bytes(), int64(r.Version().Value()))
	// // Verify this AggregateID and Version is unique in this transaction
	// if t.keys.Includes(k) {
	// 	return btree.ErrExists
	// }

	// t.keys.Add(k)

	// // Verify this AggregateID and Version is unique in Index
	// if err := t.keyValue.CanAppend(ctx, tn, r.PartitionID(), k); err != nil {
	// 	return err
	// }

	_, err := t.log.Insert(ctx, t.tid, t.id, now, data)
	if err != nil {
		return err
	}

	t.values = append(t.values, valueItem{
		r:    r,
		data: data,
	})

	return nil
}

func Run(
	ctx context.Context,
	t *Transaction,
	now time.Time,
	f func() error,
) error {
	_, err := t.run(ctx, now, func() (any, error) {
		return struct{}{}, f()
	})
	if err != nil {
		return errTransactionFailed(t.id, err)
	}

	return nil
}

func RunResult[R any](
	ctx context.Context,
	t *Transaction,
	now time.Time,
	f func() (R, error),
) (R, error) {
	r, err := t.run(ctx, now, func() (any, error) {
		return f()
	})
	if err != nil {
		var r R
		return r, errTransactionFailed(t.id, err)
	}

	return r.(R), nil
}

func (t *Transaction) run(
	ctx context.Context,
	now time.Time,
	f func() (any, error),
) (any, error) {
	_, err := t.log.Start(ctx, t.tid, t.id, now)
	if err != nil {
		return nil, err
	}

	defer t.release()

	r, err := f()
	if err != nil {
		return nil, errors.Join(err, t.rollback(ctx, now))
	}

	if t.forceRollback {
		return nil, t.rollback(ctx, now)
	}

	return r, t.commit(ctx, now)
}

func (t *Transaction) release() {
	t.concurrencyMgr.Release()
	t.bufferList.Release()
	t.manager.End(t.id)
}

func (t *Transaction) commit(
	ctx context.Context,
	now time.Time,
) error {
	_, err := t.log.Commit(ctx, t.tid, t.id, now)
	if err != nil {
		return err
	}

	if err := t.flush(ctx); err != nil {
		return err
	}

	if err := t.appendToEventLog(ctx); err != nil {
		return err
	}

	return nil
}

func (t *Transaction) appendToEventLog(ctx context.Context) error {
	tn := tablename.New(t.tid)
	for _, item := range t.events {
		k := key.NewKey2(item.e.AggregateID.Bytes(), item.e.Version.Value())
		err := t.collection.EventLog.AppendEvent(ctx, tn, item.e.PartitionID, k, item.data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Transaction) rollback(ctx context.Context, now time.Time) error {
	_, err := t.log.Rollback(ctx, t.tid, t.id, now)
	if err != nil {
		return err
	}

	if err := t.flush(ctx); err != nil {
		return err
	}

	return nil
}

func (t *Transaction) flush(ctx context.Context) error {
	return t.log.Flush(ctx)
}

func (t *Transaction) GetOutbox(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	oid value.OutboxID,
) (*entity.Outbox, error) {
	if err := t.concurrencyMgr.SLock(ctx, action.ItemID(oid.String())); err != nil {
		return nil, err
	}

	return t.collection.Outbox.Get(ctx, tn, pid, oid)
}

func (t *Transaction) InsertOutbox(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	now time.Time,
	oid value.OutboxID,
	e *entity.Outbox,
) error {
	if err := t.concurrencyMgr.XLock(ctx, action.ItemID(oid.String())); err != nil {
		return err
	}

	data := make([]byte, entity.OutboxSize)
	_, ok := e.Write(data)
	if !ok {
		return io.ErrUnexpectedEOF
	}

	_, err := t.log.Insert(ctx, t.tid, t.id, now, data)
	if err != nil {
		return err
	}

	return t.collection.Outbox.Set(ctx, tn, pid, oid, e)
}

func (t *Transaction) UpdateOutbox(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	now time.Time,
	oid value.OutboxID,
	v value.GlobalVersion,
) error {
	if err := t.concurrencyMgr.XLock(ctx, action.ItemID(oid.String())); err != nil {
		return err
	}

	o, err := t.collection.Outbox.Get(ctx, tn, pid, oid)
	if err != nil {
		return err
	}

	prev := make([]byte, entity.OutboxSize)
	_, ok := o.Write(prev)
	if !ok {
		return io.ErrUnexpectedEOF
	}

	o.UpdateGlobalVersion(v)

	data := make([]byte, entity.OutboxSize)
	_, ok = o.Write(data)
	if !ok {
		return io.ErrUnexpectedEOF
	}

	_, err = t.log.Update(ctx, t.tid, t.id, now, data, prev)
	if err != nil {
		return err
	}

	return t.collection.Outbox.Replace(ctx, tn, pid, oid, o)
}

func (t *Transaction) ListOutboxes(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
) iter.Seq2[*entity.Outbox, error] {
	// TODO: How do we lock this query?
	// if err := t.concurrencyMgr.SLock(ctx, action.ItemID(oid.String())); err != nil {
	// 	return nil, err
	// }

	return t.collection.Outbox.List(ctx, tn, pid)
}

func (t *Transaction) InsertRequestID(
	ctx context.Context,
	tn tablename.TableName,
	rqid value.RequestID,
	now value.Time,
) error {
	return t.collection.Idempotency.Set(ctx, tn, value.NewPartitionID(0), rqid, entity.NewRequestLog(rqid, now))
}

func (t *Transaction) TestRequestID(
	ctx context.Context,
	tn tablename.TableName,
	rqid value.RequestID,
) (bool, error) {
	return t.collection.Idempotency.Test(ctx, tn, value.NewPartitionID(0), rqid)
}
