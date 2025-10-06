package record

import (
	"io"
	"time"
)

type Record struct {
	logSequenceNumber LogSequenceNumber
	transactionID     TransactionID
	time              Time
	action            Action
	data              LogData
	reverse           LogData
}

func New(
	lsn LogSequenceNumber,
	tid TransactionID,
	time time.Time,
	action Action,
	data []byte,
	reverse []byte,
) *Record {
	return &Record{
		logSequenceNumber: lsn,
		transactionID:     tid,
		time:              NewTime(time),
		action:            action,
		data:              LogData{data},
		reverse:           LogData{reverse},
	}
}

func (rc *Record) LogSequenceNumber() LogSequenceNumber { return rc.logSequenceNumber }
func (rc *Record) TransactionID() TransactionID         { return rc.transactionID }
func (rc *Record) Time() time.Time                      { return rc.time.Time }
func (rc *Record) Action() Action                       { return rc.action }
func (rc *Record) Data() []byte                         { return rc.data.Bytes() }
func (rc *Record) Reverse() []byte                      { return rc.reverse.Bytes() }

func (rc *Record) Size() int {
	return size(
		rc.logSequenceNumber,
		rc.transactionID,
		rc.time,
		rc.action,
		rc.data,
		rc.reverse)
}

func (rc *Record) Write(w io.Writer) error {
	if err := rc.logSequenceNumber.Write(w); err != nil {
		return err
	}

	if err := rc.transactionID.Write(w); err != nil {
		return err
	}

	if err := rc.time.Write(w); err != nil {
		return err
	}

	if err := rc.action.Write(w); err != nil {
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

	if err := rc.time.Read(r); err != nil {
		return err
	}

	if err := rc.action.Read(r); err != nil {
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
