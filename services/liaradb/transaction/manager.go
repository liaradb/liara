package transaction

import (
	"context"
	"time"

	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/collection"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/util/set"
)

type Manager struct {
	log           *recovery.Log
	storage       *storage.Storage
	collections   *collection.Collections
	lockTable     *locktable.LockTable[action.ItemID]
	transactionID record.TransactionID
	txReqs        async.Handler[value.TenantID, *Transaction]
	returns       chan record.TransactionID
	active        set.Set[record.TransactionID]
	highWater     record.LogSequenceNumber
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
		txReqs:      make(async.Handler[value.TenantID, *Transaction]),
		returns:     make(chan record.TransactionID), // TODO: Should this be buffered?
		active:      make(set.Set[record.TransactionID]),
	}
}

func (m *Manager) Run(ctx context.Context) {
	go m.run(ctx)
}

func (m *Manager) run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			// TODO: Should we try and drain m.returns?
			m.flush(t)
		case r := <-m.txReqs:
			m.next(r)
		case r := <-m.returns:
			m.end(r)
		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) Active() []record.TransactionID {
	return m.active.Slice()
}

func (m *Manager) isDirty() bool {
	hw := m.log.HighWater()
	isDirty := hw.Value() > m.highWater.Value()
	m.highWater = hw
	return isDirty
}

func (m *Manager) Next(ctx context.Context, tid value.TenantID) (*Transaction, error) {
	return m.txReqs.Send(ctx, tid)
}

func (m *Manager) next(r *async.Request[value.TenantID, *Transaction]) {
	m.transactionID = record.NewTransactionID(m.transactionID.Value() + 1)
	m.active.Add(m.transactionID)
	r.Reply(newTransaction(
		m.transactionID,
		r.Value(),
		m.log,
		NewBufferList(m.storage),
		locktable.NewConcurrencyMgr(m.lockTable),
		m.collections,
		m,
	), nil)
}

func (m *Manager) End(txid record.TransactionID) {
	m.returns <- txid
}

func (m *Manager) end(txid record.TransactionID) {
	m.active.Remove(txid)
}

func (m *Manager) flush(now time.Time) {
	if !m.isDirty() {
		return
	}

	// TODO: What do we do with this error?
	_ = m.Flush(now)
}

func (m *Manager) Flush(now time.Time) error {
	if err := m.storage.FlushAll(); err != nil {
		return err
	}

	_, err := m.log.FlushCheckpoint(now, m.Active())
	return err
}
