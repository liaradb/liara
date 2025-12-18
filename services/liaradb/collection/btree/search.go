package btree

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/btree/keynode"
	"github.com/liaradb/liaradb/collection/btree/leafnode"
	"github.com/liaradb/liaradb/collection/btree/page"
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

func (c *search) SearchRecursive(
	ctx context.Context,
	fn string,
	k key.Key,
) (leafnode.RecordID, error) {
	return c.searchPage(ctx, storage.NewBlockID(fn, 0), k)
}

func (c *search) searchPage(
	ctx context.Context,
	bid storage.BlockID,
	k key.Key,
) (leafnode.RecordID, error) {
	p, err := c.ns.getPage(ctx, bid)
	if err != nil {
		return leafnode.RecordID{}, err
	}

	defer p.Release()

	if p.Level() == 0 {
		return c.searchLeaf(p, k)
	} else {
		return c.searchKey(ctx, bid.FileName, p, k)
	}
}

func (*search) searchLeaf(p page.Page, k key.Key) (leafnode.RecordID, error) {
	ln := leafnode.New(p)
	rid, ok := ln.Search(k)
	if !ok {
		return leafnode.RecordID{}, ErrNotFound
	}

	return rid, nil
}

func (c *search) searchKey(
	ctx context.Context,
	fn string,
	p page.Page,
	k key.Key,
) (leafnode.RecordID, error) {
	bid := storage.NewBlockID(fn, storage.Offset(keynode.New(p).Search(k)))
	return c.searchPage(ctx, bid, k)
}

func (c *search) Search(
	ctx context.Context,
	fn string,
	k key.Key,
) (leafnode.RecordID, error) {
	level, block, rid, err := c.searchRootNode(ctx, fn, k)
	if err != nil {
		return leafnode.RecordID{}, err
	}

	if level == 0 {
		return rid, nil
	}

	for i := level - 1; i > 0; i-- {
		_, block, err = c.searchKeyNode(ctx,
			storage.NewBlockID(fn, storage.Offset(block)), k)
		if err != nil {
			return leafnode.RecordID{}, err
		}
	}

	return c.searchLeafNode(ctx,
		storage.NewBlockID(fn, storage.Offset(block)), k)
}

func (c *search) searchRootNode(
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

func (c *search) searchKeyNode(
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

func (c *search) searchLeafNode(
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
