package outbox

import (
	"context"
	"iter"
	"slices"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/tablename"
	domain "github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/storage/node"
)

type Outbox struct {
	s *storage.Storage
	c *btree.Cursor
}

func New(storage *storage.Storage, cursor *btree.Cursor) *Outbox {
	return &Outbox{
		s: storage,
		c: cursor,
	}
}

// TODO: Use io.Reader?
func (o *Outbox) Get(ctx context.Context, tn tablename.TableName, key key.Key) ([]byte, error) {
	fnIdx := tn.Index(0, domain.NewPartitionID(0))
	rid, err := o.c.Search(ctx, fnIdx, key)
	if err != nil {
		return nil, err
	}

	return o.getItem(ctx, tn, rid)
}

// TODO: Test this
func (o *Outbox) List(ctx context.Context, tn tablename.TableName) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		fnIdx := tn.Index(0, domain.NewPartitionID(0))
		for rid, err := range o.c.All(ctx, fnIdx, 0, 0) {
			if err != nil {
				yield(nil, err)
				return
			}

			i, err := o.getItem(ctx, tn, rid)
			if !yield(i, err) {
				return
			}
		}
	}
}

func (o *Outbox) getItem(ctx context.Context, tn tablename.TableName, rid link.RecordLocator) ([]byte, error) {
	bid := tn.KeyValue(domain.NewPartitionID(0)).BlockID(rid.Block())
	b, err := o.s.Request(ctx, bid)
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
func (o *Outbox) Set(ctx context.Context, tn tablename.TableName, key key.Key, v []byte) error {
	fn := tn.KeyValue(domain.NewPartitionID(0))
	crc := page.NewCRC(v)

	rid, ok, err := o.setCurrent(ctx, fn, v, crc)
	if err != nil {
		return err
	} else if !ok {
		rid, ok, err = o.setNext(ctx, fn, v, crc)
		if err != nil {
			return err
		} else if !ok {
			return btree.ErrNoInsert
		}
	}

	fnIdx := tn.Index(0, domain.NewPartitionID(0))
	return o.c.Insert(ctx, fnIdx, key, rid)
}

func (o *Outbox) setCurrent(ctx context.Context, fn link.FileName, v []byte, crc page.CRC) (link.RecordLocator, bool, error) {
	b, err := o.s.RequestCurrent(ctx, fn)
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

func (o *Outbox) setNext(ctx context.Context, fn link.FileName, v []byte, crc page.CRC) (link.RecordLocator, bool, error) {
	b, err := o.s.RequestNext(ctx, fn)
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
