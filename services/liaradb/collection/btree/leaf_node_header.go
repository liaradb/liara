package btree

import "github.com/liaradb/liaradb/storage"

type LeafNodeHeader[K Key] struct {
	Level       byte
	Slots       int16
	PrevBlockID storage.BlockID
	NextBlockID storage.BlockID
	HighKey     K
}
