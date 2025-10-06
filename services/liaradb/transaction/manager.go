package transaction

import (
	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/log"
	"github.com/liaradb/liaradb/log/action"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/storage"
)

type Manager struct {
	log            *log.Log
	storage        *storage.Storage
	concurrencyMgr *locktable.ConcurrencyMgr[action.ItemID]
	transactionID  record.TransactionID
}

func NewManager(
	log *log.Log,
	storage *storage.Storage,
	concurrencyMgr *locktable.ConcurrencyMgr[action.ItemID],
) *Manager {
	return &Manager{
		log:            log,
		storage:        storage,
		concurrencyMgr: concurrencyMgr,
	}
}

func (m *Manager) Next() *Transaction {
	m.transactionID++
	return newTransaction(
		m.transactionID,
		m.log,
		storage.NewBufferList(m.storage),
		m.concurrencyMgr)
}
