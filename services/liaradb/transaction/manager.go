package transaction

import (
	"context"
	"log/slog"
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

const (
	returnSize = 100
	interval   = 10 * time.Second
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
		returns:     make(chan record.TransactionID, returnSize),
		active:      make(set.Set[record.TransactionID]),
	}
}

func (m *Manager) Run(ctx context.Context) {
	go m.run(ctx)
}

func (m *Manager) run(ctx context.Context) {
	// TODO: This may back up over time
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
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

// TODO: Should we use highWater or lowWater?
func (m *Manager) isDirty() (record.LogSequenceNumber, bool) {
	hw := m.log.HighWater()
	return hw, hw.Value() > m.highWater.Value()
}

func (m *Manager) setHighwater(hw record.LogSequenceNumber) {
	m.highWater = hw
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
	hw, isDirty := m.isDirty()
	if !isDirty {
		return
	}

	m.drainEnd()

	slog.Info("flushing...")

	// TODO: What do we do with this error?
	if err := m.Flush(now); err != nil {
		slog.Error("unable to flush",
			"error", err)
		return
	}

	m.setHighwater(hw)

	slog.Info("flushing complete")
}

func (m *Manager) drainEnd() {
	for range min(returnSize, len(m.returns)) {
		txid := <-m.returns
		m.end(txid)
	}
}

func (m *Manager) Flush(now time.Time) error {
	if err := m.storage.FlushAll(); err != nil {
		return err
	}

	_, err := m.log.FlushCheckpoint(now, m.Active())
	return err
}
