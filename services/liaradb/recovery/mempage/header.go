package mempage

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

const headerSize = 0 +
	action.PageIDSize +
	action.TimeLineIDSize +
	record.LengthSize

type header struct {
	id              action.PageID
	timeLineID      action.TimeLineID
	lengthRemaining record.Length
}

func newHeader(
	id action.PageID,
	timeLineID action.TimeLineID,
	lengthRemaining record.Length,
) *header {
	return &header{
		id:              id,
		timeLineID:      timeLineID,
		lengthRemaining: lengthRemaining,
	}
}

func (h *header) ID() action.PageID              { return h.id }
func (h *header) TimeLineID() action.TimeLineID  { return h.timeLineID }
func (h *header) LengthRemaining() record.Length { return h.lengthRemaining }

func (h header) Size() int {
	return headerSize
}

func (h *header) Read(r io.Reader) error {
	return raw.ReadAll(r,
		&h.id,
		&h.timeLineID,
		&h.lengthRemaining)
}

func (h *header) Write(w io.Writer) error {
	return raw.WriteAll(w,
		h.id,
		h.timeLineID,
		h.lengthRemaining)
}
