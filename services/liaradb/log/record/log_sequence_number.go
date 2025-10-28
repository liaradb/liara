package record

import "github.com/liaradb/liaradb/encoder/raw"

type LogSequenceNumber struct {
	baseUint64
}

const LogSequenceNumberSize = raw.BaseUint64Size

func NewLogSequenceNumber(value uint64) LogSequenceNumber {
	return LogSequenceNumber{raw.NewBaseUint64(value)}
}

func (l LogSequenceNumber) Increment() LogSequenceNumber {
	return LogSequenceNumber{l.baseUint64 + 1}
}

func (l LogSequenceNumber) Decrement() LogSequenceNumber {
	return LogSequenceNumber{l.baseUint64 - 1}
}
