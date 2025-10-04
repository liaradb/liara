package transaction

import (
	"github.com/liaradb/liaradb/log"
	"github.com/liaradb/liaradb/log/record"
)

type Manager struct {
	log           *log.Log
	transactionID record.TransactionID
}

func NewManager(
	log *log.Log,
) *Manager {
	return &Manager{
		log: log,
	}
}

func (m *Manager) Next() *Transaction {
	m.transactionID++
	return newTransaction(m.transactionID, m.log)
}
