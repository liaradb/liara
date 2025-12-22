package keyvalue

import (
	"bufio"
	"bytes"
	"context"
	"io"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/collection/tablename"
	domain "github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type KeyValue struct {
	s      *storage.Storage
	c      *btree.Cursor
	buffer *bytes.Buffer
	reader *bufio.Reader
}

func New(s *storage.Storage) *KeyValue {
	buffer := bytes.NewBuffer(nil)
	reader := bufio.NewReader(buffer)
	return &KeyValue{
		s:      s,
		c:      btree.NewCursor(s),
		buffer: buffer,
		reader: reader,
	}
}

func (kv *KeyValue) Get(ctx context.Context, fn string, key value.Key) ([]byte, error) {
	tn := tablename.New(fn)
	rid, err := kv.c.Search(ctx, tn.Index(0, domain.NewPartitionID(0)), key)
	if err != nil {
		return nil, err
	}

	b, err := kv.s.Request(ctx, link.NewBlockID(tn.KeyValue(domain.NewPartitionID(0)), page.Offset(rid.Block())))
	if err != nil {
		return nil, err
	}

	defer b.Release()

	p := NewBufferPage(b)
	// TODO: Find a simpler way
	i := 0
	for data, err := range p.Items() {
		if err != nil {
			return nil, err
		}

		if i == int(rid.Position()) {
			buf := raw.NewBufferFromSlice(data)
			var result []byte
			err := raw.Read(buf, &result)
			return result, err
		}

		i++
	}

	return nil, btree.ErrNotFound
}

func (kv *KeyValue) Set(ctx context.Context, fn string, key value.Key, v []byte) error {
	tn := tablename.New(fn)
	// TODO: Don't use io.Reader
	rid, err := kv.append(ctx, tn.KeyValue(domain.NewPartitionID(0)), v)
	if err != nil {
		return err
	}

	return kv.c.Insert(ctx,
		tn.Index(0, domain.NewPartitionID(0)),
		key,
		link.NewRecordLocator(
			rid.BlockID().Position,
			link.RecordPosition(rid.Position())))
}

func (l *KeyValue) append(ctx context.Context, fileName string, value []byte) (link.RecordID, error) {
	if err := raw.Write(l.buffer, value); err != nil {
		return link.RecordID{}, err
	}

	rid, err := l.appendData(ctx, fileName, l.reader)
	if err != nil {
		return link.RecordID{}, err
	}

	l.buffer.Reset()
	return rid, nil
}

// TODO: Should this be multiple BlockIDs?
func (kv *KeyValue) appendData(ctx context.Context, fn string, rd io.Reader) (link.RecordID, error) {
	// TODO: Find a better way to get this
	data := make([]byte, kv.s.BufferSize())
	n, err := rd.Read(data)
	if err != nil && err != io.EOF {
		return link.RecordID{}, err
	}

	bid, err := kv.appendCurrent(ctx, fn, data[:n])
	if err == raw.ErrInsufficientSpace {
		bid, err = kv.appendNext(ctx, fn, data[:n])
	}

	return bid, err
}

func (kv *KeyValue) appendCurrent(ctx context.Context, fn string, data []byte) (link.RecordID, error) {
	b, err := kv.s.RequestCurrent(ctx, fn)
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

func (kv *KeyValue) appendNext(ctx context.Context, fn string, data []byte) (link.RecordID, error) {
	b, err := kv.s.RequestNext(ctx, fn)
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
