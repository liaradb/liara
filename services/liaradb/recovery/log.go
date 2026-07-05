package recovery

import (
	"container/list"
	"context"
	"iter"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/pageiterator"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/recordqueue"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/recovery/span"
	"github.com/liaradb/liaradb/util/iterator"
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
	rq recordqueue.RecordQueue
	it *pageiterator.PageIterator
}

func NewLog(
	pageSize int64,
	segmentSize action.PageID,
	maxRecordSize int64,
	writeQueueSize int,
	fsys filecache.FileSystem,
	dir string,
) *Log {
	sl := segment.NewList(fsys, dir, pageSize, segmentSize)
	// TODO: Fix this cast
	pl := pagepool.New(int(pageSize), span.FragmentHeaderSize)
	it := pageiterator.New(sl, pl)
	return &Log{
		rq: *recordqueue.New(
			pageSize,
			maxRecordSize,
			writeQueueSize,
			sl,
			pl),
		it: it,
	}
}

func (l *Log) HighWater() record.LogSequenceNumber { return l.rq.HighWater() }
func (l *Log) LowWater() record.LogSequenceNumber  { return l.rq.LowWater() }

func (l *Log) Run(ctx context.Context) error {
	return l.rq.Run(ctx)
}

func (l *Log) Close() error {
	return l.rq.Close()
}

func (l *Log) StartWriter() error {
	lw, hw, err := l.getFlushStatus()
	if err != nil {
		return err
	}

	return l.rq.Init(lw, hw)
}

func (l *Log) getFlushStatus() (lowWater, highWater record.LogSequenceNumber, err error) {
	hw := false
	var rc *record.Record
	for rc, err = range l.it.Reverse() {
		if err != nil {
			return
		}

		if !hw {
			highWater = rc.LogSequenceNumber()
			hw = true
		}

		if rc.Action() == record.ActionCheckpoint {
			lowWater = rc.LogSequenceNumber()
			break
		}
	}

	return
}

func (l *Log) Start(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
) (record.LogSequenceNumber, error) {
	return l.rq.Append(ctx,
		tid,
		txid,
		now,
		record.ActionStart,
		record.CollectionSystem,
		nil,
		nil)
}

func (l *Log) Commit(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
) (record.LogSequenceNumber, error) {
	return l.rq.AppendAndWait(ctx,
		tid,
		txid,
		now,
		record.ActionCommit)
}

func (l *Log) Rollback(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
) (record.LogSequenceNumber, error) {
	return l.rq.AppendAndWait(ctx,
		tid,
		txid,
		now,
		record.ActionRollback)
}

func (l *Log) Insert(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
	collection record.Collection,
	data []byte,
) (record.LogSequenceNumber, error) {
	return l.rq.Append(ctx,
		tid,
		txid,
		now,
		record.ActionInsert,
		collection,
		data,
		nil)
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
	return l.rq.Append(ctx,
		tid,
		txid,
		now,
		record.ActionUpdate,
		collection,
		data,
		prev)
}

func (l *Log) FlushCheckpoint(
	ctx context.Context,
	now time.Time,
	txids ...record.TransactionID,
) (record.LogSequenceNumber, error) {
	return l.rq.FlushCheckpoint(ctx,
		now,
		txids...)
}

func (l *Log) Flush(ctx context.Context) error {
	return l.rq.Flush(ctx)
}

// Iterate in reverse until Checkpoint. Then iterate forward entil end of log.
func (l *Log) Recover() (iter.Seq[*record.Record], error) {
	// TODO: Should we use a list, or just iterate forwards?
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
