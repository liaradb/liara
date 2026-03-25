package record

import (
	"io"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/serializer"
)

type Record struct {
	logSequenceNumber LogSequenceNumber
	tenantID          value.TenantID
	transactionID     TransactionID
	time              Time
	action            Action
	data              LogData
	reverse           LogData
}

func New(
	lsn LogSequenceNumber,
	tid value.TenantID,
	txid TransactionID,
	time time.Time,
	action Action,
	data []byte,
	reverse []byte,
) *Record {
	return &Record{
		logSequenceNumber: lsn,
		tenantID:          tid,
		transactionID:     txid,
		time:              NewTime(time),
		action:            action,
		data:              LogData{data},
		reverse:           LogData{reverse},
	}
}

func (rc *Record) LogSequenceNumber() LogSequenceNumber { return rc.logSequenceNumber }
func (rc *Record) TenantID() value.TenantID             { return rc.tenantID }
func (rc *Record) TransactionID() TransactionID         { return rc.transactionID }
func (rc *Record) Time() time.Time                      { return rc.time.Time }
func (rc *Record) Action() Action                       { return rc.action }
func (rc *Record) Data() []byte                         { return rc.data.Bytes() }
func (rc *Record) Reverse() []byte                      { return rc.reverse.Bytes() }

func (rc *Record) Size() int {
	return serializer.Size(
		rc.logSequenceNumber,
		rc.tenantID,
		rc.transactionID,
		rc.time,
		rc.action,
		&rc.data,
		&rc.reverse)
}

func (rc *Record) Write(w io.Writer) error {
	if err := rc.logSequenceNumber.Write(w); err != nil {
		return err
	}

	if err := rc.tenantID.Write(w); err != nil {
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

	if err := rc.tenantID.Read(r); err != nil {
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

func (rc *Record) Compare(b *Record) bool {
	if rc == b {
		return true
	}

	return rc.logSequenceNumber == b.logSequenceNumber &&
		rc.tenantID == b.tenantID &&
		rc.transactionID == b.transactionID &&
		rc.time.Equal(b.time) &&
		rc.action == b.action &&
		rc.data.Compare(&b.data) &&
		rc.reverse.Compare(&b.reverse)
}
