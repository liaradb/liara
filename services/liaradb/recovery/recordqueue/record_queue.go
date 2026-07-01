package recordqueue

import (
	"context"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/pageiterator"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/pagequeue"
	"github.com/liaradb/liaradb/recovery/pagequeue/pagestorage"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
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
type RecordQueue struct {
	pageSize      int64
	sl            *segment.List
	pq            *pagequeue.PageQueue
	fs            flushStatus
	appendReqs    pagequeue.AppendHandler
	cancel        context.CancelFunc
	maxRecordSize int64
}

func New(
	pageSize int64,
	maxRecordSize int64,
	writeQueueSize int,
	sl *segment.List,
	pl *pagepool.PagePool,
) *RecordQueue {
	ps := pagestorage.New(sl)
	l := &RecordQueue{
		pageSize:      pageSize,
		sl:            sl,
		appendReqs:    pagequeue.NewAppendHandler(),
		maxRecordSize: maxRecordSize,
	}

	l.pq = pagequeue.New(ps, pl, writeQueueSize)
	return l
}

func (lw *RecordQueue) HighWater() record.LogSequenceNumber { return lw.fs.HighWater() }
func (lw *RecordQueue) LowWater() record.LogSequenceNumber  { return lw.fs.LowWater() }

func (lw *RecordQueue) Run(ctx context.Context) error {
	if err := lw.sl.Open(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	lw.cancel = cancel
	go lw.run(ctx)
	go func() {
		_ = lw.pq.Run(ctx)
	}()
	return nil
}

func (lw *RecordQueue) Close() error {
	if lw.cancel != nil {
		lw.cancel()
	}

	return lw.sl.Close()
}

func (lw *RecordQueue) StartWriter(it *pageiterator.PageIterator) error {
	if err := lw.fs.initHighWater(it); err != nil {
		return err
	}

	return lw.pq.Init(lw.pageSize)
}

func (lw *RecordQueue) run(ctx context.Context) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case r := <-lw.appendReqs.Reqs():
			lw.appendRequest(ctx, r)
		case <-ticker.C:
			lw.flushOrPanic(ctx)
		}
	}
}

func (lw *RecordQueue) AppendAndWait(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	time time.Time,
	action record.Action,
) (record.LogSequenceNumber, error) {
	return lw.appendReqs.AppendAndWait(ctx,
		tid,
		txid,
		time,
		action,
		record.CollectionSystem,
		nil,
		nil,
	)
}

func (lw *RecordQueue) AppendRecord(
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
	if len(data) > int(lw.maxRecordSize) {
		return record.LogSequenceNumber{}, raw.ErrInsufficientSpace
	}

	return lw.appendReqs.Append(ctx,
		tid,
		txid,
		time,
		action,
		collection,
		data,
		reverse,
	)
}

func (lw *RecordQueue) appendRequest(ctx context.Context, r *pagequeue.AppendRequest) {
	lw.fs.setHighWater(lw.pq.AppendRequest(ctx, lw.fs.HighWater(), r))
}

func (l *RecordQueue) FlushCheckpoint(
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

// # Append to PageQueue
//   - If current page was full already, do not sync
//   - Otherwise, sync current page
//   - For every page after, push new page to disk
//   - For last page, store that as current
func (lw *RecordQueue) appendCheckpoint(
	ctx context.Context,
	time time.Time,
	txids ...record.TransactionID,
) (record.LogSequenceNumber, error) {
	h := lw.fs.HighWater().Increment()
	data := lw.txIDsToData(txids)
	rc := record.New(h,
		value.TenantID{},
		record.TransactionID{},
		record.NewTime(time),
		record.ActionCheckpoint,
		record.CollectionSystem,
		data,
		nil)

	if err := lw.pq.Append(ctx, rc); err != nil {
		return record.NewLogSequenceNumber(0), err
	}

	lw.fs.setHighWater(h)
	return h, nil
}

func (*RecordQueue) txIDsToData(txids []record.TransactionID) []byte {
	data := make([]byte, len(txids)*record.TransactionIDSize)

	data0 := data
	for _, txid := range txids {
		// There will always be enough space
		data0, _ = txid.WriteData(data0)
	}

	return data
}

func (lw *RecordQueue) flushOrPanic(ctx context.Context) {
	if err := lw.flush(ctx); err != nil {
		panic(err)
	}
}

func (lw *RecordQueue) flush(ctx context.Context) error {
	if !lw.fs.isDirty() {
		return nil
	}

	if err := lw.flushPageQueue(ctx); err != nil {
		return err
	}

	lw.fs.completeFlush()
	return nil
}

func (lw *RecordQueue) flushPageQueue(ctx context.Context) error {
	return lw.pq.Flush(ctx)
}
