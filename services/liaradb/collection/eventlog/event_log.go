package eventlog

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"iter"

	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type EventLog struct {
	storage *storage.Storage
	buffer  *bytes.Buffer
	reader  *bufio.Reader
}

func New(storage *storage.Storage) *EventLog {
	buffer := bytes.NewBuffer(nil)
	reader := bufio.NewReader(buffer)
	return &EventLog{
		storage: storage,
		buffer:  buffer,
		reader:  reader,
	}
}

func (l *EventLog) Append(ctx context.Context, fn link.FileName, e *entity.Event) (link.RecordID, error) {
	if err := e.Write(l.buffer); err != nil {
		return link.RecordID{}, err
	}

	rid, err := l.AppendEvent(ctx, fn, l.reader)
	if err != nil {
		return link.RecordID{}, err
	}

	l.buffer.Reset()
	return rid, nil
}

// TODO: Should this be multiple BlockIDs?
func (l *EventLog) AppendEvent(ctx context.Context, fn link.FileName, rd io.Reader) (link.RecordID, error) {
	// TODO: Find a better way to get this
	data := make([]byte, l.storage.BufferSize())
	n, err := rd.Read(data)
	if err != nil && err != io.EOF {
		return link.RecordID{}, err
	}

	bid, err := l.appendCurrent(ctx, fn, data[:n])
	if err == raw.ErrInsufficientSpace {
		bid, err = l.appendNext(ctx, fn, data[:n])
	}

	return bid, err
}

func (l *EventLog) appendCurrent(ctx context.Context, fn link.FileName, data []byte) (link.RecordID, error) {
	b, err := l.storage.RequestCurrent(ctx, fn)
	if err != nil {
		return link.RecordID{}, err
	}

	defer b.Release()

	bp := NewBufferPage(b)
	offset, err := bp.Add(data)
	if err != nil {
		return link.RecordID{}, err
	}

	// TODO: Fix this type
	return b.BlockID().RecordID(link.RecordPosition(offset)), nil
}

func (l *EventLog) appendNext(ctx context.Context, fn link.FileName, data []byte) (link.RecordID, error) {
	b, err := l.storage.RequestNext(ctx, fn)
	if err != nil {
		return link.RecordID{}, err
	}

	defer b.Release()

	bp := NewBufferPage(b)
	offset, err := bp.Add(data)
	if err != nil {
		return link.RecordID{}, err
	}

	// TODO: Fix this type
	return b.BlockID().RecordID(link.RecordPosition(offset)), nil
}

func (l *EventLog) Find(ctx context.Context, fn link.FileName, id value.EventID) (*entity.Event, error) {
	for e, err := range l.Events(ctx, fn) {
		if err != nil {
			return nil, err
		}

		if e.ID == id {
			return e, nil
		}
	}

	return nil, page.ErrNotFound
}

func (l *EventLog) GetAggregate(ctx context.Context, fn link.FileName, id value.AggregateID) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		for e, err := range l.Events(ctx, fn) {
			if err != nil {
				yield(nil, err)
				return
			}

			if e.AggregateID == id {
				if !yield(e, nil) {
					return
				}
			}
		}
	}
}

func (l *EventLog) Events(ctx context.Context, fn link.FileName) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		for i, err := range l.items(ctx, fn) {
			if err != nil {
				yield(nil, err)
				return
			}

			// TODO: Optimize this
			buf := bytes.NewBuffer(i)

			var e entity.Event
			if err := e.Read(buf); err != nil {
				yield(nil, err)
				return
			}

			if !yield(&e, nil) {
				return
			}
		}
	}
}

func (l *EventLog) items(ctx context.Context, fn link.FileName) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		for b, err := range l.Iterate(ctx, fn) {
			if err != nil {
				yield(nil, err)
				return
			}

			for i, err := range b.Items() {
				if err != nil {
					yield(nil, err)
					return
				}

				if !yield(i, nil) {
					return
				}
			}
		}
	}
}

func (l *EventLog) Iterate(ctx context.Context, fn link.FileName) iter.Seq2[*BufferPage, error] {
	return func(yield func(*BufferPage, error) bool) {
		highBid, err := l.storage.Highwater(ctx, fn)
		if err != nil {
			yield(nil, err)
			return
		}

		bid := fn.BlockID(0)
		for bid.Position() <= highBid.Position() {
			p, ok := l.handleIteration(ctx, bid, yield)
			if !ok {
				return
			}

			bid.SetPosition(p)
		}
	}
}

func (l *EventLog) handleIteration(ctx context.Context, bid link.BlockID, yield func(*BufferPage, error) bool) (link.FilePosition, bool) {
	b, err := l.storage.Request(ctx, bid)
	if err != nil {
		yield(nil, err)
		return bid.Position(), false
	}

	defer b.Release()
	if !yield(NewBufferPage(b), err) {
		return bid.Position(), false
	}

	return bid.Position() + 1, true
}
