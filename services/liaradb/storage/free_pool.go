package storage

import (
	"iter"

	"github.com/liaradb/liaradb/storage/link"
)

type FreePool interface {
	Count() int
	Iterate() iter.Seq[*Buffer]
	Pop() (*Buffer, bool)
	Push(k link.BlockID, v *Buffer)
	Remove(k link.BlockID) (*Buffer, bool)
}
