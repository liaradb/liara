package page

import "io"

const RecordHeaderSize = CrcSize + RecordLengthSize

// TODO: Should this be private?
const PageHeaderSize = 0 +
	magicSize +
	pageIDSize +
	timeLineIDSize +
	RecordLengthSize

type PageHeader struct {
	id              PageID
	timeLineID      TimeLineID
	lengthRemaining RecordLength
}

func NewPageHeader(
	id PageID,
	timeLineID TimeLineID,
	lengthRemaining RecordLength,
) PageHeader {
	return PageHeader{
		id:              id,
		timeLineID:      timeLineID,
		lengthRemaining: lengthRemaining,
	}
}

func (ph PageHeader) ID() PageID                    { return ph.id }
func (ph PageHeader) TimeLineID() TimeLineID        { return ph.timeLineID }
func (ph PageHeader) LengthRemaining() RecordLength { return ph.lengthRemaining }

// TODO: Should we store this on the header struct?
func (ph PageHeader) Position(size int64) int64 {
	return int64(ph.id) * (size + PageHeaderSize)
}

func (ph *PageHeader) Read(r io.Reader) error {
	if err := MagicPage.ReadIsPage(r); err != nil {
		return err
	}

	if err := ph.id.Read(r); err != nil {
		return err
	}

	if err := ph.timeLineID.Read(r); err != nil {
		return err
	}

	if err := ph.lengthRemaining.Read(r); err != nil {
		return err
	}

	return nil
}

func (ph *PageHeader) Write(w io.Writer) error {
	if err := MagicPage.Write(w); err != nil {
		return err
	}

	if err := ph.id.Write(w); err != nil {
		return err
	}

	if err := ph.timeLineID.Write(w); err != nil {
		return err
	}

	if err := ph.lengthRemaining.Write(w); err != nil {
		return err
	}

	return nil
}
