package btree

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/page"
	"github.com/liaradb/liaradb/storage"
)

// TODO: Create latching support
type Cursor struct {
	s *storage.Storage
}

func NewCursor(s *storage.Storage) *Cursor {
	return &Cursor{
		s: s,
	}
}

func (c *Cursor) Insert(
	ctx context.Context,
	fileName string,
	k Key,
	rid RecordID,
) error {
	return c.insertPage(ctx,
		storage.NewBlockID(fileName, 0),
		k,
		rid)
}

func (c *Cursor) insertPage(
	ctx context.Context,
	bid storage.BlockID,
	k Key,
	rid RecordID,
) error {
	page, err := c.GetPage(ctx, bid)
	if err != nil {
		return err
	}

	if page.Level() == 0 {
		_, _, err := c.insertLeaf(ctx, bid.FileName, page, k, rid)
		return err
	} else {
		return c.insertKey(ctx, bid.FileName, page, k, rid)
	}
}

func (c *Cursor) insertLeaf(
	ctx context.Context,
	fileName string,
	p page.BTreePage,
	k Key,
	rid RecordID,
) (storage.BlockID, bool, error) {
	ln := NewLeafNode(p)
	first, second, ok := ln.Insert(k, rid)
	if ok {
		return storage.BlockID{}, false, nil
	}

	b, err := c.s.RequestNext(ctx, fileName)
	if err != nil {
		return storage.BlockID{}, false, err
	}

	p2 := page.New(b)
	ln2 := NewLeafNode(p2)

	ln2.Fill(second)
	ln.Replace(first)

	return b.BlockID(), true, nil
}

func (c *Cursor) insertKey(
	ctx context.Context,
	fileName string,
	p page.BTreePage,
	k Key,
	rid RecordID,
) error {
	kn := newKeyNode(p)
	block := kn.Search(k)

	return c.insertPage(ctx,
		storage.NewBlockID(fileName, storage.Offset(block)),
		k,
		rid)
}

func (c *Cursor) Search(
	ctx context.Context,
	fileName string,
	k Key,
) (RecordID, error) {
	return c.searchPage(ctx, storage.NewBlockID(fileName, 0), k)
}

func (c *Cursor) searchPage(
	ctx context.Context,
	bid storage.BlockID,
	k Key,
) (RecordID, error) {
	p, err := c.GetPage(ctx, bid)
	if err != nil {
		return RecordID{}, err
	}

	if p.Level() == 0 {
		return c.searchLeaf(p, k)
	} else {
		return c.searchKey(ctx, bid.FileName, p, k)
	}
}

func (*Cursor) searchLeaf(
	p page.BTreePage,
	k Key,
) (RecordID, error) {
	ln := NewLeafNode(p)
	rid, ok := ln.Search(k)
	if !ok {
		return RecordID{}, ErrNotFound
	}

	return rid, nil
}

func (c *Cursor) searchKey(
	ctx context.Context,
	fileName string,
	p page.BTreePage,
	k Key,
) (RecordID, error) {
	kn := newKeyNode(p)
	block := kn.Search(k)
	return c.searchPage(ctx,
		storage.NewBlockID(fileName, storage.Offset(block)),
		k)
}

func (c *Cursor) GetRoot(ctx context.Context, fileName string) (page.BTreePage, error) {
	return c.GetPage(ctx, storage.NewBlockID(fileName, 0))
}

func (c *Cursor) GetPage(ctx context.Context, bid storage.BlockID) (page.BTreePage, error) {
	b, err := c.GetBuffer(ctx, bid)
	if err != nil {
		return page.BTreePage{}, err
	}

	return page.New(b), nil
}

func (c *Cursor) GetBuffer(ctx context.Context, bid storage.BlockID) (*storage.Buffer, error) {
	return c.s.Request(ctx, bid)
}
