package btree

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/btree/keynode"
	"github.com/liaradb/liaradb/collection/btree/leafnode"
	"github.com/liaradb/liaradb/storage"
)

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

	if l := p.Level(); l == 0 {
		ln := leafnode.New(p)
		defer ln.Release()

		ln.RLatch()
		defer ln.RUnlatch()

		rid, ok := ln.Search(k)
		if !ok {
			return l, 0, leafnode.RecordID{}, ErrNotFound
		}

		return l, 0, rid, nil
	} else {
		kn := keynode.New(p)
		defer kn.Release()

		kn.RLatch()
		defer kn.RUnlatch()

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

	kn.RLatch()
	defer kn.RUnlatch()

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

	ln.RLatch()
	defer ln.RUnlatch()

	rid, ok := ln.Search(k)
	if !ok {
		return leafnode.RecordID{}, ErrNotFound
	}

	return rid, nil
}

func (s *search) SearchRange(
	ctx context.Context,
	fn string,
	k key.Key,
	skip int,
	limit int,
) iter.Seq2[leafnode.RecordID, error] {
	skipped := 0
	returned := 0
	return func(yield func(leafnode.RecordID, error) bool) {
		block, rids, err := s.searchRangeFirst(ctx, fn, k)
		if err != nil {
			yield(leafnode.RecordID{}, err)
			return
		}

		for rid := range rids {
			if skip > skipped {
				skipped++
				continue
			}
			if s.isLimit(limit, returned) || !yield(rid, nil) {
				return
			}
			returned++
		}

		for block != 0 {
			if s.isLimit(limit, returned) {
				return
			}

			block, rids, err = s.searchRangeNext(ctx, fn, block)
			if err != nil {
				yield(leafnode.RecordID{}, err)
				return
			}

			for rid := range rids {
				if skip > skipped {
					skipped++
					continue
				}
				if s.isLimit(limit, returned) || !yield(rid, nil) {
					return
				}
				returned++
			}
		}
	}
}

func (s *search) isLimit(limit int, returned int) bool {
	return limit > 0 && returned >= limit
}

func (s *search) searchRangeFirst(ctx context.Context, fn string, k key.Key) (keynode.BlockPosition, iter.Seq[leafnode.RecordID], error) {
	ln, err := s.searchRange(ctx, fn, k)
	if err != nil {
		return 0, nil, err
	}

	defer ln.Release()

	ln.RLatch()
	defer ln.RUnlatch()

	return ln.RightID(), ln.SearchRange(k), nil
}

func (s *search) searchRangeNext(ctx context.Context, fn string, block keynode.BlockPosition) (keynode.BlockPosition, iter.Seq[leafnode.RecordID], error) {
	ln, err := s.ns.getLeafNode(ctx, storage.NewBlockID(fn, storage.Offset(block)))
	if err != nil {
		return 0, nil, err
	}

	defer ln.Release()

	ln.RLatch()
	defer ln.RUnlatch()

	return ln.RightID(), ln.RecordIDs(), nil
}

func (c *search) searchRange(
	ctx context.Context,
	fn string,
	k key.Key,
) (*leafnode.LeafNode, error) {
	level, block, ln, err := c.searchRangeRoot(ctx, fn, k)
	if err != nil {
		return nil, err
	}

	if level == 0 {
		return ln, nil
	}

	for i := level - 1; i > 0; i-- {
		_, block, err = c.searchRangeKey(ctx,
			storage.NewBlockID(fn, storage.Offset(block)), k)
		if err != nil {
			return nil, err
		}
	}

	return c.ns.getLeafNode(ctx,
		storage.NewBlockID(fn, storage.Offset(block)))
}

func (c *search) searchRangeRoot(
	ctx context.Context,
	fn string,
	k key.Key,
) (byte, keynode.BlockPosition, *leafnode.LeafNode, error) {
	p, err := c.ns.getPage(ctx, storage.NewBlockID(fn, 0))
	if err != nil {
		return 0, 0, nil, err
	}

	if l := p.Level(); l == 0 {
		return l, 0, leafnode.New(p), ErrNotFound
	} else {
		kn := keynode.New(p)
		defer kn.Release()

		kn.RLatch()
		defer kn.RUnlatch()

		return l, kn.Search(k), nil, nil
	}
}

func (c *search) searchRangeKey(
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
