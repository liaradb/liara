package node

import (
	"github.com/liaradb/liaradb/encoder/wrap"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

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

// TODO: Fix these
func (p *header) ID() action.PageID              { return 0 }
func (p *header) TimeLineID() action.TimeLineID  { return 0 }
func (p *header) LengthRemaining() record.Length { return record.NewLength(0) }
