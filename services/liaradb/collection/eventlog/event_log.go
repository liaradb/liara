package eventlog

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/storage/node"
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
	c, err := rd.Read(data)
	if err != nil && err != io.EOF {
		return link.RecordID{}, err
	}

	v := data[:c]
	crc := page.NewCRC(v)

	rid, ok, err := l.setCurrent(ctx, fn, v, crc)
	if err != nil {
		return link.RecordID{}, err
	} else if !ok {
		rid, ok, err = l.setNext(ctx, fn, v, crc)
		if err != nil {
			return link.RecordID{}, err
		} else if !ok {
			return link.RecordID{}, btree.ErrNoInsert
		}
	}

	return link.NewRecordID(link.NewBlockID(fn, rid.Block()), link.RecordPosition(rid.Position())), nil
}

func (l *EventLog) setCurrent(ctx context.Context, fn link.FileName, v []byte, crc page.CRC) (link.RecordLocator, bool, error) {
	b, err := l.storage.RequestCurrent(ctx, fn)
	if err != nil {
		return link.RecordLocator{}, false, err
	}

	defer b.Release()

	n := node.New(b)
	rp, d, ok := n.Append(int16(len(v)), crc)
	if !ok {
		return link.RecordLocator{}, false, nil
	}

	copy(d, v)

	return link.NewRecordLocator(b.BlockID().Position(), rp), true, nil
}

func (l *EventLog) setNext(ctx context.Context, fn link.FileName, v []byte, crc page.CRC) (link.RecordLocator, bool, error) {
	b, err := l.storage.RequestNext(ctx, fn)
	if err != nil {
		return link.RecordLocator{}, false, err
	}

	defer b.Release()

	n := node.New(b)
	rp, d, ok := n.Append(int16(len(v)), crc)
	if !ok {
		return link.RecordLocator{}, false, nil
	}

	copy(d, v)

	return link.NewRecordLocator(b.BlockID().Position(), rp), true, nil
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
		for n, err := range l.Iterate(ctx, fn) {
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(n, nil) {
				return
			}
		}
	}
}

func (l *EventLog) Iterate(ctx context.Context, fn link.FileName) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		highBid, err := l.storage.Highwater(ctx, fn)
		if err != nil {
			yield(nil, err)
			return
		}

		bid := fn.BlockID(0)
		for bid.Position() <= highBid.Position() {
			children, p, err := l.handleIteration(ctx, bid)
			if err != nil {
				yield(nil, err)
				return
			}

			for c := range children {
				if !yield(c, nil) {
					return
				}
			}

			bid.SetPosition(p)
		}
	}
}

func (l *EventLog) handleIteration(ctx context.Context, bid link.BlockID) (iter.Seq[[]byte], link.FilePosition, error) {
	b, err := l.storage.Request(ctx, bid)
	if err != nil {
		return nil, 0, err
	}

	return func(yield func([]byte) bool) {
		defer b.Release()
		for c := range node.New(b).Children() {
			if !yield(c) {
				return
			}
		}
	}, bid.Position() + 1, nil
}
