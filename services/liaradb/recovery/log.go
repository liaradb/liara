package recovery

import (
	"container/list"
	"context"
	"iter"
	"time"

	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pageiterator"
	"github.com/liaradb/liaradb/recovery/pagequeue"
	"github.com/liaradb/liaradb/recovery/pagestorage"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/util/iterator"
)

const (
	interval = 100 * time.Millisecond
)

// Append process
//  1. Append log record
//  2. Fill page
//  3. If page is full, flush
//  4. Request flush to LSN
//  5. Wait for timeout
//  6. Flush to LSN, notify requester
//
// What happens if we flush previous page or segment?
// Do we flush current page when closing segment?
type Log struct {
	pageSize      int64
	sl            *segment.List
	pq            *pagequeue.PageQueue
	ps            *pagestorage.PageStorage
	it            *pageiterator.PageIterator
	highWater     record.LogSequenceNumber
	lowWater      record.LogSequenceNumber
	appendReqs    async.Handler[appendValue, record.LogSequenceNumber]
	flushReqs     async.CommandHandler[record.LogSequenceNumber]
	cancel        context.CancelFunc
	queue         requestQueue
	maxRecordSize int64
}

type flushRequest = async.Command[record.LogSequenceNumber]

func NewLog(
	pageSize int64,
	segmentSize action.PageID,
	maxRecordSize int64,
	writeQueueSize int,
	fsys filecache.FileSystem,
	dir string,
) *Log {
	sl := segment.NewList(fsys, dir, pageSize, segmentSize)
	ps := pagestorage.New(sl)
	l := &Log{
		pageSize:      pageSize,
		sl:            sl,
		ps:            ps,
		it:            pageiterator.New(sl, int16(pageSize), logpage.HeaderSize),
		appendReqs:    make(chan *appendRequest),
		flushReqs:     make(chan *flushRequest),
		maxRecordSize: maxRecordSize,
	}

	l.pq = pagequeue.New(ps, l, int16(pageSize), logpage.HeaderSize, writeQueueSize)
	return l
}

func (l *Log) HighWater() record.LogSequenceNumber { return l.highWater }
func (l *Log) LowWater() record.LogSequenceNumber  { return l.lowWater }
func (l *Log) IsDirty() bool                       { return l.lowWater != l.highWater }

func (l *Log) Run(ctx context.Context) error {
	if err := l.sl.Open(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	l.cancel = cancel
	go l.run(ctx)
	go func() {
		_ = l.pq.Run(ctx)
	}()
	return nil
}

func (l *Log) Close() error {
	if l.cancel != nil {
		l.cancel()
	}

	return l.sl.Close()
}

func (l *Log) StartWriter() error {
	if err := l.initHighWater(); err != nil {
		return err
	}

	// TODO: Don't create a page, just copy the data
	data := make([]byte, l.pageSize)
	if err := l.ps.Init(data); err != nil {
		return err
	}

	l.pq.Init(data)
	return nil
}

func (l *Log) initHighWater() error {
	l.lowWater = record.NewLogSequenceNumber(0)
	l.highWater = record.NewLogSequenceNumber(0)

	hw := false
	for rc, err := range l.it.Reverse() {
		if err != nil {
			return err
		}

		if !hw {
			l.highWater = rc.LogSequenceNumber()
			hw = true
		}

		if rc.Action() == record.ActionCheckpoint {
			l.lowWater = rc.LogSequenceNumber()
			break
		}
	}

	return nil
}

func (l *Log) Start(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
) (record.LogSequenceNumber, error) {
	return l.appendRecord(ctx, tid, txid, now, record.ActionStart, record.CollectionSystem, nil, nil)
}

func (l *Log) Commit(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
) (record.LogSequenceNumber, error) {
	lsn, err := l.appendRecord(ctx, tid, txid, now, record.ActionCommit, record.CollectionSystem, nil, nil)
	if err != nil {
		return lsn, err
	}

	return lsn, l.requestFlush(ctx, lsn)
}

func (l *Log) Rollback(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
) (record.LogSequenceNumber, error) {
	lsn, err := l.appendRecord(ctx, tid, txid, now, record.ActionRollback, record.CollectionSystem, nil, nil)
	if err != nil {
		return lsn, err
	}

	return lsn, l.requestFlush(ctx, lsn)
}

func (l *Log) Insert(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
	collection record.Collection,
	data []byte,
) (record.LogSequenceNumber, error) {
	return l.appendRecord(ctx, tid, txid, now, record.ActionInsert, collection, data, nil)
}

func (l *Log) Update(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
	collection record.Collection,
	data []byte,
	prev []byte,
) (record.LogSequenceNumber, error) {
	return l.appendRecord(ctx, tid, txid, now, record.ActionUpdate, collection, data, prev)
}

func (l *Log) FlushCheckpoint(
	ctx context.Context,
	now time.Time,
	txids ...record.TransactionID,
) (record.LogSequenceNumber, error) {
	lsn, err := l.appendCheckpoint(
		ctx,
		now,
		txids...)
	if err != nil {
		return record.LogSequenceNumber{}, err
	}

	if err := l.flushPageQueue(ctx); err != nil {
		return record.LogSequenceNumber{}, err
	}

	return lsn, nil
}

func (l *Log) Iterate(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return l.it.Forward(lsn)
}

// Iterate in reverse until Checkpoint. Then iterate forward entil end of log.
func (l *Log) Recover() (iter.Seq[*record.Record], error) {
	rcs := list.New()

	for rc, err := range l.it.Reverse() {
		if err != nil {
			return nil, err
		}

		if rc.IsCheckpoint() {
			break
		}

		rcs.PushBack(rc)
	}

	return iterator.Reverse[*record.Record](rcs), nil
}

func (l *Log) Reverse() iter.Seq2[*record.Record, error] {
	return l.it.Reverse()
}

func (l *Log) run(ctx context.Context) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case r := <-l.appendReqs:
			l.appendRequest(ctx, r)
		case r := <-l.flushReqs:
			l.flushRequest(r)
		case <-ticker.C:
			l.flushOrPanic(ctx)
		}
	}
}

func (l *Log) appendRecord(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	time time.Time,
	action record.Action,
	collection record.Collection,
	data []byte,
	reverse []byte,
) (record.LogSequenceNumber, error) {
	// Verify that record can fit at all
	if len(data) > int(l.maxRecordSize) {
		return record.LogSequenceNumber{}, raw.ErrInsufficientSpace
	}

	return l.appendReqs.Send(ctx, appendValue{
		tid:        tid,
		txid:       txid,
		time:       time,
		action:     action,
		collection: collection,
		data:       data,
		reverse:    reverse,
	})
}

func (l *Log) appendRequest(ctx context.Context, r *appendRequest) {
	v := r.Value()
	h := l.highWater.Increment()

	err := l.pq.Append(ctx, v.Record(h))
	if err == nil {
		l.highWater = h
	}

	r.Reply(l.highWater, err)
}

// # Append to PageQueue
//   - If current page was full already, do not sync
//   - Otherwise, sync current page
//   - For every page after, push new page to disk
//   - For last page, store that as current
func (l *Log) appendCheckpoint(
	ctx context.Context,
	time time.Time,
	txids ...record.TransactionID,
) (record.LogSequenceNumber, error) {
	h := l.highWater.Increment()
	data := l.txIDsToData(txids)
	rc := record.New(h,
		value.TenantID{},
		record.TransactionID{},
		record.NewTime(time),
		record.ActionCheckpoint,
		record.CollectionSystem,
		data,
		nil)

	if err := l.pq.Append(ctx, rc); err != nil {
		return record.NewLogSequenceNumber(0), err
	}

	l.highWater = h
	return l.highWater, nil
}

func (*Log) txIDsToData(txids []record.TransactionID) []byte {
	data := make([]byte, len(txids)*record.TransactionIDSize)

	data0 := data
	for _, txid := range txids {
		// There will always be enough space
		data0, _ = txid.WriteData(data0)
	}

	return data
}

func (l *Log) requestFlush(ctx context.Context, lsn record.LogSequenceNumber) error {
	return l.flushReqs.Send(ctx, lsn)
}

func (l *Log) flushRequest(r *flushRequest) {
	l.queue.add(r)
}

func (l *Log) flushOrPanic(ctx context.Context) {
	if err := l.flush(ctx); err != nil {
		panic(err)
	}
}

func (l *Log) flush(ctx context.Context) error {
	if !l.IsDirty() {
		return nil
	}

	if err := l.flushPageQueue(ctx); err != nil {
		return err
	}

	l.completeFlush()
	return nil
}

func (l *Log) flushPageQueue(ctx context.Context) error {
	return l.pq.Flush(ctx)
}

func (l *Log) completeFlush() {
	l.lowWater = l.highWater
	l.queue.sendUpToLSN(l.highWater)
}

func (pq *Log) OnFlush(lsn record.LogSequenceNumber) {
}
