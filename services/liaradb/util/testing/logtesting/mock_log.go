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
	err     error
	locks   map[tablename.TableName]map[string]int
}

func (t *MockLog) SetError(err error) {
	t.err = err
}

func (t *MockLog) LogSequenceNumber() logpage.LogSequenceNumber { return t.lsn }
func (t *MockLog) Records() []*record.Record                    { return t.records }

func (t *MockLog) SLock(ctx context.Context, tn tablename.TableName, key string) error {
	if t.locks == nil {
		t.locks = make(map[tablename.TableName]map[string]int)
	}
	locks, ok := t.locks[tn]
	if !ok {
		locks = make(map[string]int)
		t.locks[tn] = locks
	}

	if locks[key] < 2 {
		locks[key] = 1
	}

	return nil
}

func (t *MockLog) XLock(ctx context.Context, tn tablename.TableName, key string) error {
	if t.locks == nil {
		t.locks = make(map[tablename.TableName]map[string]int)
	}
	locks, ok := t.locks[tn]
	if !ok {
		locks = make(map[string]int)
		t.locks[tn] = locks
	}

	if locks[key] < 2 {
		locks[key] = 2
	}

	return nil
}

func (t *MockLog) IsLocked(tn tablename.TableName, key string) bool {
	if t.locks == nil {
		return false
	}

	locks, ok := t.locks[tn]
	if !ok {
		return false
	}

	return locks[key] > 0
}

func (t *MockLog) IsSLocked(tn tablename.TableName, key string) bool {
	if t.locks == nil {
		return false
	}

	locks, ok := t.locks[tn]
	if !ok {
		return false
	}

	return locks[key] == 1
}

func (t *MockLog) IsXLocked(tn tablename.TableName, key string) bool {
	if t.locks == nil {
		return false
	}

	locks, ok := t.locks[tn]
	if !ok {
		return false
	}

	return locks[key] == 2
}

func (t *MockLog) Append(ctx context.Context, frid link.FileRecordID, data []byte) (logpage.LogSequenceNumber, error) {
	if t.err != nil {
		return logpage.LogSequenceNumber{}, t.err
	}

	t.increment()
	t.addRecord(record.New(t.txid, frid, record.ActionAppend, data, nil))
	return t.lsn, nil
}

func (t *MockLog) Insert(ctx context.Context, frid link.FileRecordID, data []byte) (logpage.LogSequenceNumber, error) {
	if t.err != nil {
		return logpage.LogSequenceNumber{}, t.err
	}

	t.increment()
	t.addRecord(record.New(t.txid, frid, record.ActionInsert, data, nil))
	return t.lsn, nil
}

func (t *MockLog) Update(ctx context.Context, frid link.FileRecordID, data []byte, prev []byte) (logpage.LogSequenceNumber, error) {
	if t.err != nil {
		return logpage.LogSequenceNumber{}, t.err
	}

	t.increment()
	t.addRecord(record.New(t.txid, frid, record.ActionUpdate, data, prev))
	return t.lsn, nil
}

func (t *MockLog) Clear(ctx context.Context, frid link.FileRecordID, prev []byte) (logpage.LogSequenceNumber, error) {
	if t.err != nil {
		return logpage.LogSequenceNumber{}, t.err
	}

	t.increment()
	t.addRecord(record.New(t.txid, frid, record.ActionClear, nil, prev))
	return t.lsn, nil
}

func (t *MockLog) UpdateHeader(ctx context.Context, frid link.FileRecordID, data []byte, prev []byte) (logpage.LogSequenceNumber, error) {
	if t.err != nil {
		return logpage.LogSequenceNumber{}, t.err
	}

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
