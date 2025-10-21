package eventlog

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"iter"

	"github.com/liaradb/liaradb/raw"
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

func (l *EventLog) Append(ctx context.Context, fileName string, e *Event) error {
	var data []byte
	if _, err := l.buffer.Write(data); err != nil {
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

func (s *EventLog) handleIteration(ctx context.Context, bid storage.BlockID, yield func(*storage.Buffer, error) bool) (raw.Offset, bool) {
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
