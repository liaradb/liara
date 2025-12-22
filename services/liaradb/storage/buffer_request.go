package storage

import (
	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/storage/link"
)

type bufferRequest = async.Request[bufferQuery, *Buffer]

type bufferQuery struct {
	bid       link.BlockID
	fileName  link.FileName
	queryType bufferQueryType
}

type bufferQueryType int

const (
	bufferQueryTypeByID = iota
	bufferQueryTypeCurrent
	bufferQueryTypeNext
)

func newBufferByIDQuery(bid link.BlockID) bufferQuery {
	return bufferQuery{
		bid:       bid,
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
