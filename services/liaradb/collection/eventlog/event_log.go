package eventlog

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"iter"

	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/storage"
)

type EventLog struct {
	storage *storage.Storage
	buffer  *bytes.Buffer
	reader  *bufio.Reader
}

func New(
	storage *storage.Storage,
) *EventLog {
	buffer := bytes.NewBuffer(nil)
	reader := bufio.NewReader(buffer)
	return &EventLog{
		storage: storage,
		buffer:  buffer,
		reader:  reader,
	}
}

func (l *EventLog) Append(ctx context.Context, fileName string, e *entity.Event) error {
	if err := e.Write(l.buffer); err != nil {
		return err
	}

	_, err := l.AppendEvent(ctx, fileName, l.reader)
	if err != nil {
		return err
	}

	l.buffer.Reset()
	return nil
}

// TODO: Should this be multiple BlockIDs?
func (l *EventLog) AppendEvent(ctx context.Context, fileName string, rd io.Reader) (storage.BlockID, error) {
	// TODO: Find a better way to get this
	data := make([]byte, l.storage.BufferSize())
	n, err := rd.Read(data)
	if err != nil && err != io.EOF {
		return storage.BlockID{}, err
	}

	bid, err := l.appendCurrent(ctx, fileName, data[:n])
	if err == raw.ErrInsufficientSpace {
		bid, err = l.appendNext(ctx, fileName, data[:n])
	}

	return bid, err
}

func (l *EventLog) appendCurrent(ctx context.Context, fileName string, data []byte) (storage.BlockID, error) {
	b, err := l.storage.RequestCurrent(ctx, fileName)
	if err != nil {
		return storage.BlockID{}, err
	}

	defer b.Release()

	bp := storage.NewBufferPage(b)
	if err := bp.Add(data); err != nil {
		return storage.BlockID{}, err
	}

	return b.BlockID(), nil
}

func (l *EventLog) appendNext(ctx context.Context, fileName string, data []byte) (storage.BlockID, error) {
	b, err := l.storage.RequestNext(ctx, fileName)
	if err != nil {
		return storage.BlockID{}, err
	}

	defer b.Release()

	bp := storage.NewBufferPage(b)
	if err := bp.Add(data); err != nil {
		return storage.BlockID{}, err
	}

	return b.BlockID(), nil
}

func (l *EventLog) Find(ctx context.Context, fn string, id value.EventID) (*entity.Event, error) {
	for e, err := range l.Events(ctx, fn) {
		if err != nil {
			return nil, err
		}

		if e.ID == id {
			return e, nil
		}
	}

	return nil, errors.New("not found")
}

func (l *EventLog) GetAggregate(ctx context.Context, fn string, id value.AggregateID) iter.Seq2[*entity.Event, error] {
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

func (l *EventLog) Events(ctx context.Context, fn string) iter.Seq2[*entity.Event, error] {
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

func (l *EventLog) items(ctx context.Context, fn string) iter.Seq2[[]byte, error] {
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

func (l *EventLog) Iterate(ctx context.Context, fn string) iter.Seq2[*storage.BufferPage, error] {
	return func(yield func(*storage.BufferPage, error) bool) {
		highBid, err := l.storage.Highwater(ctx, fn)
		if err != nil {
			yield(nil, err)
			return
		}

		bid := storage.NewBlockID(fn, 0)
		for bid.Position <= highBid.Position {
			p, ok := l.handleIteration(ctx, bid, yield)
			if !ok {
				return
			}

			bid.Position = p
		}
	}
}

func (l *EventLog) handleIteration(ctx context.Context, bid storage.BlockID, yield func(*storage.BufferPage, error) bool) (storage.Offset, bool) {
	b, err := l.storage.Request(ctx, bid)
	if err != nil {
		yield(nil, err)
		return bid.Position, false
	}

	defer b.Release()
	if !yield(storage.NewBufferPage(b), err) {
		return bid.Position, false
	}

	return bid.Position + 1, true
}
