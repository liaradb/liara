package record

import (
	"io"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/serializer"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type BufferRecord struct {
	logSequenceNumber logpage.LogSequenceNumber
	tenantID          value.TenantID
	transactionID     TransactionID
	fileRecordID      link.FileRecordID
	action            Action
	data              LogData
	reverse           LogData
}

func NewBufferRecord(
	tid value.TenantID,
	txid TransactionID,
	fileRecordID link.FileRecordID,
	action Action,
	data []byte,
	reverse []byte,
) *BufferRecord {
	return &BufferRecord{
		tenantID:      tid,
		transactionID: txid,
		fileRecordID:  fileRecordID,
		action:        action,
		data:          LogData{data},
		reverse:       LogData{reverse},
	}
}

func (br *BufferRecord) LogSequenceNumber() logpage.LogSequenceNumber {
	return br.logSequenceNumber
}

func (br *BufferRecord) SetLogSequenceNumber(lsn logpage.LogSequenceNumber) {
	br.logSequenceNumber = lsn
}

func (br *BufferRecord) Size() int {
	return serializer.Size(
		br.logSequenceNumber,
		br.tenantID,
		br.transactionID,
		br.fileRecordID,
		br.action,
		&br.data,
		&br.reverse)
}

func (br *BufferRecord) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		br.logSequenceNumber,
		br.tenantID,
		br.transactionID,
		br.fileRecordID,
		br.action,
		&br.data,
		&br.reverse)
}

func (br *BufferRecord) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&br.logSequenceNumber,
		&br.tenantID,
		&br.transactionID,
		&br.fileRecordID,
		&br.action,
		&br.data,
		&br.reverse)
}
