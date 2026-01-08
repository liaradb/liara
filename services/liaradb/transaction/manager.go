package transaction

import (
	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/manager"
	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/storage"
)

type Manager struct {
	log           *recovery.Log
	storage       *storage.Storage
	manager       *manager.Manager
	cursor        *btree.Cursor
	eventLog      *eventlog.EventLog
	keyValue      *keyvalue.KeyValue
	lockTable     *locktable.LockTable[action.ItemID]
	transactionID record.TransactionID
}

func NewManager(
	log *recovery.Log,
	storage *storage.Storage,
	lockTable *locktable.LockTable[action.ItemID],
) *Manager {
	cursor := btree.NewCursor(storage)
	return &Manager{
		log:       log,
		storage:   storage,
		manager:   manager.New(storage),
		cursor:    cursor,
		eventLog:  eventlog.New(storage, cursor),
		keyValue:  keyvalue.New(storage),
		lockTable: lockTable,
	}
}

func (m *Manager) Next() *Transaction {
	m.transactionID = record.NewTransactionID(m.transactionID.Value() + 1)
	return newTransaction(
		m.transactionID,
		m.log,
		NewBufferList(m.storage),
		locktable.NewConcurrencyMgr(m.lockTable),
		m.manager,
		m.cursor,
		m.eventLog,
		m.keyValue,
	)
}
