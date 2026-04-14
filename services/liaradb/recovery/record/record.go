package record

import (
	"io"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/serializer"
)

type Record struct {
	logSequenceNumber LogSequenceNumber
	tenantID          value.TenantID
	transactionID     TransactionID
	time              Time
	action            Action
	collection        Collection
	data              LogData
	reverse           LogData
}

func New(
	lsn LogSequenceNumber,
	tid value.TenantID,
	txid TransactionID,
	time Time,
	action Action,
	collection Collection,
	data []byte,
	reverse []byte,
) *Record {
	return &Record{
		logSequenceNumber: lsn,
		tenantID:          tid,
		transactionID:     txid,
		time:              time,
		action:            action,
		collection:        collection,
		data:              LogData{data},
		reverse:           LogData{reverse},
	}
}

func (rc *Record) LogSequenceNumber() LogSequenceNumber { return rc.logSequenceNumber }
func (rc *Record) TenantID() value.TenantID             { return rc.tenantID }
func (rc *Record) TransactionID() TransactionID         { return rc.transactionID }
func (rc *Record) Time() Time                           { return rc.time }
func (rc *Record) Action() Action                       { return rc.action }
func (rc *Record) Collection() Collection               { return rc.collection }
func (rc *Record) Data() []byte                         { return rc.data.Bytes() }
func (rc *Record) Reverse() []byte                      { return rc.reverse.Bytes() }
func (rc *Record) IsCheckpoint() bool                   { return rc.action == ActionCheckpoint }

func (rc *Record) Size() int {
	return serializer.Size(
		rc.logSequenceNumber,
		rc.tenantID,
		rc.transactionID,
		rc.time,
		rc.action,
		rc.collection,
		&rc.data,
		&rc.reverse)
}

func (rc *Record) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		rc.logSequenceNumber,
		rc.tenantID,
		rc.transactionID,
		rc.time,
		rc.action,
		rc.collection,
		&rc.data,
		&rc.reverse)
}

func (rc *Record) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&rc.logSequenceNumber,
		&rc.tenantID,
		&rc.transactionID,
		&rc.time,
		&rc.action,
		&rc.collection,
		&rc.data,
		&rc.reverse)
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
		rc.collection == b.collection &&
		rc.data.Compare(&b.data) &&
		rc.reverse.Compare(&b.reverse)
}
