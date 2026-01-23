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
	next wrap.Int16
}

func newHeader(data []byte) (header, []byte) {
	// TODO: Should this have a magic entry?
	_, data0 := wrap.NewInt32(data) // Magic
	next, data1 := wrap.NewInt16(data0)

	return header{
		next: next,
	}, data1
}

func (h *header) Next() int16 {
	return h.next.Get()
}

func (h *header) setNext(o int16) {
	h.next.Set(o)
}
