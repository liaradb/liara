package transaction

import (
	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/manager"
	"github.com/liaradb/liaradb/collection/outbox"
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
	eventLog      *eventlog.EventLog
	keyValue      *keyvalue.KeyValue
	outbox        *outbox.Outbox
	lockTable     *locktable.LockTable[action.ItemID]
	transactionID record.TransactionID
}

func NewManager(
	log *recovery.Log,
	storage *storage.Storage,
	lockTable *locktable.LockTable[action.ItemID],
) *Manager {
	cursor := btree.NewCursor(storage)
	kv := keyvalue.New(storage, cursor)
	return &Manager{
		log:       log,
		storage:   storage,
		manager:   manager.New(kv),
		eventLog:  eventlog.New(storage, cursor),
		keyValue:  kv,
		outbox:    outbox.New(storage, cursor),
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
		m.eventLog,
		m.keyValue,
		m.outbox,
	)
}
