package log

import (
	"io"

	"github.com/liaradb/liaradb/raw"
)

type LogRecord struct {
	header  LogRecordHeader
	data    []byte
	reverse []byte
}

func newLogRecord(
	lsn LogSequenceNumber,
	tid TransactionID,
	data []byte,
	reverse []byte,
) *LogRecord {
	return &LogRecord{
		header: LogRecordHeader{
			logSequenceNumber: lsn,
			transactionID:     tid,
		},
		data:    data,
		reverse: reverse,
	}
}

func (lr *LogRecord) LogSequenceNumber() LogSequenceNumber { return lr.header.LogSequenceNumber() }
func (lr *LogRecord) TransactionID() TransactionID         { return lr.header.TransactionID() }
func (lr *LogRecord) Data() []byte                         { return lr.data }
func (lr *LogRecord) Reverse() []byte                      { return lr.reverse }

func (lr *LogRecord) Write(w io.Writer) error {
	if err := lr.header.Write(w); err != nil {
		return err
	}

	if _, err := w.Write(lr.data); err != nil {
		return err
	}

	if _, err := w.Write(lr.reverse); err != nil {
		return err
	}

	return nil
}

func (lr *LogRecord) Read(r io.Reader) error {
	if err := lr.header.Read(r); err != nil {
		return err
	}

	if _, err := r.Read(lr.data); err != nil {
		return err
	}

	if _, err := r.Read(lr.reverse); err != nil {
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
		raw.Offset(len(lr.data)) +
		// [LogRecord.Reverse]
		raw.Offset(len(lr.reverse))
}
