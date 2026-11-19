package storage

import "github.com/liaradb/liaradb/storage/link"

type queryNext struct {
	fileName link.FileName
}

func newQueryNext(fn link.FileName) bufferQuery {
	return queryNext{
		fileName: fn,
	}
}

func (q queryNext) getBufferID(hw highWater) (link.FileBlockID, bool) {
	return hw.incrementHighFileBlockID(q.fileName), true
}
