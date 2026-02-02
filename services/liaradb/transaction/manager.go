package transaction

import (
	"github.com/liaradb/liaradb/collection"
	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/idempotency"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/manager"
	"github.com/liaradb/liaradb/collection/outbox"
	"github.com/liaradb/liaradb/collection/tenant"
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
	manager       *manager.Manager
	tenant        *tenant.Tenant
	eventLog      *eventlog.EventLog
	keyValue      *keyvalue.KeyValue
	outbox        *outbox.Outbox
	idempotency   *idempotency.Idempotency
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
		log:         log,
		storage:     storage,
		collections: collection.NewCollections(storage),
		manager:     manager.New(kv),
		tenant:      tenant.New(storage, cursor),
		eventLog:    eventlog.New(storage, cursor),
		keyValue:    kv,
		outbox:      outbox.New(storage, cursor),
		idempotency: idempotency.New(storage, cursor),
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
		m.manager,
		m.tenant,
		m.eventLog,
		m.keyValue,
		m.outbox,
		m.idempotency,
	)
}
