package btree

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/btree/keynode"
	"github.com/liaradb/liaradb/collection/btree/leafnode"
	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
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
	fn link.FileName,
	k value.Key,
) (link.RecordLocator, error) {
	level, block, rid, err := c.searchRoot(ctx, fn, k)
	if err != nil {
		return link.RecordLocator{}, err
	}

	if level == 0 {
		return rid, nil
	}

	for i := level - 1; i > 0; i-- {
		_, block, err = c.searchKey(ctx, fn.BlockID(block), k)
		if err != nil {
			return link.RecordLocator{}, err
		}
	}

	return c.searchLeaf(ctx, fn.BlockID(block), k)
}

func (c *search) searchRoot(
	ctx context.Context,
	fn link.FileName,
	k value.Key,
) (byte, link.FilePosition, link.RecordLocator, error) {
	p, err := c.ns.getPage(ctx, fn.BlockID(0))
	if err != nil {
		return 0, 0, link.RecordLocator{}, err
	}

	if l := p.Level(); l == 0 {
		ln := leafnode.New(p)
		defer ln.Release()

		ln.RLatch()
		defer ln.RUnlatch()

		rid, ok := ln.Search(k)
		if !ok {
			return l, 0, link.RecordLocator{}, ErrNotFound
		}

		return l, 0, rid, nil
	} else {
		kn := keynode.New(p)
		defer kn.Release()

		kn.RLatch()
		defer kn.RUnlatch()

		return l, kn.Search(k), link.RecordLocator{}, nil
	}
}

func (c *search) searchKey(
	ctx context.Context,
	bid link.BlockID,
	k value.Key,
) (byte, link.FilePosition, error) {
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
	bid link.BlockID,
	k value.Key,
) (link.RecordLocator, error) {
	ln, err := c.ns.getLeafNode(ctx, bid)
	if err != nil {
		return link.RecordLocator{}, err
	}

	defer ln.Release()

	ln.RLatch()
	defer ln.RUnlatch()

	rid, ok := ln.Search(k)
	if !ok {
		return link.RecordLocator{}, ErrNotFound
	}

	return rid, nil
}

func (s *search) SearchRange(
	ctx context.Context,
	fn link.FileName,
	k value.Key,
	skip int,
	limit int,
) iter.Seq2[link.RecordLocator, error] {
	skipped := 0
	returned := 0
	return func(yield func(link.RecordLocator, error) bool) {
		block, rids, err := s.searchRangeFirst(ctx, fn, k)
		if err != nil {
			yield(link.RecordLocator{}, err)
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
				yield(link.RecordLocator{}, err)
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

func (s *search) searchRangeFirst(
	ctx context.Context,
	fn link.FileName,
	k value.Key,
) (link.FilePosition, iter.Seq[link.RecordLocator], error) {
	ln, err := s.searchRange(ctx, fn, k)
	if err != nil {
		return 0, nil, err
	}

	defer ln.Release()

	ln.RLatch()
	defer ln.RUnlatch()

	return ln.RightID(), ln.SearchRange(k), nil
}

func (s *search) searchRangeNext(
	ctx context.Context,
	fn link.FileName,
	block link.FilePosition,
) (link.FilePosition, iter.Seq[link.RecordLocator], error) {
	ln, err := s.ns.getLeafNode(ctx, fn.BlockID(block))
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
	fn link.FileName,
	k value.Key,
) (*leafnode.LeafNode, error) {
	level, block, ln, err := c.searchRangeRoot(ctx, fn, k)
	if err != nil {
		return nil, err
	}

	if level == 0 {
		return ln, nil
	}

	for i := level - 1; i > 0; i-- {
		_, block, err = c.searchRangeKey(ctx, fn.BlockID(block), k)
		if err != nil {
			return nil, err
		}
	}

	return c.ns.getLeafNode(ctx, fn.BlockID(block))
}

func (c *search) searchRangeRoot(
	ctx context.Context,
	fn link.FileName,
	k value.Key,
) (byte, link.FilePosition, *leafnode.LeafNode, error) {
	p, err := c.ns.getPage(ctx, fn.BlockID(0))
	if err != nil {
		return 0, 0, nil, err
	}

	if l := p.Level(); l == 0 {
		return 0, 0, leafnode.New(p), nil
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
	bid link.BlockID,
	k value.Key,
) (byte, link.FilePosition, error) {
	kn, err := c.ns.getKeyNode(ctx, bid)
	if err != nil {
		return 0, 0, err
	}

	defer kn.Release()

	return kn.Level(), kn.Search(k), nil
}

func (s *search) All(
	ctx context.Context,
	fn link.FileName,
	skip int,
	limit int,
) iter.Seq2[link.RecordLocator, error] {
	skipped := 0
	returned := 0
	return func(yield func(link.RecordLocator, error) bool) {
		block, rids, err := s.allFirst(ctx, fn)
		if err != nil {
			yield(link.RecordLocator{}, err)
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

			block, rids, err = s.allNext(ctx, fn, block)
			if err != nil {
				yield(link.RecordLocator{}, err)
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

func (s *search) allFirst(
	ctx context.Context,
	fn link.FileName,
) (link.FilePosition, iter.Seq[link.RecordLocator], error) {
	ln, err := s.all(ctx, fn)
	if err != nil {
		return 0, nil, err
	}

	defer ln.Release()

	ln.RLatch()
	defer ln.RUnlatch()

	return ln.RightID(), ln.RecordIDs(), nil
}

func (s *search) allNext(
	ctx context.Context,
	fn link.FileName,
	block link.FilePosition,
) (link.FilePosition, iter.Seq[link.RecordLocator], error) {
	ln, err := s.ns.getLeafNode(ctx, fn.BlockID(block))
	if err != nil {
		return 0, nil, err
	}

	defer ln.Release()

	ln.RLatch()
	defer ln.RUnlatch()

	return ln.RightID(), ln.RecordIDs(), nil
}

func (c *search) all(
	ctx context.Context,
	fn link.FileName,
) (*leafnode.LeafNode, error) {
	level, block, ln, err := c.allRoot(ctx, fn)
	if err != nil {
		return nil, err
	}

	if level == 0 {
		return ln, nil
	}

	for i := level - 1; i > 0; i-- {
		_, block, err = c.allKey(ctx, fn.BlockID(block))
		if err != nil {
			return nil, err
		}
	}

	return c.ns.getLeafNode(ctx, fn.BlockID(block))
}

func (c *search) allRoot(
	ctx context.Context,
	fn link.FileName,
) (byte, link.FilePosition, *leafnode.LeafNode, error) {
	p, err := c.ns.getPage(ctx, fn.BlockID(0))
	if err != nil {
		return 0, 0, nil, err
	}

	if l := p.Level(); l == 0 {
		return l, 0, leafnode.New(p), nil
	} else {
		kn := keynode.New(p)
		defer kn.Release()

		kn.RLatch()
		defer kn.RUnlatch()

		_, fp, ok := kn.Child(0)
		if !ok {
			return 0, 0, nil, ErrNotFound
		}

		return l, fp, nil, nil
	}
}

func (c *search) allKey(
	ctx context.Context,
	bid link.BlockID,
) (byte, link.FilePosition, error) {
	kn, err := c.ns.getKeyNode(ctx, bid)
	if err != nil {
		return 0, 0, err
	}

	defer kn.Release()

	_, fp, ok := kn.Child(0)
	if !ok {
		return 0, 0, ErrNotFound
	}

	return kn.Level(), fp, nil
}
