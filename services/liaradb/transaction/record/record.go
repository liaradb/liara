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
	fileRecordID      link.FileRecordID
	action            Action
	data              LogData
	reverse           LogData
}

func New(
	txid TransactionID,
	fileRecordID link.FileRecordID,
	action Action,
	data []byte,
	reverse []byte,
) *Record {
	return &Record{
		transactionID: txid,
		fileRecordID:  fileRecordID,
		action:        action,
		data:          LogData{data},
		reverse:       LogData{reverse},
	}
}

func (rc *Record) LogSequenceNumber() logpage.LogSequenceNumber { return rc.logSequenceNumber }
func (rc *Record) TransactionID() TransactionID                 { return rc.transactionID }
func (rc *Record) FileRecordID() link.FileRecordID              { return rc.fileRecordID }
func (rc *Record) Action() Action                               { return rc.action }
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
		rc.fileRecordID,
		rc.action,
		&rc.data,
		&rc.reverse)
}

func (rc *Record) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		rc.logSequenceNumber,
		rc.transactionID,
		rc.fileRecordID,
		rc.action,
		&rc.data,
		&rc.reverse)
}

func (rc *Record) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&rc.logSequenceNumber,
		&rc.transactionID,
		&rc.fileRecordID,
		&rc.action,
		&rc.data,
		&rc.reverse)
}

func (rc *Record) Compare(b *Record) bool {
	if rc == b {
		return true
	}

	return rc.logSequenceNumber == b.logSequenceNumber &&
		rc.transactionID == b.transactionID &&
		rc.fileRecordID == b.fileRecordID &&
		rc.action == b.action &&
		rc.data.Compare(&b.data) &&
		rc.reverse.Compare(&b.reverse)
}
