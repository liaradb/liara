package btree

import (
	"context"

	"github.com/liaradb/liaradb/collection/btree/key"
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

func (l *maximum) Maximum(ctx context.Context, fn link.FileName) (key.Key, error) {
	p, err := l.ns.getPage(ctx, fn.BlockID(0))
	if err != nil {
		return key.Key{}, err
	}

	defer p.Release()

	p.RLatch()
	defer p.RUnlatch()

	// TODO: Implement this
	return key.Key{}, ErrNotFound
}
