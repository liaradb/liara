package lrupool

import (
	"iter"

	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type LFUPool struct {
}

var _ storage.FreePool = (*LFUPool)(nil)

func New() *LFUPool {
	return &LFUPool{}
}

func (p *LFUPool) Count() int {
	panic("unimplemented")
}

func (p *LFUPool) Iterate() iter.Seq[*storage.Buffer] {
	panic("unimplemented")
}

func (p *LFUPool) Pop() (*storage.Buffer, bool) {
	panic("unimplemented")
}

func (p *LFUPool) Push(k link.BlockID, v *storage.Buffer) {
	panic("unimplemented")
}

func (p *LFUPool) Remove(k link.BlockID) (*storage.Buffer, bool) {
	panic("unimplemented")
}
