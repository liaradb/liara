package transaction

import (
	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/log"
	"github.com/liaradb/liaradb/log/action"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/eventlog"
)

type Manager struct {
	log           *log.Log
	storage       *storage.Storage
	lockTable     *locktable.LockTable[action.ItemID]
	transactionID record.TransactionID
}

func NewManager(
	log *log.Log,
	storage *storage.Storage,
	lockTable *locktable.LockTable[action.ItemID],
) *Manager {
	return &Manager{
		log:       log,
		storage:   storage,
		lockTable: lockTable,
	}
}

func (m *Manager) Next() *Transaction {
	m.transactionID++
	return newTransaction(
		m.transactionID,
		m.log,
		NewBufferList(m.storage),
		locktable.NewConcurrencyMgr(m.lockTable),
		eventlog.New(m.storage))
}
