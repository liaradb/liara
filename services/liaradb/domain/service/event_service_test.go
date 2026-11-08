package service

import (
	"errors"
	"testing"
	"testing/synctest"
	"time"
	"uuid"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/fixed"
	"github.com/liaradb/liaradb/collection/idempotency"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/outbox"
	"github.com/liaradb/liaradb/collection/record"
	"github.com/liaradb/liaradb/collection/tenant"
	"github.com/liaradb/liaradb/domain/command"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/transaction"
	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/util/testing/filetesting"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestEventService_Append(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventService_Append)
}

func testEventService_Append(t *testing.T) {
	m, s := createManager(t, 1024, 1024)
	c := btree.NewCursor(s)
	fc := fixed.New(record.New(s), c)
	ts := tenant.New(fc)
	es := NewEventService(m,
		ts,
		idempotency.New(fc, c),
		eventlog.New(fc, c),
		keyvalue.New(fc),
		outbox.New(fc))

	tid := initTenant(t, m, ts)
	aggregateID := value.NewAggregateID(uuid.NewV7().String())
	if err := es.Append(t.Context(), command.AppendEvent{
		TenantID: tid,
		Options:  command.AppendOptions{},
		Events: []command.EventOptions{{
			// ID:         value.NewEventID(),
			AggregateName: value.NewAggregateName("example"),
			AggregateID:   aggregateID,
			Version:       value.NewVersion(1),
		}},
	}); err != nil {
		t.Fatal(err)
	}

	if err := es.Append(t.Context(), command.AppendEvent{
		TenantID: tid,
		Options:  command.AppendOptions{},
		Events: []command.EventOptions{{
			// ID:         value.NewEventID(),
			AggregateName: value.NewAggregateName("example"),
			AggregateID:   aggregateID,
			Version:       value.NewVersion(2),
		}},
	}); err != nil {
		t.Fatal(err)
	}

	if err := es.Append(t.Context(), command.AppendEvent{
		TenantID: tid,
		Options:  command.AppendOptions{},
		Events: []command.EventOptions{{
			// ID:         value.NewEventID(),
			AggregateName: value.NewAggregateName("example"),
			AggregateID:   aggregateID,
			Version:       value.NewVersion(3),
		}},
	}); err != nil {
		t.Fatal(err)
	}

	synctest.Wait()
}

func TestEventService_Append__Large(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventService_Append__Large)
}

func testEventService_Append__Large(t *testing.T) {
	m, s := createManager(t, 1024, 1024)
	c := btree.NewCursor(s)
	fc := fixed.New(record.New(s), c)
	ts := tenant.New(fc)

	es := NewEventService(m,
		ts,
		idempotency.New(fc, c),
		eventlog.New(fc, c),
		keyvalue.New(fc),
		outbox.New(fc))

	tid := initTenant(t, m, ts)
	aggregateID := value.NewAggregateID(uuid.NewV7().String())
	if err := es.Append(t.Context(), command.AppendEvent{
		TenantID: tid,
		Options:  command.AppendOptions{},
		Events: []command.EventOptions{{
			// ID:         value.NewEventID(),
			AggregateName: value.NewAggregateName("example"),
			AggregateID:   aggregateID,
			Version:       value.NewVersion(1),
		}},
	}); err != nil {
		t.Fatal(err)
	}

	if err := es.Append(t.Context(), command.AppendEvent{
		TenantID: tid,
		Options:  command.AppendOptions{},
		Events: []command.EventOptions{{
			// ID:         value.NewEventID(),
			AggregateName: value.NewAggregateName("example"),
			AggregateID:   aggregateID,
			Version:       value.NewVersion(2),
		}},
	}); err != nil {
		t.Fatal(err)
	}

	if err := es.Append(t.Context(), command.AppendEvent{
		TenantID: tid,
		Options:  command.AppendOptions{},
		Events: []command.EventOptions{{
			// ID:         value.NewEventID(),
			AggregateName: value.NewAggregateName("example"),
			AggregateID:   aggregateID,
			Version:       value.NewVersion(3),
		}},
	}); err != nil {
		t.Fatal(err)
	}

	synctest.Wait()
}

func TestEventService_Append__Invalid(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventService_Append__Invalid)
}

func testEventService_Append__Invalid(t *testing.T) {
	es := NewEventService(nil, nil, nil, nil, nil, nil)

	aggregateID := value.NewAggregateID(uuid.NewV7().String())
	want := command.EventOptions{
		// ID:         value.NewEventID(),
		AggregateName: value.NewAggregateName("example"),
		AggregateID:   aggregateID,
		Version:       value.NewVersion(0),
	}

	if err := es.Append(t.Context(), command.AppendEvent{
		TenantID: value.NewTenantID(),
		Options:  command.AppendOptions{},
		Events:   []command.EventOptions{want}},
	); !errors.Is(err, value.ErrAggregateVersionInvalid) {
		t.Error("should return error")
	}

	synctest.Wait()
}

func initTenant(t *testing.T, m *transaction.Manager, ts *tenant.Tenant) value.TenantID {
	tx, err := m.Next(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	tid := value.NewTenantID()
	if err := tx.Run(t.Context(), func() error {
		return ts.Insert(t.Context(), tx, entity.NewTenant(tid, value.NewTenantName("")))
	}); err != nil {
		t.Fatal(err)
	}
	return tid
}

func createManager(t *testing.T, pageSize int16, maxRecordSize int64) (*transaction.Manager, *storage.Storage) {
	t.Helper()

	fsys, dir := createFiles()
	l := createLog(t, fsys, dir, pageSize, maxRecordSize)
	s := storagetesting.CreateStorageWithFileSystem(t, 2, 1024, fsys, dir)

	m := transaction.NewManager(l, s, 1)
	m.Run(t.Context())

	t.Cleanup(func() {
		m.Shutdown(t.Context(), time.Now())
	})

	return m, s
}

func createLog(t *testing.T, fsys filecache.FileSystem, dir string, pageSize int16, maxRecordSize int64) *log.Log {
	t.Helper()

	l := log.New(pageSize, 10, maxRecordSize, 100, fsys, dir)
	if err := l.Run(t.Context()); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := l.Close(); err != nil {
			t.Error(err)
		}
	})

	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	return l
}

func createFiles() (filecache.FileSystem, string) {
	// return &disk.FileSystem{}, t.TempDir()
	return filetesting.New(nil), "."
}
