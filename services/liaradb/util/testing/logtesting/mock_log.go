package logtesting

import (
	"context"

	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/record"
)

type MockLog struct {
	txid    record.TransactionID
	lsn     logpage.LogSequenceNumber
	records []*record.Record
}

func (t *MockLog) LogSequenceNumber() logpage.LogSequenceNumber { return t.lsn }
func (t *MockLog) Records() []*record.Record                    { return t.records }

func (t *MockLog) SLock(context.Context, tablename.TableName, string) error {
	return nil
}

func (t *MockLog) XLock(context.Context, tablename.TableName, string) error {
	return nil
}

func (t *MockLog) Append(ctx context.Context, frid link.FileRecordID, data []byte) (logpage.LogSequenceNumber, error) {
	t.increment()
	t.addRecord(record.New(t.txid, frid, record.ActionAppend, data, nil))
	return t.lsn, nil
}

func (t *MockLog) Insert(ctx context.Context, frid link.FileRecordID, data []byte) (logpage.LogSequenceNumber, error) {
	t.increment()
	t.addRecord(record.New(t.txid, frid, record.ActionInsert, data, nil))
	return t.lsn, nil
}

func (t *MockLog) Update(ctx context.Context, frid link.FileRecordID, data []byte, prev []byte) (logpage.LogSequenceNumber, error) {
	t.increment()
	t.addRecord(record.New(t.txid, frid, record.ActionUpdate, data, prev))
	return t.lsn, nil
}

func (t *MockLog) Clear(ctx context.Context, frid link.FileRecordID, prev []byte) (logpage.LogSequenceNumber, error) {
	t.increment()
	t.addRecord(record.New(t.txid, frid, record.ActionClear, nil, prev))
	return t.lsn, nil
}

func (t *MockLog) UpdateHeader(ctx context.Context, frid link.FileRecordID, data []byte, prev []byte) (logpage.LogSequenceNumber, error) {
	t.increment()
	t.addRecord(record.New(t.txid, frid, record.ActionUpdateHeader, data, prev))
	return t.lsn, nil
}

func (t *MockLog) increment() {
	t.lsn = t.lsn.Increment()
}

func (t *MockLog) addRecord(r *record.Record) {
	t.records = append(t.records, r)
}
