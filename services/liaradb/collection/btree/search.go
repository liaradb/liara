package btree

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/btree/keynode"
	"github.com/liaradb/liaradb/collection/btree/leafnode"
	"github.com/liaradb/liaradb/storage"
)

// TODO: Create latching support
type search struct {
	ns *nodeStorage
}

func newSearch(s *storage.Storage) search {
	return search{
		ns: newNodeStorage(s),
	}
}

func (c *search) Search(
	ctx context.Context,
	fn string,
	k key.Key,
) (leafnode.RecordID, error) {
	level, block, rid, err := c.searchRoot(ctx, fn, k)
	if err != nil {
		return leafnode.RecordID{}, err
	}

	if level == 0 {
		return rid, nil
	}

	for i := level - 1; i > 0; i-- {
		_, block, err = c.searchKey(ctx,
			storage.NewBlockID(fn, storage.Offset(block)), k)
		if err != nil {
			return leafnode.RecordID{}, err
		}
	}

	return c.searchLeaf(ctx,
		storage.NewBlockID(fn, storage.Offset(block)), k)
}

func (c *search) searchRoot(
	ctx context.Context,
	fn string,
	k key.Key,
) (byte, keynode.BlockPosition, leafnode.RecordID, error) {
	p, err := c.ns.getPage(ctx, storage.NewBlockID(fn, 0))
	if err != nil {
		return 0, 0, leafnode.RecordID{}, err
	}

	l := p.Level()
	if l == 0 {
		ln := leafnode.New(p)
		defer ln.Release()
		rid, ok := ln.Search(k)
		if !ok {
			return l, 0, leafnode.RecordID{}, ErrNotFound
		}

		return l, 0, rid, nil
	} else {
		kn := keynode.New(p)
		defer kn.Release()
		return l, kn.Search(k), leafnode.RecordID{}, nil
	}
}

func (c *search) searchKey(
	ctx context.Context,
	bid storage.BlockID,
	k key.Key,
) (byte, keynode.BlockPosition, error) {
	kn, err := c.ns.getKeyNode(ctx, bid)
	if err != nil {
		return 0, 0, err
	}

	defer kn.Release()

	return kn.Level(), kn.Search(k), nil
}

func (c *search) searchLeaf(
	ctx context.Context,
	bid storage.BlockID,
	k key.Key,
) (leafnode.RecordID, error) {
	ln, err := c.ns.getLeafNode(ctx, bid)
	if err != nil {
		return leafnode.RecordID{}, err
	}

	defer ln.Release()

	rid, ok := ln.Search(k)
	if !ok {
		return leafnode.RecordID{}, ErrNotFound
	}

	return rid, nil
}

func (s *search) SearchRange(ctx context.Context, fn string, k key.Key) iter.Seq2[leafnode.RecordID, error] {
	return func(yield func(leafnode.RecordID, error) bool) {
		p, err := s.ns.getPage(ctx, storage.NewBlockID(fn, 0))
		if err != nil {
			yield(leafnode.RecordID{}, err)
			return
		}

		defer p.Release()
	}
}
