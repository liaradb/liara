package page

import (
	"io"

	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/raw"
)

const headerSize = 0 +
	magicSize +
	pageIDSize +
	timeLineIDSize +
	record.LengthSize

type Header struct {
	id              PageID
	timeLineID      TimeLineID
	lengthRemaining record.Length
}

func NewHeader(
	id PageID,
	timeLineID TimeLineID,
	lengthRemaining record.Length,
) Header {
	return Header{
		id:              id,
		timeLineID:      timeLineID,
		lengthRemaining: lengthRemaining,
	}
}

func (h Header) ID() PageID                     { return h.id }
func (h Header) TimeLineID() TimeLineID         { return h.timeLineID }
func (h Header) LengthRemaining() record.Length { return h.lengthRemaining }

func (h Header) Size() int {
	return headerSize
}

func (h *Header) Read(r io.Reader) error {
	var m Magic
	return raw.ReadAll(r,
		&m,
		&h.id,
		&h.timeLineID,
		&h.lengthRemaining)
}

func (h *Header) Write(w io.Writer) error {
	return raw.WriteAll(w,
		MagicPage,
		h.id,
		h.timeLineID,
		h.lengthRemaining)
}
