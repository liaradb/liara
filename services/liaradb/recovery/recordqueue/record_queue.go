package recordqueue

import (
	"context"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/pagequeue"
	"github.com/liaradb/liaradb/recovery/pagequeue/pagestorage"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/util/async"
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
	ps            *pagestorage.PageStorage
	fs            flushStatus
	appendReqs    AppendHandler
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
		appendReqs:    NewAppendHandler(),
		flushReqs:     make(async.CommandHandler[struct{}], 1),
		maxRecordSize: maxRecordSize,
	}
}

func (rq *RecordQueue) HighWater() logpage.LogSequenceNumber { return rq.fs.HighWater() }
func (rq *RecordQueue) LowWater() logpage.LogSequenceNumber  { return rq.fs.LowWater() }

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

func (rq *RecordQueue) Init(lw, hw logpage.LogSequenceNumber) error {
	rq.fs.init(lw, hw)

	// TODO: Don't create a page, just copy the data
	data := make([]byte, rq.pageSize)
	if err := rq.ps.Init(data); err != nil {
		return err
	}

	return rq.pq.Init(data)
}

func (rq *RecordQueue) run(ctx context.Context) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case r := <-rq.appendReqs.Reqs():
			rq.appendRequest(ctx, r)
		case r := <-rq.flushReqs:
			rq.flushRequest(ctx, r)
		case <-ticker.C:
			rq.flushOrPanic(ctx)
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
) (logpage.LogSequenceNumber, error) {
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
) (logpage.LogSequenceNumber, error) {
	// Verify that record can fit at all
	if len(data) > int(rq.maxRecordSize) {
		return logpage.LogSequenceNumber{}, raw.ErrInsufficientSpace
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

func (rq *RecordQueue) appendRequest(ctx context.Context, r *AppendRequest) {
	hw := rq.append(ctx, rq.fs.HighWater(), r)
	rq.fs.setHighWater(hw)
}

func (rq *RecordQueue) append(
	ctx context.Context,
	lsn logpage.LogSequenceNumber,
	r *AppendRequest,
) logpage.LogSequenceNumber {
	if r.Value().IsWait() {
		return rq.appendWait(ctx, lsn, r)
	} else {
		return rq.appendNoWait(ctx, lsn, r)
	}
}

func (rq *RecordQueue) appendWait(
	ctx context.Context,
	lsn logpage.LogSequenceNumber,
	r *AppendRequest,
) logpage.LogSequenceNumber {
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
	lsn logpage.LogSequenceNumber,
	r *AppendRequest,
) logpage.LogSequenceNumber {
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

// TODO: This is unused
func (rq *RecordQueue) Flush(ctx context.Context) error {
	return rq.flushReqs.Send(ctx, struct{}{})
}

func (rq *RecordQueue) flushRequest(
	ctx context.Context,
	r *async.Command[struct{}],
) {
	r.Reply(rq.flush(ctx))
}

func (rq *RecordQueue) flushOrPanic(ctx context.Context) {
	if err := rq.flush(ctx); err != nil {
		panic(err)
	}
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
