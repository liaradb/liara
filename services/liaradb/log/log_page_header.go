package log

import "io"

type LogPageHeader struct {
	id              LogPageID
	timeLineID      TimeLineID
	lengthRemaining LogRecordLength
}

func newLogPageHeader(
	id LogPageID,
	timeLineID TimeLineID,
	lengthRemaining LogRecordLength,
) *LogPageHeader {
	return &LogPageHeader{
		id:              id,
		timeLineID:      timeLineID,
		lengthRemaining: lengthRemaining,
	}
}

func (lph *LogPageHeader) ID() LogPageID                    { return lph.id }
func (lph *LogPageHeader) TimeLineID() TimeLineID           { return lph.timeLineID }
func (lph *LogPageHeader) LengthRemaining() LogRecordLength { return lph.lengthRemaining }

func (lph *LogPageHeader) Read(r io.Reader) error {
	if err := LogMagicPage.ReadIsPage(r); err != nil {
		return err
	}

	if err := lph.id.Read(r); err != nil {
		return err
	}

	if err := lph.timeLineID.Read(r); err != nil {
		return err
	}

	// if err := lph.lengthRemaining.Read(r); err != nil {
	// 	return err
	// }

	return nil
}

func (lph *LogPageHeader) Write(w io.Writer) error {
	if err := LogMagicPage.Write(w); err != nil {
		return err
	}

	if err := lph.id.Write(w); err != nil {
		return err
	}

	if err := lph.timeLineID.Write(w); err != nil {
		return err
	}

	// if err := lph.lengthRemaining.Write(w); err != nil {
	// 	return err
	// }

	return nil
}
