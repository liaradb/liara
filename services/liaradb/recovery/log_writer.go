package recovery

import (
	"container/list"
	"context"
	"iter"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/pageiterator"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/pagequeue"
	"github.com/liaradb/liaradb/recovery/pagequeue/pagestorage"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/recovery/span"
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
type logWriter struct {
	pageSize      int64
	sl            *segment.List
	pq            *pagequeue.PageQueue
	it            *pageiterator.PageIterator
	highWater     record.LogSequenceNumber
	lowWater      record.LogSequenceNumber
	appendReqs    pagequeue.AppendHandler
	cancel        context.CancelFunc
	maxRecordSize int64
}

func newLogWriter(
	pageSize int64,
	segmentSize action.PageID,
	maxRecordSize int64,
	writeQueueSize int,
	fsys filecache.FileSystem,
	dir string,
) *logWriter {
	sl := segment.NewList(fsys, dir, pageSize, segmentSize)
	ps := pagestorage.New(sl)
	pl := pagepool.New(int16(pageSize), span.FragmentHeaderSize)
	l := &logWriter{
		pageSize:      pageSize,
		sl:            sl,
		it:            pageiterator.New(sl, pl),
		appendReqs:    pagequeue.NewAppendHandler(),
		maxRecordSize: maxRecordSize,
	}

	l.pq = pagequeue.New(ps, pl, writeQueueSize)
	return l
}

func (lw *logWriter) HighWater() record.LogSequenceNumber { return lw.highWater }
func (lw *logWriter) LowWater() record.LogSequenceNumber  { return lw.lowWater }
func (lw *logWriter) IsDirty() bool                       { return lw.lowWater != lw.highWater }

func (lw *logWriter) Run(ctx context.Context) error {
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

func (lw *logWriter) Close() error {
	if lw.cancel != nil {
		lw.cancel()
	}

	return lw.sl.Close()
}

func (lw *logWriter) StartWriter() error {
	if err := lw.initHighWater(); err != nil {
		return err
	}

	// TODO: Don't create a page, just copy the data
	data := make([]byte, lw.pageSize)
	return lw.pq.Init(data)
}

func (lw *logWriter) initHighWater() error {
	lw.lowWater = record.NewLogSequenceNumber(0)
	lw.highWater = record.NewLogSequenceNumber(0)

	hw := false
	for rc, err := range lw.it.Reverse() {
		if err != nil {
			return err
		}

		if !hw {
			lw.highWater = rc.LogSequenceNumber()
			hw = true
		}

		if rc.Action() == record.ActionCheckpoint {
			lw.lowWater = rc.LogSequenceNumber()
			break
		}
	}

	return nil
}

func (lw *logWriter) Iterate(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return lw.it.Forward(lsn)
}

// Iterate in reverse until Checkpoint. Then iterate forward entil end of log.
func (lw *logWriter) Recover() (iter.Seq[*record.Record], error) {
	rcs := list.New()

	for rc, err := range lw.it.Reverse() {
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

func (lw *logWriter) Reverse() iter.Seq2[*record.Record, error] {
	return lw.it.Reverse()
}

func (lw *logWriter) run(ctx context.Context) {
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

func (lw *logWriter) appendAndWait(
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

func (lw *logWriter) appendRecord(
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

func (lw *logWriter) appendRequest(ctx context.Context, r *pagequeue.AppendRequest) {
	lw.highWater = lw.pq.AppendRequest(ctx, lw.highWater, r)
}

func (l *logWriter) flushCheckpoint(
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
func (lw *logWriter) appendCheckpoint(
	ctx context.Context,
	time time.Time,
	txids ...record.TransactionID,
) (record.LogSequenceNumber, error) {
	h := lw.highWater.Increment()
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

	lw.highWater = h
	return lw.highWater, nil
}

func (*logWriter) txIDsToData(txids []record.TransactionID) []byte {
	data := make([]byte, len(txids)*record.TransactionIDSize)

	data0 := data
	for _, txid := range txids {
		// There will always be enough space
		data0, _ = txid.WriteData(data0)
	}

	return data
}

func (lw *logWriter) flushOrPanic(ctx context.Context) {
	if err := lw.flush(ctx); err != nil {
		panic(err)
	}
}

func (lw *logWriter) flush(ctx context.Context) error {
	if !lw.IsDirty() {
		return nil
	}

	if err := lw.flushPageQueue(ctx); err != nil {
		return err
	}

	lw.completeFlush()
	return nil
}

func (lw *logWriter) flushPageQueue(ctx context.Context) error {
	return lw.pq.Flush(ctx)
}

func (lw *logWriter) completeFlush() {
	lw.lowWater = lw.highWater
}
