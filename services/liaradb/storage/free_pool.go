package storage

import (
	"iter"

	"github.com/liaradb/liaradb/storage/link"
)

type FreePool interface {
	Count() int
	Iterate() iter.Seq[*Buffer]
	Pop() (*Buffer, bool)
	Push(k link.FileBlockID, v *Buffer)
	Remove(k link.FileBlockID) (*Buffer, bool)
}
