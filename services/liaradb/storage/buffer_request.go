package storage

import (
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/async"
)

type bufferRequest = async.Request[bufferQuery, *Buffer]

type bufferQuery interface {
	getBufferID(highWater) (link.FileBlockID, bool)
}

type highWater interface {
	highFileBlockID(link.FileName) link.FileBlockID
	incrementHighFileBlockID(link.FileName) link.FileBlockID
}
