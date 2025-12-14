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

// Insert key value pair into tree
func (c *Cursor) Insert(
	ctx context.Context,
	fileName string,
	k Key,
	rid RecordID,
) error {
	_, _, err := c.insertPage(ctx, storage.NewBlockID(fileName, 0), k, rid)
	return err
}

// Get page from buffer pool and insert key value pair
func (c *Cursor) insertPage(
	ctx context.Context,
	bid storage.BlockID,
	k Key,
	rid RecordID,
) (storage.BlockID, bool, error) {
	page, err := c.GetPage(ctx, bid)
	if err != nil {
		return storage.BlockID{}, false, err
	}

	if page.Level() == 0 {
		return c.insertLeaf(ctx, bid.FileName, page, k, rid)
	}

	return c.insertKey(ctx, bid.FileName, page, k, rid)
}

// This is a leaf level page.
//   - Insert, and handle a split.
func (c *Cursor) insertLeaf(
	ctx context.Context,
	fileName string,
	p page.BTreePage,
	k Key,
	rid RecordID,
) (storage.BlockID, bool, error) {
	ln := NewLeafNode(p)
	first, second, split := ln.Insert(k, rid)
	if split {
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

// This is a key level page.
//   - Find the correct child.
//   - Insert into the child.
//   - If child splits, insert into this page.
//   - Handle a split
func (c *Cursor) insertKey(
	ctx context.Context,
	fileName string,
	p page.BTreePage,
	k Key,
	rid RecordID,
) (storage.BlockID, bool, error) {
	kn := newKeyNode(p)
	childID := kn.Search(k)

	b2, split, err := c.insertPage(ctx,
		storage.NewBlockID(fileName, storage.Offset(childID)),
		k,
		rid)
	if err != nil {
		return storage.BlockID{}, false, err
	} else if !split {
		return storage.BlockID{}, false, nil
	}

	first, second, split := kn.Insert(k, BlockPosition(b2.Position))
	if !split {
		return storage.BlockID{}, false, nil
	}

	b, err := c.s.RequestNext(ctx, fileName)
	if err != nil {
		return storage.BlockID{}, false, err
	}

	p2 := page.New(b)
	kn2 := newKeyNode(p2)

	kn2.Fill(second)
	kn.Replace(first)

	return b.BlockID(), true, nil
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
