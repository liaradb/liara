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
type searchCursor struct {
	ns *nodeStorage
}

func newSearchCursor(s *storage.Storage) searchCursor {
	return searchCursor{
		ns: newNodeStorage(s),
	}
}

func (c *searchCursor) Search(
	ctx context.Context,
	fn string,
	k key.Key,
) (leafnode.RecordID, error) {
	return c.searchPage(ctx, storage.NewBlockID(fn, 0), k)
}

func (c *searchCursor) searchPage(
	ctx context.Context,
	bid storage.BlockID,
	k key.Key,
) (leafnode.RecordID, error) {
	p, err := c.GetPage(ctx, bid)
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

func (*searchCursor) searchLeaf(p page.BTreePage, k key.Key) (leafnode.RecordID, error) {
	ln := leafnode.New(p)
	rid, ok := ln.Search(k)
	if !ok {
		return leafnode.RecordID{}, ErrNotFound
	}

	return rid, nil
}

func (c *searchCursor) searchKey(
	ctx context.Context,
	fn string,
	p page.BTreePage,
	k key.Key,
) (leafnode.RecordID, error) {
	bid := storage.NewBlockID(fn, storage.Offset(keynode.New(p).Search(k)))
	return c.searchPage(ctx, bid, k)
}

func (c *searchCursor) GetPage(ctx context.Context, bid storage.BlockID) (page.BTreePage, error) {
	return c.ns.getPage(ctx, bid)
}
