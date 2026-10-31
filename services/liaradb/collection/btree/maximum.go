package btree

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/btree/keynode"
	"github.com/liaradb/liaradb/collection/btree/leafnode"
	"github.com/liaradb/liaradb/collection/btree/node"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type maximum struct {
	ns *nodeStorage
}

func newMaximum(s *storage.Storage) maximum {
	return maximum{
		ns: newNodeStorage(s),
	}
}

func (m *maximum) Maximum(ctx context.Context, fn link.FileName) (key.Key, error) {
	return m.getKey(ctx, fn, 0)
}

func (m *maximum) getKey(ctx context.Context, fn link.FileName, fp link.BlockID) (key.Key, error) {
	p, err := m.ns.getPage(ctx, fn.FileBlockID(fp))
	if err != nil {
		return key.Key{}, err
	}

	defer p.Release()

	p.RLatch()
	defer p.RUnlatch()

	if p.Level() > 0 {
		return m.getKeyFromKeyNode(ctx, p, fn)
	}

	return m.getKeyFromLeafNode(p)
}

func (*maximum) getKeyFromLeafNode(p node.Node) (key.Key, error) {
	if p.Count() == 0 {
		return key.Key{}, ErrNotFound
	}

	ln := leafnode.New(p)
	le := ln.Child(p.Count() - 1)
	return le.Key(), nil
}

func (m *maximum) getKeyFromKeyNode(ctx context.Context, p node.Node, fn link.FileName) (key.Key, error) {
	kn := keynode.New(p)
	_, fp := kn.Child(p.Count() - 1)
	return m.getKey(ctx, fn, fp)
}
