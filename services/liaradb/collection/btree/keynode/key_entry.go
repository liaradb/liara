package keynode

import (
	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/encoder/wrap"
)

type keyEntry struct {
	key   value.Key
	block value.BlockPosition
}

func (ke keyEntry) Key() value.Key             { return ke.key }
func (ke keyEntry) Block() value.BlockPosition { return ke.block }

func newKeyEntry(key value.Key, block value.BlockPosition) keyEntry {
	return keyEntry{
		key:   key,
		block: block,
	}
}

func (ke keyEntry) Size() int { return ke.key.Size() + value.BlockPositionSize }

func (ke keyEntry) Write(data []byte) {
	block, data0 := wrap.NewInt64(data)
	block.Set(ke.block.Value())
	copy(data0, []byte(ke.key))
}

func (ke *keyEntry) Read(data []byte) {
	block, data0 := wrap.NewInt64(data)
	ke.block = value.BlockPosition(block.Get())
	ke.key = value.Key(data0)
}
