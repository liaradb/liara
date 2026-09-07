package bufferpage

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/recovery/logpage"
)

const (
	headerSize = 0 +
		page.MagicSize +
		logpage.LogSequenceNumberSize
)

type header struct {
	magic wrap.Int32
	lsn   logpage.LogSequenceNumber
}

func newHeader(data []byte) (header, []byte) {
	magic, data0 := wrap.NewInt32(data)
	var lsn logpage.LogSequenceNumber
	data1, _ := lsn.ReadData(data0)

	return header{
		magic: magic,
		lsn:   lsn,
	}, data1
}

// TODO: Make sure to call this
func (h *header) init() {
	h.magic.Set(int32(page.MagicPage))
}

func (h *header) isEmpty() bool {
	return page.Magic(h.magic.Get()).IsEmpty()
}

// TODO: We are not using this
func (h *header) isPage() bool {
	return page.Magic(h.magic.Get()).IsPage()
}

func (h *header) LogSequenceNumber() logpage.LogSequenceNumber { return h.lsn }
