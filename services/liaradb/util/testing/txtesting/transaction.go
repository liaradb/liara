package txtesting

import (
	"context"

	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/record"
)

type Transaction struct {
	txid    record.TransactionID
	lsn     logpage.LogSequenceNumber
	records []*record.Record
	err     error
	locks   map[tablename.TableName]map[string]int
}

func (t *Transaction) SetError(err error) {
	t.err = err
}

func (t *Transaction) LogSequenceNumber() logpage.LogSequenceNumber { return t.lsn }
func (t *Transaction) Records() []*record.Record                    { return t.records }

func (t *Transaction) SLock(ctx context.Context, tn tablename.TableName, key string) error {
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

func (t *Transaction) XLock(ctx context.Context, tn tablename.TableName, key string) error {
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

func (t *Transaction) IsLocked(tn tablename.TableName, key string) bool {
	if t.locks == nil {
		return false
	}

	locks, ok := t.locks[tn]
	if !ok {
		return false
	}

	return locks[key] > 0
}

func (t *Transaction) IsSLocked(tn tablename.TableName, key string) bool {
	if t.locks == nil {
		return false
	}

	locks, ok := t.locks[tn]
	if !ok {
		return false
	}

	return locks[key] == 1
}

func (t *Transaction) IsXLocked(tn tablename.TableName, key string) bool {
	if t.locks == nil {
		return false
	}

	locks, ok := t.locks[tn]
	if !ok {
		return false
	}

	return locks[key] == 2
}

func (t *Transaction) Append(ctx context.Context, frid link.FileRecordID, data []byte) (logpage.LogSequenceNumber, error) {
	if t.err != nil {
		return logpage.LogSequenceNumber{}, t.err
	}

	t.increment()
	t.addRecord(record.New(t.txid, frid, record.ActionAppend, data, nil))
	return t.lsn, nil
}

func (t *Transaction) Insert(ctx context.Context, frid link.FileRecordID, data []byte) (logpage.LogSequenceNumber, error) {
	if t.err != nil {
		return logpage.LogSequenceNumber{}, t.err
	}

	t.increment()
	t.addRecord(record.New(t.txid, frid, record.ActionInsert, data, nil))
	return t.lsn, nil
}

func (t *Transaction) Update(ctx context.Context, frid link.FileRecordID, data []byte, prev []byte) (logpage.LogSequenceNumber, error) {
	if t.err != nil {
		return logpage.LogSequenceNumber{}, t.err
	}

	t.increment()
	t.addRecord(record.New(t.txid, frid, record.ActionUpdate, data, prev))
	return t.lsn, nil
}

func (t *Transaction) Clear(ctx context.Context, frid link.FileRecordID, prev []byte) (logpage.LogSequenceNumber, error) {
	if t.err != nil {
		return logpage.LogSequenceNumber{}, t.err
	}

	t.increment()
	t.addRecord(record.New(t.txid, frid, record.ActionClear, nil, prev))
	return t.lsn, nil
}

func (t *Transaction) UpdateHeader(ctx context.Context, frid link.FileRecordID, data []byte, prev []byte) (logpage.LogSequenceNumber, error) {
	if t.err != nil {
		return logpage.LogSequenceNumber{}, t.err
	}

	t.increment()
	t.addRecord(record.New(t.txid, frid, record.ActionUpdateHeader, data, prev))
	return t.lsn, nil
}

func (t *Transaction) increment() {
	t.lsn = t.lsn.Increment()
}

func (t *Transaction) addRecord(r *record.Record) {
	t.records = append(t.records, r)
}

func (t *Transaction) Pin(ctx context.Context, fbid link.FileBlockID) (*storage.Buffer, error) {
	panic("unimplemented")
}

func (t *Transaction) PinCurrent(ctx context.Context, fn link.FileName) (*storage.Buffer, error) {
	panic("unimplemented")
}

func (t *Transaction) PinNext(ctx context.Context, fn link.FileName) (*storage.Buffer, error) {
	panic("unimplemented")
}
