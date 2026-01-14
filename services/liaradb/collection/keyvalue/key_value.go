package keyvalue

import (
	"context"
	"iter"
	"slices"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/storage/node"
)

// TODO: Create latching
type KeyValue struct {
	s *storage.Storage
	c *btree.Cursor
}

func New(s *storage.Storage, c *btree.Cursor) *KeyValue {
	return &KeyValue{
		s: s,
		c: c,
	}
}

// TODO: Use io.Reader?
func (kv *KeyValue) Get(ctx context.Context, tn tablename.TableName, key key.Key) ([]byte, error) {
	fnIdx := tn.Index(0, value.NewPartitionID(0))
	rid, err := kv.c.Search(ctx, fnIdx, key)
	if err != nil {
		return nil, err
	}

	return kv.getItem(ctx, tn, rid)
}

func (kv *KeyValue) List(ctx context.Context, tn tablename.TableName) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		fnIdx := tn.Index(0, value.NewPartitionID(0))
		for rid, err := range kv.c.All(ctx, fnIdx, 0, 0) {
			if err != nil {
				yield(nil, err)
				return
			}

			i, err := kv.getItem(ctx, tn, rid)
			if !yield(i, err) {
				return
			}
		}
	}
}

func (kv *KeyValue) getItem(ctx context.Context, tn tablename.TableName, rid link.RecordLocator) ([]byte, error) {
	bid := tn.KeyValue(value.NewPartitionID(0)).BlockID(rid.Block())
	b, err := kv.s.Request(ctx, bid)
	if err != nil {
		return nil, err
	}

	defer b.Release()

	n := node.New(b)
	// TODO: Fix this type
	d, ok := n.Child(int16(rid.Position()))
	if !ok {
		return nil, btree.ErrNotFound
	}

	// TODO: Should we clone?
	return slices.Clone(d), nil
}

// TODO: Use io.Writer?
func (kv *KeyValue) Set(ctx context.Context, tn tablename.TableName, key key.Key, v []byte) error {
	fn := tn.KeyValue(value.NewPartitionID(0))
	crc := page.NewCRC(v)

	rid, ok, err := kv.setCurrent(ctx, fn, v, crc)
	if err != nil {
		return err
	} else if !ok {
		rid, ok, err = kv.setNext(ctx, fn, v, crc)
		if err != nil {
			return err
		} else if !ok {
			return btree.ErrNoInsert
		}
	}

	fnIdx := tn.Index(0, value.NewPartitionID(0))
	return kv.c.Insert(ctx, fnIdx, key, rid)
}

func (kv *KeyValue) setCurrent(ctx context.Context, fn link.FileName, v []byte, crc page.CRC) (link.RecordLocator, bool, error) {
	b, err := kv.s.RequestCurrent(ctx, fn)
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

func (kv *KeyValue) setNext(ctx context.Context, fn link.FileName, v []byte, crc page.CRC) (link.RecordLocator, bool, error) {
	b, err := kv.s.RequestNext(ctx, fn)
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
