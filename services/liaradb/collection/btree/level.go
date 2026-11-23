package btree

import (
	"context"

	"github.com/liaradb/liaradb/collection/log"
	"github.com/liaradb/liaradb/storage/link"
)

type level struct {
	ns *nodeStorage
}

func newLevel(l log.Log) level {
	return level{
		ns: newNodeStorage(l),
	}
}

func (l *level) Level(ctx context.Context, fn link.FileName) (byte, error) {
	p, err := l.ns.getPage(ctx, fn.FileBlockID(0))
	if err != nil {
		return 0, err
	}

	defer p.Release()

	p.RLatch()
	defer p.RUnlatch()

	return p.Level(), nil
}
