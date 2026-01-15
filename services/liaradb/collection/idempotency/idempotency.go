package idempotency

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

type Idempotency struct {
	s *storage.Storage
	c *btree.Cursor
}

func New(storage *storage.Storage, cursor *btree.Cursor) *Idempotency {
	return &Idempotency{
		s: storage,
		c: cursor,
	}
}

// TODO: Use io.Reader?
func (i *Idempotency) Get(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
	rqid value.RequestID,
) (*entity.RequestLog, error) {
	k := key.NewKey(rqid.Bytes())
	fnIdx := tn.Index(0, pid)
	rid, err := i.c.Search(ctx, fnIdx, k)
	if err != nil {
		return nil, err
	}

	return i.getItem(ctx, tn, rid)
}

func (i *Idempotency) List(
	ctx context.Context,
	tn tablename.TableName,
	pid value.PartitionID,
) iter.Seq2[*entity.RequestLog, error] {
	return func(yield func(*entity.RequestLog, error) bool) {
		fnIdx := tn.Index(0, pid)
		for rid, err := range i.c.All(ctx, fnIdx, 0, 0) {
			if err != nil {
				yield(nil, err)
				return
			}

			i, err := i.getItem(ctx, tn, rid)
			if !yield(i, err) {
				return
			}
		}
	}
}

func (i *Idempotency) getItem(ctx context.Context, tn tablename.TableName, rid link.RecordLocator) (*entity.RequestLog, error) {
	bid := tn.RequestLog().BlockID(rid.Block())
	b, err := i.s.Request(ctx, bid)
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
	e := &entity.RequestLog{}
	_ = e.Read(d)
	return e, nil
}

// TODO: Use io.Writer?
func (i *Idempotency) Set(
	ctx context.Context,
	tn tablename.TableName,
	rqid value.RequestID,
	e *entity.RequestLog,
) error {
	fn := tn.RequestLog()
	k := key.NewKey(rqid.Bytes())

	v := make([]byte, entity.RequestLogSize)
	_ = e.Write(v)
	crc := page.NewCRC(v)

	rid, ok, err := i.setCurrent(ctx, fn, v, crc)
	if err != nil {
		return err
	} else if !ok {
		rid, ok, err = i.setNext(ctx, fn, v, crc)
		if err != nil {
			return err
		} else if !ok {
			return btree.ErrNoInsert
		}
	}

	fnIdx := tn.Index(0, value.NewPartitionID(0))
	return i.c.Insert(ctx, fnIdx, k, rid)
}

func (i *Idempotency) setCurrent(ctx context.Context, fn link.FileName, v []byte, crc page.CRC) (link.RecordLocator, bool, error) {
	b, err := i.s.RequestCurrent(ctx, fn)
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

func (i *Idempotency) setNext(ctx context.Context, fn link.FileName, v []byte, crc page.CRC) (link.RecordLocator, bool, error) {
	b, err := i.s.RequestNext(ctx, fn)
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
