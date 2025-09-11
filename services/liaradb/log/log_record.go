package log

import (
	"bytes"
	"io"

	"github.com/liaradb/liaradb/raw"
)

type LogRecord struct {
	header  LogRecordHeader
	data    []byte
	reverse []byte
	crc     CRC
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
func (lr *LogRecord) CRC() CRC                             { return lr.crc }
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

func (lr *LogRecord) WriteFull(b *bytes.Buffer) error {
	if err := lr.Write(b); err != nil {
		return err
	}

	lr.crc = NewCRC(b.Bytes())
	return lr.crc.Write(b)
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
		// [LogRecord.CRC]
		raw.Uint32Length +
		// [LogRecord.Data]
		raw.Offset(len(lr.data)) +
		// [LogRecord.Reverse]
		raw.Offset(len(lr.reverse))
}
