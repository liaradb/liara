package transaction

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/transaction/locktable"
	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/transaction/record"
	"github.com/liaradb/liaradb/transaction/replay"
	"github.com/liaradb/liaradb/util/async"
	"github.com/liaradb/liaradb/util/set"
)

const (
	returnSize = 100
	interval   = 10 * time.Second
)

type Manager struct {
	log           *log.Log
	storage       *storage.Storage
	lockTable     *locktable.LockTable[ItemID]
	transactionID record.TransactionID
	txReqs        async.Handler[struct{}, *Transaction]
	returns       chan record.TransactionID
	active        set.Set[record.TransactionID]
	checkpoint    logpage.LogSequenceNumber
}

func NewManager(
	log *log.Log,
	storage *storage.Storage,
	ltInSize int,
) *Manager {
	return &Manager{
		log:       log,
		storage:   storage,
		lockTable: locktable.New[ItemID](ltInSize),
		txReqs:    make(async.Handler[struct{}, *Transaction]),
		returns:   make(chan record.TransactionID, returnSize),
		active:    make(set.Set[record.TransactionID]),
	}
}

func (m *Manager) Run(ctx context.Context, onError func(error)) {
	slog.Info("starting...")

	go m.lockTable.Run(ctx)
	go func() { onError(m.run(ctx)) }()

	slog.Info("transaction manager running")
}

func (m *Manager) Recover(ctx context.Context) error {
	if err := m.storage.Run(ctx); err != nil {
		return err
	}

	slog.Info("storage running")

	if err := m.log.Run(ctx); err != nil {
		return err
	}

	slog.Info("recovering...")

	r := replay.NewReplay(m.storage, m.log)
	if err := r.Recover(ctx); err != nil {
		return err
	}

	slog.Info("recovered")

	if err := m.log.StartWriter(); err != nil {
		return err
	}

	slog.Info("log running")

	return nil
}

func (m *Manager) run(ctx context.Context) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := m.appendCheckpoint(ctx); err != nil {
				return err
			}
		case r := <-m.txReqs:
			m.next(r)
		case r := <-m.returns:
			m.end(r)
		case <-ctx.Done():
			return nil
		}
	}
}

func (m *Manager) Active() []record.TransactionID {
	return m.active.Slice()
}

func (m *Manager) isDirty() bool {
	return m.log.HighWater().Value() > m.checkpoint.Value()
}

func (m *Manager) setCheckpoint(cp logpage.LogSequenceNumber) {
	m.checkpoint = cp
}

func (m *Manager) Next(ctx context.Context) (*Transaction, error) {
	return m.txReqs.Send(ctx, struct{}{})
}

func (m *Manager) next(r *async.Request[struct{}, *Transaction]) {
	m.transactionID = record.NewTransactionID(m.transactionID.Value() + 1)
	m.active.Add(m.transactionID)
	r.Reply(newTransaction(
		m.transactionID,
		m.log,
		NewBufferList(m.storage),
		locktable.NewConcurrencyMgr(m.lockTable),
		m,
	), nil)
}

func (m *Manager) End(txid record.TransactionID) {
	m.returns <- txid
}

func (m *Manager) end(txid record.TransactionID) {
	m.active.Remove(txid)
}

func (m *Manager) appendCheckpoint(ctx context.Context) error {
	m.drainEnd()
	if !m.isDirty() {
		return nil
	}

	slog.Info("flushing...")

	if err := m.flush(ctx); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("unable to flush",
			"error", err)
		return err
	}

	slog.Info("flushing complete")

	return nil
}

func (m *Manager) drainEnd() {
	for range min(returnSize, len(m.returns)) {
		txid := <-m.returns
		m.end(txid)
	}
}

func (m *Manager) drainAll() {
	for {
		select {
		case txid := <-m.returns:
			m.end(txid)
		default:
			return
		}
	}
}

func (m *Manager) flush(ctx context.Context) error {
	if err := m.storage.FlushUnpinned(ctx); err != nil {
		return err
	}

	lsn, err := m.log.Checkpoint(ctx, m.Active()...)
	if err != nil {
		return err
	}

	m.setCheckpoint(lsn)
	return nil
}

func (m *Manager) Shutdown(ctx context.Context, now time.Time) error {
	m.lockTable.Close()

	m.drainAll()

	if err := m.storage.FlushUnpinned(ctx); err != nil {
		return err
	}

	if !m.isDirty() {
		return nil
	}

	lsn, err := m.log.Checkpoint(ctx, m.Active()...)
	if err != nil {
		return err
	}

	m.setCheckpoint(lsn)
	return nil
}
