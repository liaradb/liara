package btree

import (
	"context"

	"github.com/liaradb/liaradb/storage"
)

// TODO: Create latching support
type level struct {
	ns *nodeStorage
}

func newLevel(s *storage.Storage) level {
	return level{
		ns: newNodeStorage(s),
	}
}

func (l *level) Level(ctx context.Context, fn string) (byte, error) {
	p, err := l.ns.getPage(ctx, storage.NewBlockID(fn, 0))
	if err != nil {
		return 0, err
	}

	defer p.Release()

	p.RLatch()
	defer p.RUnlatch()

	return p.Level(), nil
}
