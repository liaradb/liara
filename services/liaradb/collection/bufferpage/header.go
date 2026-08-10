package bufferpage

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/wrap"
)

const (
	headerSize = 0 +
		page.MagicSize
)

type header struct {
	magic wrap.Int32
}

func newHeader(data []byte) (header, []byte) {
	magic, data0 := wrap.NewInt32(data)

	return header{
		magic: magic,
	}, data0
}

func (h *header) init() {
	h.magic.Set(int32(page.MagicPage))
}

func (h *header) isEmpty() bool {
	return page.Magic(h.magic.Get()).IsEmpty()
}

func (h *header) isPage() bool {
	return page.Magic(h.magic.Get()).IsPage()
}
