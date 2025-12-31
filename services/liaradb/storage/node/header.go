package node

import "github.com/liaradb/liaradb/encoder/wrap"

// TODO: Should this have a magic entry?
const (
	nextSize = 2

	headerSize = nextSize
)

type header struct {
	next wrap.Int16
}

func newHeader(data []byte) (header, []byte) {
	next, data0 := wrap.NewInt16(data)

	return header{
		next: next,
	}, data0
}

func (h *header) Next() int16 {
	return h.next.Get()
}

func (h *header) setNext(o int16) {
	h.next.Set(o)
}
