package page

import "io"

const headerSize = 0 +
	magicSize +
	pageIDSize +
	timeLineIDSize +
	RecordLengthSize

type Header struct {
	id              PageID
	timeLineID      TimeLineID
	lengthRemaining RecordLength
}

func NewHeader(
	id PageID,
	timeLineID TimeLineID,
	lengthRemaining RecordLength,
) Header {
	return Header{
		id:              id,
		timeLineID:      timeLineID,
		lengthRemaining: lengthRemaining,
	}
}

func (h Header) ID() PageID                    { return h.id }
func (h Header) TimeLineID() TimeLineID        { return h.timeLineID }
func (h Header) LengthRemaining() RecordLength { return h.lengthRemaining }

func (h Header) Size() int {
	return headerSize
}

func (h *Header) Read(r io.Reader) error {
	if err := MagicPage.ReadIsPage(r); err != nil {
		return err
	}

	if err := h.id.Read(r); err != nil {
		return err
	}

	if err := h.timeLineID.Read(r); err != nil {
		return err
	}

	if err := h.lengthRemaining.Read(r); err != nil {
		return err
	}

	return nil
}

func (h *Header) Write(w io.Writer) error {
	if err := MagicPage.Write(w); err != nil {
		return err
	}

	if err := h.id.Write(w); err != nil {
		return err
	}

	if err := h.timeLineID.Write(w); err != nil {
		return err
	}

	if err := h.lengthRemaining.Write(w); err != nil {
		return err
	}

	return nil
}
