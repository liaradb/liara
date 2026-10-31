package keynode

import (
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/encoder/scan"
	"github.com/liaradb/liaradb/storage/link"
)

type keyEntry struct {
	key   key.Key
	block link.BlockID
}

func (ke keyEntry) Key() key.Key        { return ke.key }
func (ke keyEntry) Block() link.BlockID { return ke.block }

func newKeyEntry(key key.Key, block link.BlockID) keyEntry {
	return keyEntry{
		key:   key,
		block: block,
	}
}

func (ke keyEntry) Size() int { return ke.key.Size() + link.BlockIDSize }

func (ke keyEntry) Write(data []byte) {
	data0 := scan.SetInt64(data, ke.block.Value())
	ke.key.Write(data0)
}

func (ke *keyEntry) Read(data []byte) {
	block, data0 := scan.Int64(data)
	ke.block = link.BlockID(block)
	ke.key.Read(data0)
}
