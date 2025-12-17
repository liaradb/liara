package btree

import (
	"context"

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

func (c *search) Search(
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
