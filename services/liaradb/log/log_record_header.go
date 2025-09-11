package log

import "io"

type LogRecordHeader struct {
	logSequenceNumber LogSequenceNumber
	transactionID     TransactionID
	dataLength        LogRecordLength
	reverseLength     LogRecordLength
}

func (lrh *LogRecordHeader) LogSequenceNumber() LogSequenceNumber { return lrh.logSequenceNumber }
func (lrh *LogRecordHeader) TransactionID() TransactionID         { return lrh.transactionID }
func (lrh *LogRecordHeader) DataLength() LogRecordLength          { return lrh.dataLength }
func (lrh *LogRecordHeader) ReverseLength() LogRecordLength       { return lrh.reverseLength }

func (lrh *LogRecordHeader) Write(w io.Writer) error {
	if err := lrh.logSequenceNumber.Write(w); err != nil {
		return err
	}

	if err := lrh.transactionID.Write(w); err != nil {
		return err
	}

	if err := lrh.dataLength.Write(w); err != nil {
		return err
	}

	return lrh.reverseLength.Write(w)
}

func (lrh *LogRecordHeader) Read(r io.Reader) error {
	if err := lrh.logSequenceNumber.Read(r); err != nil {
		return err
	}

	if err := lrh.transactionID.Read(r); err != nil {
		return err
	}

	if err := lrh.dataLength.Read(r); err != nil {
		return err
	}

	return lrh.reverseLength.Read(r)
}
