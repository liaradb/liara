package log

import (
	"context"
	"iter"
	"time"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

type Log struct {
	sl        *segment.List
	reader    *reader
	writer    *writer
	highWater record.LogSequenceNumber
	lowWater  record.LogSequenceNumber
	requests  chan *request
	cancel    context.CancelFunc
}

func NewLog(
	pageSize int64,
	segmentSize page.PageID,
	fsys file.FileSystem,
	dir string,
) *Log {
	sl := segment.NewList(fsys, dir)
	return &Log{
		sl:       sl,
		reader:   newReader(pageSize, sl),
		writer:   newWriter(pageSize, segmentSize, sl),
		requests: make(chan *request),
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
		case r := <-l.requests:
			l.request(r)
		}
	}
}

func (l *Log) request(r *request) {
	lsn, err := l.append(r.tid, r.time, r.action, r.data, r.reverse)
	r.Reply(lsn, err)
}

func (l *Log) Append(
	ctx context.Context,
	tid record.TransactionID,
	time time.Time,
	action record.Action,
	data []byte,
	reverse []byte,
) (record.LogSequenceNumber, error) {
	req := newRequest(tid, time, action, data, reverse)

	select {
	case l.requests <- req:
	case <-ctx.Done():
	}

	select {
	case res := <-req.response:
		return res.lsn, res.err
	case <-ctx.Done():
		return 0, context.Canceled
	}
}

func (l *Log) append(
	tid record.TransactionID,
	time time.Time,
	action record.Action,
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

func (l *Log) Flush(lsn record.LogSequenceNumber) error {
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
