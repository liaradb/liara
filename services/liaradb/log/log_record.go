package log

import (
	"io"

	"github.com/liaradb/liaradb/raw"
)

type LogRecord struct {
	logSequenceNumber LogSequenceNumber
	transactionID     TransactionID
	data              LogData
	reverse           LogData
}

func newLogRecord(
	lsn LogSequenceNumber,
	tid TransactionID,
	data []byte,
	reverse []byte,
) *LogRecord {
	return &LogRecord{
		logSequenceNumber: lsn,
		transactionID:     tid,
		data:              LogData{data},
		reverse:           LogData{reverse},
	}
}

func (lr *LogRecord) LogSequenceNumber() LogSequenceNumber { return lr.logSequenceNumber }
func (lr *LogRecord) TransactionID() TransactionID         { return lr.transactionID }
func (lr *LogRecord) Data() []byte                         { return lr.data.Bytes() }
func (lr *LogRecord) Reverse() []byte                      { return lr.reverse.Bytes() }

func (lr *LogRecord) Write(w io.Writer) error {
	if err := lr.logSequenceNumber.Write(w); err != nil {
		return err
	}

	if err := lr.transactionID.Write(w); err != nil {
		return err
	}
	if err := lr.data.Write(w); err != nil {
		return err
	}

	if err := lr.reverse.Write(w); err != nil {
		return err
	}

	return nil
}

func (lr *LogRecord) Read(r io.Reader) error {
	if err := lr.logSequenceNumber.Read(r); err != nil {
		return err
	}

	if err := lr.transactionID.Read(r); err != nil {
		return err
	}

	if err := lr.data.Read(r); err != nil {
		return err
	}

	if err := lr.reverse.Read(r); err != nil {
		return err
	}

	return nil
}

func (lr *LogRecord) Value() []byte {
	data := make([]byte, lr.Length())
	return data
}

func (lr *LogRecord) Length() raw.Offset {
	// [LogRecord.LogSequenceNumber]
	return raw.Uint64Length +
		// [LogRecord.TransactionID]
		raw.Uint64Length +
		// [LogRecord.Length]
		raw.Uint32Length +
		// [LogRecord.Data]
		raw.Offset(lr.data.Length()) +
		// [LogRecord.Reverse]
		raw.Offset(lr.reverse.Length())
}
