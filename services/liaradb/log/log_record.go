package log

import (
	"encoding/binary"
	"io"
	"time"
)

type LogRecord struct {
	logSequenceNumber LogSequenceNumber
	transactionID     TransactionID
	time              time.Time
	data              LogData
	reverse           LogData
}

func newLogRecord(
	lsn LogSequenceNumber,
	tid TransactionID,
	time time.Time,
	data []byte,
	reverse []byte,
) *LogRecord {
	return &LogRecord{
		logSequenceNumber: lsn,
		transactionID:     tid,
		time:              time,
		data:              LogData{data},
		reverse:           LogData{reverse},
	}
}

func (lr *LogRecord) LogSequenceNumber() LogSequenceNumber { return lr.logSequenceNumber }
func (lr *LogRecord) TransactionID() TransactionID         { return lr.transactionID }
func (lr *LogRecord) Time() time.Time                      { return lr.time }
func (lr *LogRecord) Data() []byte                         { return lr.data.Bytes() }
func (lr *LogRecord) Reverse() []byte                      { return lr.reverse.Bytes() }

func (lr *LogRecord) Write(w io.Writer) error {
	if err := lr.logSequenceNumber.Write(w); err != nil {
		return err
	}

	if err := lr.transactionID.Write(w); err != nil {
		return err
	}

	if err := lr.writeTime(w); err != nil {
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

	if err := lr.readTime(r); err != nil {
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

func (lr *LogRecord) writeTime(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, lr.time.UnixMicro())
}

func (lr *LogRecord) readTime(r io.Reader) error {
	var t int64
	if err := binary.Read(r, binary.BigEndian, &t); err != nil {
		return err
	}

	lr.time = time.UnixMicro(t)
	return nil
}
