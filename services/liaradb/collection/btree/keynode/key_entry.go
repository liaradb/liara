package keynode

import (
	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/encoder/scan"
	"github.com/liaradb/liaradb/storage/link"
)

type keyEntry struct {
	key   value.Key
	block link.FilePosition
}

func (ke keyEntry) Key() value.Key           { return ke.key }
func (ke keyEntry) Block() link.FilePosition { return ke.block }

func newKeyEntry(key value.Key, block link.FilePosition) keyEntry {
	return keyEntry{
		key:   key,
		block: block,
	}
}

func (ke keyEntry) Size() int { return ke.key.Size() + link.FilePositionSize }

func (ke keyEntry) Write(data []byte) {
	data0 := scan.SetInt64(data, ke.block.Value())
	ke.key.Write(data0)
}

func (ke *keyEntry) Read(data []byte) {
	block, data0 := scan.Int64(data)
	ke.block = link.FilePosition(block)
	ke.key.Read(data0)
}
