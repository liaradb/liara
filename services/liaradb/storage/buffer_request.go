package storage

import (
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/async"
)

type bufferRequest = async.Request[bufferQuery, *Buffer]

type bufferQuery interface {
	getBufferID(*Storage) (link.FileBlockID, bool)
}
