package storage

import "github.com/liaradb/liaradb/storage/link"

type queryCurrent struct {
	fileName link.FileName
}

func newQueryCurrent(fn link.FileName) bufferQuery {
	return queryCurrent{
		fileName: fn,
	}
}

func (q queryCurrent) getBufferID(s *Storage) (link.FileBlockID, bool) {
	return s.highFileBlockID(q.fileName), false
}
