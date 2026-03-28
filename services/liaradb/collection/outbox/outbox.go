package outbox

import (
	"context"
	"io"
	"iter"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
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
func (o *Outbox) Get(
	ctx context.Context,
	tn tablename.TableName,
	oid value.OutboxID,
) (*entity.Outbox, error) {
	k := key.NewKey(oid.Bytes())
	fnIdx := tn.Index(0, value.NewPartitionID(0))
	rid, err := o.c.Search(ctx, fnIdx, k)
	if err != nil {
		return nil, err
	}

	return o.getItem(ctx, tn, rid)
}

func (o *Outbox) List(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
) iter.Seq2[*entity.Outbox, error] {
	return func(yield func(*entity.Outbox, error) bool) {
		fnIdx := tn.Index(0, pid)
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

func (o *Outbox) getItem(ctx context.Context, tn tablename.TableName, rid link.RecordLocator) (*entity.Outbox, error) {
	bid := tn.Outbox(value.NewPartitionID(0)).BlockID(rid.Block())
	b, err := o.s.Request(ctx, bid)
	if err != nil {
		return nil, err
	}

	defer b.Release()

	n := node.New(b)

	if !n.IsPage() {
		return nil, page.ErrNotPage
	}

	d, ok := n.Child(rid.Position())
	if !ok {
		return nil, btree.ErrNotFound
	}

	e := &entity.Outbox{}
	if _, ok := e.Read(d); !ok {
		return nil, io.EOF
	}

	return e, nil
}

// TODO: Use io.Writer?
func (o *Outbox) Set(
	ctx context.Context,
	tn tablename.TableName,
	oid value.OutboxID,
	e *entity.Outbox,
) error {
	fn := tn.Outbox(value.NewPartitionID(0))
	k := key.NewKey(oid.Bytes())

	v := make([]byte, entity.OutboxSize)
	if _, ok := e.Write(v); !ok {
		return io.EOF
	}

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

	fnIdx := tn.Index(0, value.NewPartitionID(0))
	return o.c.Insert(ctx, fnIdx, k, rid)
}

func (o *Outbox) setCurrent(ctx context.Context, fn link.FileName, v []byte, crc page.CRC) (link.RecordLocator, bool, error) {
	b, err := o.s.RequestCurrent(ctx, fn)
	if err != nil {
		return link.RecordLocator{}, false, err
	}

	defer b.Release()

	n := node.New(b)

	if !n.IsPage() {
		return link.RecordLocator{}, false, page.ErrNotPage
	}

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

// TODO: Use io.Writer?
func (o *Outbox) Replace(
	ctx context.Context,
	tn tablename.TableName,
	oid value.OutboxID,
	e *entity.Outbox,
) error {
	k := key.NewKey(oid.Bytes())
	fnIdx := tn.Index(0, value.NewPartitionID(0))
	rid, err := o.c.Search(ctx, fnIdx, k)
	if err != nil {
		return err
	}

	bid := tn.Outbox(value.NewPartitionID(0)).BlockID(rid.Block())
	b, err := o.s.Request(ctx, bid)
	if err != nil {
		return err
	}

	defer b.Release()

	n := node.New(b)

	if n.IsPage() {
		return page.ErrNotPage
	}

	v := make([]byte, entity.OutboxSize)
	if _, ok := e.Write(v); !ok {
		return btree.ErrNoUpdate
	}

	if !n.ReplaceChild(int16(rid.Position()), v) {
		return btree.ErrNoUpdate
	}

	return nil
}
