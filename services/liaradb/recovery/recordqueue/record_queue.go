package recordqueue

import (
	"context"
	"time"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/pagequeue"
	"github.com/liaradb/liaradb/recovery/pagequeue/pagestorage"
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
type RecordQueue[R pagequeue.Record] struct {
	pageSize      int64
	sl            *segment.List
	pq            *pagequeue.PageQueue[R]
	ps            *pagestorage.PageStorage
	fs            flushStatus
	appendReqs    AppendHandler[R]
	flushReqs     async.CommandHandler[struct{}]
	cancel        context.CancelFunc
	maxRecordSize int64
}

func New[R pagequeue.Record](
	pageSize int64,
	maxRecordSize int64,
	writeQueueSize int,
	sl *segment.List,
	pl *pagepool.PagePool,
) *RecordQueue[R] {
	ps := pagestorage.New(sl)
	return &RecordQueue[R]{
		pageSize:      pageSize,
		sl:            sl,
		pq:            pagequeue.New[R](ps, pl, writeQueueSize),
		ps:            ps,
		appendReqs:    NewAppendHandler[R](),
		flushReqs:     make(async.CommandHandler[struct{}], 1),
		maxRecordSize: maxRecordSize,
	}
}

func (rq *RecordQueue[R]) HighWater() logpage.LogSequenceNumber { return rq.fs.HighWater() }
func (rq *RecordQueue[R]) LowWater() logpage.LogSequenceNumber  { return rq.fs.LowWater() }

func (rq *RecordQueue[R]) Run(ctx context.Context) error {
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

func (rq *RecordQueue[R]) Close() error {
	if rq.cancel != nil {
		rq.cancel()
	}

	return rq.sl.Close()
}

func (rq *RecordQueue[R]) Init(lw, hw logpage.LogSequenceNumber) error {
	rq.fs.init(lw, hw)

	// TODO: Don't create a page, just copy the data
	data := make([]byte, rq.pageSize)
	if err := rq.ps.Init(data); err != nil {
		return err
	}

	return rq.pq.Init(data)
}

func (rq *RecordQueue[R]) run(ctx context.Context) {
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
func (rq *RecordQueue[R]) AppendAndWait(
	ctx context.Context,
	record R,
) (logpage.LogSequenceNumber, error) {
	return rq.appendReqs.AppendAndWait(ctx, record)
}

// Request thread
func (rq *RecordQueue[R]) Append(
	ctx context.Context,
	record R,
) (logpage.LogSequenceNumber, error) {
	// Verify that record can fit at all
	if record.Size() > int(rq.maxRecordSize) {
		return logpage.LogSequenceNumber{}, raw.ErrInsufficientSpace
	}

	return rq.appendReqs.Append(ctx, record)
}

func (rq *RecordQueue[R]) appendRequest(ctx context.Context, r *AppendRequest[R]) {
	hw := rq.append(ctx, rq.fs.HighWater(), r)
	rq.fs.setHighWater(hw)
}

func (rq *RecordQueue[R]) append(
	ctx context.Context,
	lsn logpage.LogSequenceNumber,
	r *AppendRequest[R],
) logpage.LogSequenceNumber {
	if r.Value().IsWait() {
		return rq.appendWait(ctx, lsn, r)
	} else {
		return rq.appendNoWait(ctx, lsn, r)
	}
}

func (rq *RecordQueue[R]) appendWait(
	ctx context.Context,
	lsn logpage.LogSequenceNumber,
	r *AppendRequest[R],
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

func (rq *RecordQueue[R]) appendNoWait(
	ctx context.Context,
	lsn logpage.LogSequenceNumber,
	r *AppendRequest[R],
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
func (rq *RecordQueue[R]) Flush(ctx context.Context) error {
	return rq.flushReqs.Send(ctx, struct{}{})
}

func (rq *RecordQueue[R]) flushRequest(
	ctx context.Context,
	r *async.Command[struct{}],
) {
	r.Reply(rq.flush(ctx))
}

func (rq *RecordQueue[R]) flushOrPanic(ctx context.Context) {
	if err := rq.flush(ctx); err != nil {
		panic(err)
	}
}

func (rq *RecordQueue[R]) flush(ctx context.Context) error {
	if !rq.fs.isDirty() {
		return nil
	}

	if err := rq.flushPageQueue(ctx); err != nil {
		return err
	}

	rq.fs.completeFlush()
	return nil
}

func (rq *RecordQueue[R]) flushPageQueue(ctx context.Context) error {
	return rq.pq.Flush(ctx)
}
