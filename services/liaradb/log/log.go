package log

import (
	"context"
	"iter"
	"time"

	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/log/action"
	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

type Log struct {
	sl         *segment.List
	reader     *reader
	writer     *writer
	highWater  record.LogSequenceNumber
	lowWater   record.LogSequenceNumber
	appendReqs chan *appendRequest
	flushReqs  chan *flushRequest
	cancel     context.CancelFunc
}

type flushRequest = async.Command[record.LogSequenceNumber]

type appendRequest = async.Request[appendValue, record.LogSequenceNumber]

type appendValue struct {
	tid     record.TransactionID
	time    time.Time
	action  action.Action
	data    []byte
	reverse []byte
}

func NewLog(
	pageSize int64,
	segmentSize page.PageID,
	fsys file.FileSystem,
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
func (l *Log) PageID() page.PageID                 { return l.writer.PageID() }

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
	action action.Action,
	data []byte,
	reverse []byte,
) (record.LogSequenceNumber, error) {
	req := async.NewRequest[appendValue, record.LogSequenceNumber](ctx, appendValue{
		tid:     tid,
		time:    time,
		action:  action,
		data:    data,
		reverse: reverse,
	})

	select {
	case l.appendReqs <- req:
	case <-ctx.Done():
		return 0, context.Canceled
	}

	return req.Wait(ctx)
}

func (l *Log) appendRequest(r *appendRequest) {
	v := r.Value()
	lsn, err := l.append(v.tid, v.time, v.action, v.data, v.reverse)
	r.Reply(lsn, err)
}

func (l *Log) append(
	tid record.TransactionID,
	time time.Time,
	action action.Action,
	data []byte,
	reverse []byte,
) (record.LogSequenceNumber, error) {
	rc := record.New(l.highWater+1, tid, time, action, data, reverse)
	if err := l.writer.Append(rc); err != nil {
		return 0, err
	}

	l.highWater++
	return l.highWater, nil
}

func (l *Log) Close() error {
	if l.cancel != nil {
		l.cancel()
	}

	return l.sl.Close()
}

func (l *Log) Flush(ctx context.Context, lsn record.LogSequenceNumber) error {
	req := async.NewCommand(lsn)

	select {
	case l.flushReqs <- req:
	case <-ctx.Done():
		return context.Canceled
	}

	return req.Wait(ctx)
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
	lsn = min(lsn, l.highWater)
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

func (l *Log) StartWriter() error {
	return l.writer.Start()
}
