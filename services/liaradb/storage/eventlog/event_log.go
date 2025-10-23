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
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/page"
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
func (s *EventLog) AppendEvent(ctx context.Context, fileName string, rd io.Reader) (storage.BlockID, error) {
	// TODO: Find a better way to get this
	data := make([]byte, s.storage.BufferSize())
	n, err := rd.Read(data)
	if err != nil && err != io.EOF {
		return storage.BlockID{}, err
	}

	bid, err := s.appendCurrent(ctx, fileName, data[:n])
	if err == page.ErrInsufficientSpace {
		bid, err = s.appendNext(ctx, fileName, data[:n])
	}

	return bid, err
}

func (s *EventLog) appendCurrent(ctx context.Context, fileName string, data []byte) (storage.BlockID, error) {
	b, err := s.storage.RequestCurrent(ctx, fileName)
	if err != nil {
		return storage.BlockID{}, err
	}

	defer b.Release()

	if err := b.Add(data); err != nil {
		return storage.BlockID{}, err
	}

	return b.BlockID(), nil
}

func (s *EventLog) appendNext(ctx context.Context, fileName string, data []byte) (storage.BlockID, error) {
	b, err := s.storage.RequestNext(ctx, fileName)
	if err != nil {
		return storage.BlockID{}, err
	}

	defer b.Release()

	if err := b.Add(data); err != nil {
		return storage.BlockID{}, err
	}

	return b.BlockID(), nil
}

func (s *EventLog) Find(ctx context.Context, fn string, id value.EventID) (*entity.Event, error) {
	for e, err := range s.Events(ctx, fn) {
		if err != nil {
			return nil, err
		}

		if e.ID == id {
			return e, nil
		}
	}

	return nil, errors.New("not found")
}

func (s *EventLog) Events(ctx context.Context, fn string) iter.Seq2[*entity.Event, error] {
	return func(yield func(*entity.Event, error) bool) {
		for i, err := range s.items(ctx, fn) {
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

func (s *EventLog) items(ctx context.Context, fn string) iter.Seq2[page.Item, error] {
	return func(yield func(page.Item, error) bool) {
		for b, err := range s.Iterate(ctx, fn) {
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

func (s *EventLog) Iterate(ctx context.Context, fn string) iter.Seq2[*storage.Buffer, error] {
	return func(yield func(*storage.Buffer, error) bool) {
		highBid, err := s.storage.Highwater(ctx, fn)
		if err != nil {
			yield(nil, err)
			return
		}

		bid := storage.NewBlockID(fn, 0)
		for bid.Position <= highBid.Position {
			p, ok := s.handleIteration(ctx, bid, yield)
			if !ok {
				return
			}

			bid.Position = p
		}
	}
}

func (s *EventLog) handleIteration(ctx context.Context, bid storage.BlockID, yield func(*storage.Buffer, error) bool) (storage.Offset, bool) {
	b, err := s.storage.Request(ctx, bid)
	if err != nil {
		yield(nil, err)
		return bid.Position, false
	}

	defer b.Release()
	if !yield(b, err) {
		return bid.Position, false
	}

	return bid.Position + 1, true
}
