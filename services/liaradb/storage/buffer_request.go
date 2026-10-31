package storage

import (
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/async"
)

type bufferRequest = async.Request[bufferQuery, *Buffer]

type bufferQuery struct {
	fbid      link.FileBlockID
	fileName  link.FileName
	queryType bufferQueryType
}

type bufferQueryType int

const (
	bufferQueryTypeByID = iota
	bufferQueryTypeCurrent
	bufferQueryTypeNext
)

func newBufferByIDQuery(fbid link.FileBlockID) bufferQuery {
	return bufferQuery{
		fbid:      fbid,
		queryType: bufferQueryTypeByID,
	}
}

func newCurrentBufferQuery(fn link.FileName) bufferQuery {
	return bufferQuery{
		fileName:  fn,
		queryType: bufferQueryTypeCurrent,
	}
}

func newNextBufferQuery(fn link.FileName) bufferQuery {
	return bufferQuery{
		fileName:  fn,
		queryType: bufferQueryTypeNext,
	}
}

func (q bufferQuery) isNext() bool {
	return q.queryType == bufferQueryTypeNext
}
