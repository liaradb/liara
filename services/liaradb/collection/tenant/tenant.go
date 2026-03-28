package tenant

import (
	"context"
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

type Tenant struct {
	s *storage.Storage
	c *btree.Cursor
}

func New(storage *storage.Storage, cursor *btree.Cursor) *Tenant {
	return &Tenant{
		s: storage,
		c: cursor,
	}
}

// TODO: Use io.Reader?
func (o *Tenant) Get(
	ctx context.Context,
	tn tablename.TableName,
	tid value.TenantID,
) (*entity.Tenant, error) {
	k := key.NewKey(tid.Bytes())
	fnIdx := tn.Index(0, value.NewPartitionID(0))
	rid, err := o.c.Search(ctx, fnIdx, k)
	if err != nil {
		return nil, err
	}

	return o.getItem(ctx, tn, rid)
}

func (o *Tenant) List(
	ctx context.Context,
	tn tablename.TableName,
) iter.Seq2[*entity.Tenant, error] {
	return func(yield func(*entity.Tenant, error) bool) {
		fnIdx := tn.Index(0, value.PartitionID{})
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

func (o *Tenant) getItem(ctx context.Context, tn tablename.TableName, rid link.RecordLocator) (*entity.Tenant, error) {
	bid := tn.KeyValue(value.NewPartitionID(0)).BlockID(rid.Block())
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

	e := &entity.Tenant{}
	_, _ = e.Read(d)
	return e, nil
}

// TODO: Use io.Writer?
func (o *Tenant) Set(
	ctx context.Context,
	tn tablename.TableName,
	tid value.TenantID,
	e *entity.Tenant,
) error {
	fn := tn.KeyValue(value.NewPartitionID(0))
	k := key.NewKey(tid.Bytes())

	v := make([]byte, entity.TenantSize)
	_, _ = e.Write(v)
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

func (o *Tenant) setCurrent(ctx context.Context, fn link.FileName, v []byte, crc page.CRC) (link.RecordLocator, bool, error) {
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

func (o *Tenant) setNext(ctx context.Context, fn link.FileName, v []byte, crc page.CRC) (link.RecordLocator, bool, error) {
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
func (o *Tenant) Replace(
	ctx context.Context,
	tn tablename.TableName,
	tid value.TenantID,
	e *entity.Tenant,
) error {
	k := key.NewKey(tid.Bytes())
	fnIdx := tn.Index(0, value.NewPartitionID(0))
	rid, err := o.c.Search(ctx, fnIdx, k)
	if err != nil {
		return err
	}

	bid := tn.KeyValue(value.NewPartitionID(0)).BlockID(rid.Block())
	b, err := o.s.Request(ctx, bid)
	if err != nil {
		return err
	}

	defer b.Release()

	n := node.New(b)

	if !n.IsPage() {
		return page.ErrNotPage
	}

	v := make([]byte, entity.TenantSize)
	if _, ok := e.Write(v); !ok {
		return btree.ErrNoUpdate
	}

	if !n.ReplaceChild(int16(rid.Position()), v) {
		return btree.ErrNoUpdate
	}

	return nil
}
