package btree

import "github.com/liaradb/liaradb/storage"

type KeyNode[K Key] struct {
	Header KeyNodeHeader[K]
	Keys   []KeyEntry[K]
	Values []storage.BlockID
}
