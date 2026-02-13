package transaction

import (
	"context"

	"github.com/liaradb/liaradb/async"
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
	reqs          async.Handler[value.TenantID, *Transaction]
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
		reqs:        make(async.Handler[value.TenantID, *Transaction]),
	}
}

func (m *Manager) Run(ctx context.Context) {
	go m.run(ctx)
}

func (m *Manager) run(ctx context.Context) {
	for {
		select {
		case r := <-m.reqs:
			m.next(r)
		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) Next(ctx context.Context, tid value.TenantID) (*Transaction, error) {
	return m.reqs.Send(ctx, tid)
}

func (m *Manager) next(r *async.Request[value.TenantID, *Transaction]) {
	m.transactionID = record.NewTransactionID(m.transactionID.Value() + 1)
	r.Reply(newTransaction(
		m.transactionID,
		r.Value(),
		m.log,
		NewBufferList(m.storage),
		locktable.NewConcurrencyMgr(m.lockTable),
		m.collections,
	), nil)
}
