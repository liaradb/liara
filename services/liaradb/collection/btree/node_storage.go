package btree

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/keynode"
	"github.com/liaradb/liaradb/collection/btree/leafnode"
	"github.com/liaradb/liaradb/collection/btree/node"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type nodeStorage struct {
	s *storage.Storage
}

func newNodeStorage(s *storage.Storage) *nodeStorage {
	return &nodeStorage{s}
}

func (ns *nodeStorage) getBuffer(ctx context.Context, bid link.BlockID) (*storage.Buffer, error) {
	return ns.s.Request(ctx, bid)
}

func (ns *nodeStorage) getKeyNode(ctx context.Context, bid link.BlockID) (*keynode.KeyNode, error) {
	b, err := ns.s.Request(ctx, bid)
	if err != nil {
		return nil, err
	}

	return keynode.New(node.New(b)), nil
}

func (ns *nodeStorage) getLeafNode(ctx context.Context, bid link.BlockID) (*leafnode.LeafNode, error) {
	b, err := ns.s.Request(ctx, bid)
	if err != nil {
		return nil, err
	}

	return leafnode.New(node.New(b)), nil
}

func (ns *nodeStorage) getNextBuffer(ctx context.Context, fn link.FileName) (*storage.Buffer, error) {
	return ns.s.RequestNext(ctx, fn)
}

func (ns *nodeStorage) getNextKeyNode(ctx context.Context, fn link.FileName) (*keynode.KeyNode, link.BlockID, error) {
	b, err := ns.s.RequestNext(ctx, fn)
	if err != nil {
		return nil, link.BlockID{}, err
	}

	return keynode.New(node.New(b)), b.BlockID(), nil
}

func (ns *nodeStorage) getNextLeafNode(ctx context.Context, fn link.FileName) (*leafnode.LeafNode, link.BlockID, error) {
	b, err := ns.s.RequestNext(ctx, fn)
	if err != nil {
		return nil, link.BlockID{}, err
	}

	return leafnode.New(node.New(b)), b.BlockID(), nil
}

func (ns *nodeStorage) getPage(ctx context.Context, bid link.BlockID) (node.Node, error) {
	b, err := ns.getBuffer(ctx, bid)
	if err != nil {
		return node.Node{}, err
	}

	return node.New(b), nil
}
