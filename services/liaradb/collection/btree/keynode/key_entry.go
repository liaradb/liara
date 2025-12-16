package keynode

import (
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/encoder/wrap"
)

// TODO: Test this
type KeyEntry struct {
	key   key.Key
	block BlockPosition
}

func (ke KeyEntry) Key() key.Key         { return ke.key }
func (ke KeyEntry) Block() BlockPosition { return ke.block }

func newKeyEntry(key key.Key, block BlockPosition) KeyEntry {
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
	ke.key = key.Key(data0)
}
