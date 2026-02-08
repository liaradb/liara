package transaction

import (
	"github.com/liaradb/liaradb/collection"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/storage"
)

type Manager struct {
	log           *recovery.Log
	storage       *storage.Storage
	collections   *collection.Collections
	lockTable     *locktable.LockTable[action.ItemID]
	transactionID record.TransactionID
}

func NewManager(
	log *recovery.Log,
	storage *storage.Storage,
	lockTable *locktable.LockTable[action.ItemID],
) *Manager {
	return &Manager{
		log:         log,
		storage:     storage,
		collections: collection.NewCollections(storage),
		lockTable:   lockTable,
	}
}

func (m *Manager) Next(tid value.TenantID) *Transaction {
	m.transactionID = record.NewTransactionID(m.transactionID.Value() + 1)
	return newTransaction(
		m.transactionID,
		tid,
		m.log,
		NewBufferList(m.storage),
		locktable.NewConcurrencyMgr(m.lockTable),
		m.collections,
	)
}
