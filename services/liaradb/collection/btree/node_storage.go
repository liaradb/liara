package btree

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/keynode"
	"github.com/liaradb/liaradb/collection/btree/leafnode"
	"github.com/liaradb/liaradb/collection/btree/page"
	"github.com/liaradb/liaradb/storage"
)

type nodeStorage struct {
	s *storage.Storage
}

func newNodeStorage(s *storage.Storage) *nodeStorage {
	return &nodeStorage{s}
}

func (ns *nodeStorage) getBuffer(ctx context.Context, bid storage.BlockID) (*storage.Buffer, error) {
	return ns.s.Request(ctx, bid)
}

func (ns *nodeStorage) getKeyNode(ctx context.Context, bid storage.BlockID) (*keynode.KeyNode, error) {
	b, err := ns.s.Request(ctx, bid)
	if err != nil {
		return nil, err
	}

	return keynode.New(page.New(b)), nil
}

func (ns *nodeStorage) getLeafNode(ctx context.Context, bid storage.BlockID) (*leafnode.LeafNode, error) {
	b, err := ns.s.Request(ctx, bid)
	if err != nil {
		return nil, err
	}

	return leafnode.New(page.New(b)), nil
}

func (ns *nodeStorage) getNextBuffer(ctx context.Context, fn string) (*storage.Buffer, error) {
	return ns.s.RequestNext(ctx, fn)
}

func (ns *nodeStorage) getNextKeyNode(ctx context.Context, fn string) (*keynode.KeyNode, storage.BlockID, error) {
	b, err := ns.s.RequestNext(ctx, fn)
	if err != nil {
		return nil, storage.BlockID{}, err
	}

	return keynode.New(page.New(b)), b.BlockID(), nil
}

func (ns *nodeStorage) getNextLeafNode(ctx context.Context, fn string) (*leafnode.LeafNode, storage.BlockID, error) {
	b, err := ns.s.RequestNext(ctx, fn)
	if err != nil {
		return nil, storage.BlockID{}, err
	}

	return leafnode.New(page.New(b)), b.BlockID(), nil
}

func (ns *nodeStorage) getPage(ctx context.Context, bid storage.BlockID) (page.Page, error) {
	b, err := ns.getBuffer(ctx, bid)
	if err != nil {
		return page.Page{}, err
	}

	return page.New(b), nil
}
