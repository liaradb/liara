package storage

import "github.com/liaradb/liaradb/storage/link"

type queryByID struct {
	fbid link.FileBlockID
}

func newQueryByID(fbid link.FileBlockID) bufferQuery {
	return queryByID{
		fbid: fbid,
	}
}

func (q queryByID) getBufferID(hw highWater) (link.FileBlockID, bool) {
	return q.fbid, false
}
