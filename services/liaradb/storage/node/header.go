package node

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/wrap"
)

const (
	nextSize = 2

	headerSize = 0 +
		page.MagicSize +
		nextSize
)

type header struct {
	magic wrap.Int32
	next  wrap.Int16
}

func newHeader(data []byte) (header, []byte) {
	magic, data0 := wrap.NewInt32(data)
	next, data1 := wrap.NewInt16(data0)

	return header{
		magic: magic,
		next:  next,
	}, data1
}

func (h *header) init() {
	h.magic.Set(int32(page.MagicPage))
}

func (h *header) Next() int16 {
	return h.next.Get()
}

func (h *header) setNext(o int16) {
	h.next.Set(o)
}

func (h *header) validate() error {
	return page.Magic(h.magic.Get()).Validate()
}
