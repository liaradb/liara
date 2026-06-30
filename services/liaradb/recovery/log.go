package recovery

import (
	"context"
	"iter"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
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
	lw logWriter
}

func NewLog(
	pageSize int64,
	segmentSize action.PageID,
	maxRecordSize int64,
	writeQueueSize int,
	fsys filecache.FileSystem,
	dir string,
) *Log {
	return &Log{
		lw: *newLogWriter(pageSize, segmentSize, maxRecordSize, writeQueueSize, fsys, dir),
	}
}

func (l *Log) HighWater() record.LogSequenceNumber { return l.lw.HighWater() }
func (l *Log) LowWater() record.LogSequenceNumber  { return l.lw.LowWater() }
func (l *Log) IsDirty() bool                       { return l.lw.IsDirty() }

func (l *Log) Run(ctx context.Context) error {
	return l.lw.Run(ctx)
}

func (l *Log) Close() error {
	return l.lw.Close()
}

func (l *Log) StartWriter() error {
	return l.lw.StartWriter()
}

func (l *Log) Start(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
) (record.LogSequenceNumber, error) {
	return l.lw.appendRecord(ctx, tid, txid, now, record.ActionStart, record.CollectionSystem, nil, nil)
}

func (l *Log) Commit(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
) (record.LogSequenceNumber, error) {
	return l.lw.appendAndWait(ctx, tid, txid, now, record.ActionCommit)
}

func (l *Log) Rollback(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
) (record.LogSequenceNumber, error) {
	return l.lw.appendAndWait(ctx, tid, txid, now, record.ActionRollback)
}

func (l *Log) Insert(
	ctx context.Context,
	tid value.TenantID,
	txid record.TransactionID,
	now time.Time,
	collection record.Collection,
	data []byte,
) (record.LogSequenceNumber, error) {
	return l.lw.appendRecord(ctx, tid, txid, now, record.ActionInsert, collection, data, nil)
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
	return l.lw.appendRecord(ctx, tid, txid, now, record.ActionUpdate, collection, data, prev)
}

func (l *Log) FlushCheckpoint(
	ctx context.Context,
	now time.Time,
	txids ...record.TransactionID,
) (record.LogSequenceNumber, error) {
	lsn, err := l.lw.appendCheckpoint(
		ctx,
		now,
		txids...)
	if err != nil {
		return record.LogSequenceNumber{}, err
	}

	if err := l.lw.flushPageQueue(ctx); err != nil {
		return record.LogSequenceNumber{}, err
	}

	return lsn, nil
}

func (l *Log) Iterate(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return l.lw.Iterate(lsn)
}

// Iterate in reverse until Checkpoint. Then iterate forward entil end of log.
func (l *Log) Recover() (iter.Seq[*record.Record], error) {
	return l.lw.Recover()
}

func (l *Log) Reverse() iter.Seq2[*record.Record, error] {
	return l.lw.Reverse()
}
