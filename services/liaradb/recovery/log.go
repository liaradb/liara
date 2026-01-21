package recovery

import (
	"context"
	"iter"
	"time"

	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/node"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
)

type Log struct {
	sl         *segment.List
	reader     *reader
	writer     *writer
	highWater  record.LogSequenceNumber
	lowWater   record.LogSequenceNumber
	appendReqs async.Handler[appendValue, record.LogSequenceNumber]
	flushReqs  async.CommandHandler[record.LogSequenceNumber]
	cancel     context.CancelFunc
}

type flushRequest = async.Command[record.LogSequenceNumber]

type appendRequest = async.Request[appendValue, record.LogSequenceNumber]

type appendValue struct {
	tid     record.TransactionID
	time    time.Time
	action  record.Action
	data    []byte
	reverse []byte
}

func NewLog(
	pageSize int64,
	segmentSize action.PageID,
	fsys file.FileSystem,
	dir string,
) *Log {
	sl := segment.NewList(fsys, dir)
	p := node.New(pageSize)
	return &Log{
		sl:         sl,
		reader:     newReader(pageSize, sl, p),
		writer:     newWriter(pageSize, segmentSize, sl, p),
		appendReqs: make(chan *appendRequest),
		flushReqs:  make(chan *flushRequest),
	}
}

func (l *Log) HighWater() record.LogSequenceNumber { return l.highWater }
func (l *Log) LowWater() record.LogSequenceNumber  { return l.lowWater }
func (l *Log) PageID() action.PageID               { return l.writer.PageID() }

func (l *Log) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case r := <-l.appendReqs:
			l.appendRequest(r)
		case r := <-l.flushReqs:
			l.flushRequest(r)
		}
	}
}

func (l *Log) Append(
	ctx context.Context,
	tid record.TransactionID,
	time time.Time,
	action record.Action,
	data []byte,
	reverse []byte,
) (record.LogSequenceNumber, error) {
	return l.appendReqs.Send(ctx, appendValue{
		tid:     tid,
		time:    time,
		action:  action,
		data:    data,
		reverse: reverse,
	})
}

func (l *Log) appendRequest(r *appendRequest) {
	v := r.Value()
	lsn, err := l.append(v.tid, v.time, v.action, v.data, v.reverse)
	r.Reply(lsn, err)
}

func (l *Log) append(
	tid record.TransactionID,
	time time.Time,
	action record.Action,
	data []byte,
	reverse []byte,
) (record.LogSequenceNumber, error) {
	h := l.highWater.Increment()
	rc := record.New(h, tid, time, action, data, reverse)
	if err := l.writer.Append(rc); err != nil {
		return record.NewLogSequenceNumber(0), err
	}

	l.highWater = h
	return l.highWater, nil
}

func (l *Log) Close() error {
	if l.cancel != nil {
		l.cancel()
	}

	return l.sl.Close()
}

func (l *Log) Flush(ctx context.Context, lsn record.LogSequenceNumber) error {
	return l.flushReqs.Send(ctx, lsn)
}

func (l *Log) flushRequest(r *flushRequest) {
	err := l.flush(r.Value())
	r.Reply(err)
}

func (l *Log) flush(lsn record.LogSequenceNumber) error {
	if err := l.writer.Flush(lsn); err != nil {
		return err
	}

	// TODO: Is this correct?
	lsn = record.NewLogSequenceNumber(min(lsn.Value(), l.highWater.Value()))
	l.lowWater = lsn
	return nil
}

func (l *Log) Iterate(lsn record.LogSequenceNumber) iter.Seq2[*record.Record, error] {
	return l.reader.Iterate(lsn)
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
	return l.reader.Recover()
}

func (l *Log) Reverse() iter.Seq2[*record.Record, error] {
	return l.reader.Reverse()
}

// TODO: Should this be merged with Open?
func (l *Log) StartWriter() error {
	return l.writer.Start()
}
