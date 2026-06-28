package logpage

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

type LogPage struct {
	*page.Page
	lsn record.LogSequenceNumber
}

func (lp *LogPage) Clear(lsn record.LogSequenceNumber) {
	lp.Page.Clear()
	lp.lsn = lsn
}

func NewLogPage(size int16, headerSize int16, slotHeaderSize int16) *LogPage {
	return &LogPage{
		Page: page.New(size, headerSize, slotHeaderSize),
	}
}
