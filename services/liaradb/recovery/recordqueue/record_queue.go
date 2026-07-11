package recordqueue

import (
	"context"
	"time"

	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/pagequeue"
	"github.com/liaradb/liaradb/recovery/pagequeue/pagestorage"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
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
	ps            *pagestorage.PageStorage
	fs            flushStatus
	appendReqs    pagequeue.AppendHandler
	flushReqs     async.CommandHandler[struct{}]
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
	return &RecordQueue{
		pageSize:      pageSize,
		sl:            sl,
		pq:            pagequeue.New(ps, pl, writeQueueSize),
		ps:            ps,
		appendReqs:    pagequeue.NewAppendHandler(),
		flushReqs:     make(async.CommandHandler[struct{}], 1),
		maxRecordSize: maxRecordSize,
	}
}

func (rq *RecordQueue) HighWater() record.LogSequenceNumber { return rq.fs.HighWater() }
func (rq *RecordQueue) LowWater() record.LogSequenceNumber  { return rq.fs.LowWater() }

func (rq *RecordQueue) Run(ctx context.Context) error {
	if err := rq.sl.Open(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	rq.cancel = cancel
	go rq.run(ctx)
	go func() {
		_ = rq.pq.Run(ctx)
	}()
	return nil
}

func (rq *RecordQueue) Close() error {
	if rq.cancel != nil {
		rq.cancel()
	}

	return rq.sl.Close()
}

func (rq *RecordQueue) Init(lw, hw record.LogSequenceNumber) error {
	rq.fs.init(lw, hw)

	// TODO: Don't create a page, just copy the data
	data := make([]byte, rq.pageSize)
	if err := rq.ps.Init(data); err != nil {
		return err
	}

	return rq.pq.Init(data)
}

func (rq *RecordQueue) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case r := <-rq.appendReqs.Reqs():
			rq.appendRequest(ctx, r)
		case r := <-rq.flushReqs:
			rq.flushRequest(ctx, r)
		}
	}
}

// Request thread
func (rq *RecordQueue) AppendAndWait(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
	action record.Action,
) (record.LogSequenceNumber, error) {
	return rq.appendReqs.AppendAndWait(ctx,
		tid,
		txid,
		now,
		action,
		record.CollectionSystem,
		nil,
		nil,
	)
}

// Request thread
func (rq *RecordQueue) Append(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
	action record.Action,
	collection record.Collection,
	data []byte,
	reverse []byte,
) (record.LogSequenceNumber, error) {
	// Verify that record can fit at all
	if len(data) > int(rq.maxRecordSize) {
		return record.LogSequenceNumber{}, raw.ErrInsufficientSpace
	}

	return rq.appendReqs.Append(ctx,
		tid,
		txid,
		now,
		action,
		collection,
		data,
		reverse,
	)
}

func (rq *RecordQueue) appendRequest(ctx context.Context, r *pagequeue.AppendRequest) {
	hw := rq.AppendRequest(ctx, rq.fs.HighWater(), r)
	rq.fs.setHighWater(hw)
}

func (rq *RecordQueue) AppendRequest(
	ctx context.Context,
	lsn record.LogSequenceNumber,
	r *pagequeue.AppendRequest,
) record.LogSequenceNumber {
	if r.Value().IsWait() {
		return rq.appendWait(ctx, lsn, r)
	} else {
		return rq.appendNoWait(ctx, lsn, r)
	}
}

func (rq *RecordQueue) appendWait(
	ctx context.Context,
	lsn record.LogSequenceNumber,
	r *pagequeue.AppendRequest,
) record.LogSequenceNumber {
	h := lsn.Increment()
	v := r.Value()
	err := rq.pq.AppendWait(ctx, v.Record(h), func() {
		r.Reply(h, nil)
	})

	if err == nil {
		return h
	} else {
		r.Reply(lsn, err)
		return lsn
	}
}

func (rq *RecordQueue) appendNoWait(
	ctx context.Context,
	lsn record.LogSequenceNumber,
	r *pagequeue.AppendRequest,
) record.LogSequenceNumber {
	h := lsn.Increment()
	v := r.Value()
	err := rq.pq.Append(ctx, v.Record(h))
	if err == nil {
		r.Reply(h, err)
		return h
	} else {
		r.Reply(lsn, err)
		return lsn
	}
}

// Manager thread
func (rq *RecordQueue) FlushCheckpoint(
	ctx context.Context,
	now time.Time,
	txids ...record.TransactionID,
) (record.LogSequenceNumber, error) {
	lsn, err := rq.appendCheckpoint(
		ctx,
		now,
		txids...)
	if err != nil {
		return record.LogSequenceNumber{}, err
	}

	if err := rq.flushPageQueue(ctx); err != nil {
		return record.LogSequenceNumber{}, err
	}

	return lsn, nil
}

// # Append to PageQueue
//   - If current page was full already, do not sync
//   - Otherwise, sync current page
//   - For every page after, push new page to disk
//   - For last page, store that as current
func (rq *RecordQueue) appendCheckpoint(
	ctx context.Context,
	now time.Time,
	txids ...record.TransactionID,
) (record.LogSequenceNumber, error) {
	h := rq.fs.HighWater().Increment()
	data := rq.txIDsToData(txids)
	rc := record.New(h,
		value.TenantID{},
		record.TransactionID{},
		record.NewTime(now),
		record.ActionCheckpoint,
		record.CollectionSystem,
		data,
		nil)

	if err := rq.pq.Append(ctx, rc); err != nil {
		return record.NewLogSequenceNumber(0), err
	}

	rq.fs.setHighWater(h)
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

// Manager thread
func (rq *RecordQueue) Flush(ctx context.Context) error {
	return rq.flushReqs.Send(ctx, struct{}{})
}

func (rq *RecordQueue) flushRequest(
	ctx context.Context,
	r *async.Command[struct{}],
) {
	r.Reply(rq.flush(ctx))
}

func (rq *RecordQueue) flush(ctx context.Context) error {
	if !rq.fs.isDirty() {
		return nil
	}

	if err := rq.flushPageQueue(ctx); err != nil {
		return err
	}

	rq.fs.completeFlush()
	return nil
}

func (rq *RecordQueue) flushPageQueue(ctx context.Context) error {
	return rq.pq.Flush(ctx)
}
