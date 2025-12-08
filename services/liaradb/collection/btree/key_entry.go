package btree

import (
	"github.com/liaradb/liaradb/encoder/wrap"
)

// TODO: Test this
type KeyEntry struct {
	key   Key
	block BlockPosition
}

func newKeyEntry(key Key, block BlockPosition) KeyEntry {
	return KeyEntry{
		key:   key,
		block: block,
	}
}

func (ke KeyEntry) Size() int { return ke.key.Size() + BlockPositionSize }

func (ke KeyEntry) Write(data []byte) {
	block, data0 := wrap.NewInt64(data)
	block.Set(ke.block.Value())
	copy(data0, []byte(ke.key))
}

func (ke *KeyEntry) Read(data []byte) {
	block, data0 := wrap.NewInt64(data)
	ke.block = BlockPosition(block.Get())
	ke.key = Key(data0)
}
