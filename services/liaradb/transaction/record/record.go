package record

import (
	"io"

	"github.com/liaradb/liaradb/encoder/serializer"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type Record struct {
	logSequenceNumber logpage.LogSequenceNumber
	transactionID     TransactionID
	recordLocator     link.RecordLocator
	action            Action
	collection        Collection
	data              LogData
	reverse           LogData
}

func New(
	txid TransactionID,
	recordLocator link.RecordLocator,
	action Action,
	collection Collection,
	data []byte,
	reverse []byte,
) *Record {
	return &Record{
		transactionID: txid,
		recordLocator: recordLocator,
		action:        action,
		collection:    collection,
		data:          LogData{data},
		reverse:       LogData{reverse},
	}
}

func (rc *Record) LogSequenceNumber() logpage.LogSequenceNumber { return rc.logSequenceNumber }
func (rc *Record) TransactionID() TransactionID                 { return rc.transactionID }
func (rc *Record) RecordLocator() link.RecordLocator            { return rc.recordLocator }
func (rc *Record) Action() Action                               { return rc.action }
func (rc *Record) Collection() Collection                       { return rc.collection }
func (rc *Record) Data() []byte                                 { return rc.data.Bytes() }
func (rc *Record) Reverse() []byte                              { return rc.reverse.Bytes() }
func (rc *Record) IsCheckpoint() bool                           { return rc.action == ActionCheckpoint }

func (rc *Record) SetLogSequenceNumber(lsn logpage.LogSequenceNumber) {
	rc.logSequenceNumber = lsn
}

func (rc *Record) Size() int {
	return serializer.Size(
		rc.logSequenceNumber,
		rc.transactionID,
		rc.recordLocator,
		rc.action,
		rc.collection,
		&rc.data,
		&rc.reverse)
}

func (rc *Record) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		rc.logSequenceNumber,
		rc.transactionID,
		rc.recordLocator,
		rc.action,
		rc.collection,
		&rc.data,
		&rc.reverse)
}

func (rc *Record) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&rc.logSequenceNumber,
		&rc.transactionID,
		&rc.recordLocator,
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
		rc.transactionID == b.transactionID &&
		rc.recordLocator == b.recordLocator &&
		rc.action == b.action &&
		rc.collection == b.collection &&
		rc.data.Compare(&b.data) &&
		rc.reverse.Compare(&b.reverse)
}
