package storage

import "github.com/liaradb/liaradb/async"

type bufferRequest = async.Request[bufferQuery, *Buffer]

type bufferQuery struct {
	bid       BlockID
	fileName  string
	queryType bufferQueryType
}

type bufferQueryType int

const (
	bufferQueryTypeByID = iota
	bufferQueryTypeCurrent
	bufferQueryTypeNext
)

func newBufferByIDQuery(bid BlockID) bufferQuery {
	return bufferQuery{
		bid:       bid,
		queryType: bufferQueryTypeByID,
	}
}

func newCurrentBufferQuery(fileName string) bufferQuery {
	return bufferQuery{
		fileName:  fileName,
		queryType: bufferQueryTypeCurrent,
	}
}

func newNextBufferQuery(fileName string) bufferQuery {
	return bufferQuery{
		fileName:  fileName,
		queryType: bufferQueryTypeNext,
	}
}
