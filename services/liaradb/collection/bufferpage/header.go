package bufferpage

import (
	"github.com/liaradb/liaradb/recovery/logpage"
)

const (
	headerSize = 0 +
		logpage.LogSequenceNumberSize
)

type header struct {
	lsn logpage.LogSequenceNumber
}

func newHeader(data []byte) (header, []byte) {
	var lsn logpage.LogSequenceNumber
	data0 := lsn.ReadData(data)

	return header{
		lsn: lsn,
	}, data0
}

func (h *header) LogSequenceNumber() logpage.LogSequenceNumber       { return h.lsn }
func (h *header) SetLogSequenceNumber(lsn logpage.LogSequenceNumber) { h.lsn = lsn }
