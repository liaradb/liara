package keynode

import (
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/encoder/scan"
	"github.com/liaradb/liaradb/storage/link"
)

type keyEntry struct {
	key   key.Key
	block link.FilePosition
}

func (ke keyEntry) Key() key.Key             { return ke.key }
func (ke keyEntry) Block() link.FilePosition { return ke.block }

func newKeyEntry(key key.Key, block link.FilePosition) keyEntry {
	return keyEntry{
		key:   key,
		block: block,
	}
}

func (ke keyEntry) Size() int { return ke.key.Size() + link.FilePositionSize }

func (ke keyEntry) Write(data []byte) bool {
	data0, ok := scan.SetInt64(data, ke.block.Value())
	if !ok {
		return false
	}

	return ke.key.Write(data0)
}

func (ke *keyEntry) Read(data []byte) bool {
	if len(data) < 8 {
		return false
	}

	block, data0, ok := scan.Int64(data)
	if !ok {
		return false
	}

	ke.block = link.FilePosition(block)
	return ke.key.Read(data0)
}
