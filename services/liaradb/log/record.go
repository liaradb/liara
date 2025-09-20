package log

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/liaradb/liaradb/log/record"
)

type Record struct {
	logSequenceNumber record.LogSequenceNumber
	transactionID     TransactionID
	time              time.Time
	data              LogData
	reverse           LogData
}

func newRecord(
	lsn record.LogSequenceNumber,
	tid TransactionID,
	time time.Time,
	data []byte,
	reverse []byte,
) *Record {
	return &Record{
		logSequenceNumber: lsn,
		transactionID:     tid,
		time:              time,
		data:              LogData{data},
		reverse:           LogData{reverse},
	}
}

func (rc *Record) LogSequenceNumber() record.LogSequenceNumber { return rc.logSequenceNumber }
func (rc *Record) TransactionID() TransactionID                { return rc.transactionID }
func (rc *Record) Time() time.Time                             { return rc.time }
func (rc *Record) Data() []byte                                { return rc.data.Bytes() }
func (rc *Record) Reverse() []byte                             { return rc.reverse.Bytes() }

func (rc *Record) Write(w io.Writer) error {
	if err := rc.logSequenceNumber.Write(w); err != nil {
		return err
	}

	if err := rc.transactionID.Write(w); err != nil {
		return err
	}

	if err := rc.writeTime(w); err != nil {
		return err
	}

	if err := rc.data.Write(w); err != nil {
		return err
	}

	if err := rc.reverse.Write(w); err != nil {
		return err
	}

	return nil
}

func (rc *Record) Read(r io.Reader) error {
	if err := rc.logSequenceNumber.Read(r); err != nil {
		return err
	}

	if err := rc.transactionID.Read(r); err != nil {
		return err
	}

	if err := rc.readTime(r); err != nil {
		return err
	}

	if err := rc.data.Read(r); err != nil {
		return err
	}

	if err := rc.reverse.Read(r); err != nil {
		return err
	}

	return nil
}

func (rc *Record) writeTime(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, rc.time.UnixMicro())
}

func (rc *Record) readTime(r io.Reader) error {
	var t int64
	if err := binary.Read(r, binary.BigEndian, &t); err != nil {
		return err
	}

	rc.time = time.UnixMicro(t)
	return nil
}
