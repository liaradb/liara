package recovery

import (
	"context"
	"iter"
	"time"

	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
)

const (
	interval = 100 * time.Millisecond
)

type Log struct {
	sl         *segment.List
	reader     *reader
	writer     *writer
	highWater  record.LogSequenceNumber
	lowWater   record.LogSequenceNumber
	appendReqs async.Handler[appendValue, record.LogSequenceNumber]
	flushReqs  async.CommandHandler[struct{}]
	cancel     context.CancelFunc
}

type flushRequest = async.Command[struct{}]

type appendRequest = async.Request[appendValue, record.LogSequenceNumber]

type appendValue struct {
	tid        value.TenantID
	txid       record.TransactionID
	time       time.Time
	action     record.Action
	collection record.Collection
	data       []byte
	reverse    []byte
}

func NewLog(
	pageSize int64,
	segmentSize action.PageID,
	fsys filecache.FileSystem,
	dir string,
) *Log {
	sl := segment.NewList(fsys, dir)
	return &Log{
		sl:         sl,
		reader:     newReader(pageSize, sl),
		writer:     newWriter(pageSize, segmentSize, sl),
		appendReqs: make(chan *appendRequest),
		flushReqs:  make(chan *flushRequest),
	}
}

func (l *Log) HighWater() record.LogSequenceNumber { return l.highWater }
func (l *Log) LowWater() record.LogSequenceNumber  { return l.lowWater }
func (l *Log) PageID() action.PageID               { return l.writer.PageID() }
func (l *Log) IsDirty() bool                       { return l.lowWater != l.highWater }

func (l *Log) run(ctx context.Context) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case r := <-l.appendReqs:
			l.appendRequest(r)
		case <-ticker.C:
			if err := l.flush(); err != nil {
				panic(err)
			}
		case r := <-l.flushReqs:
			l.flushRequest(r)
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

func (l *Log) appendRequest(r *appendRequest) {
	v := r.Value()
	lsn, err := l.append(v.tid, v.txid, v.time, v.action, v.collection, v.data, v.reverse)
	r.Reply(lsn, err)
}

func (l *Log) append(
	tid value.TenantID,
	txid record.TransactionID,
	time time.Time,
	action record.Action,
	collection record.Collection,
	data []byte,
	reverse []byte,
) (record.LogSequenceNumber, error) {
	h := l.highWater.Increment()
	rc := record.New(h, tid, txid, record.NewTime(time), action, collection, data, reverse)
	if err := l.writer.Append(rc); err != nil {
		return record.NewLogSequenceNumber(0), err
	}

	l.highWater = h
	return l.highWater, nil
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

	return lsn, l.requestFlush(ctx)
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

	return lsn, l.requestFlush(ctx)
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

func (l *Log) Close() error {
	if l.cancel != nil {
		l.cancel()
	}

	return l.sl.Close()
}

func (l *Log) requestFlush(ctx context.Context) error {
	return l.flushReqs.Send(ctx, struct{}{})
}

func (l *Log) flushRequest(r *flushRequest) {
	err := l.flush()
	r.Reply(err)
}

func (l *Log) flush() error {
	if !l.IsDirty() {
		return nil
	}

	if err := l.writer.Flush(); err != nil {
		return err
	}

	l.lowWater = l.highWater
	return nil
}

func (l *Log) Iterate(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return l.reader.iterate(lsn)
}

func (l *Log) Open(ctx context.Context) error {
	if err := l.sl.Open(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	l.cancel = cancel
	go l.run(ctx)
	return nil
}

func (l *Log) Recover() (iter.Seq[*record.Record], error) {
	return l.reader.recover()
}

func (l *Log) Reverse() iter.Seq2[*record.Record, error] {
	return l.reader.reverse()
}

func (l *Log) StartWriter() error {
	return l.writer.Start()
}

func (l *Log) FlushCheckpoint(
	now time.Time,
	txids ...record.TransactionID,
) (record.LogSequenceNumber, error) {
	data := l.txIDsToData(txids)
	lsn, err := l.append(value.TenantID{}, record.TransactionID{}, now, record.ActionCheckpoint, record.CollectionSystem, data, nil)
	if err != nil {
		return record.LogSequenceNumber{}, err
	}

	if err := l.writer.Flush(); err != nil {
		return record.LogSequenceNumber{}, err
	}

	return lsn, nil
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
