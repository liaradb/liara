package service

import (
	"errors"
	"testing"
	"testing/synctest"
	"time"

	"github.com/google/uuid"
	"github.com/liaradb/liaradb/domain/command"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/filecache"
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
	m, _ := createManager(t, 256)
	es := NewEventService(m)

	aggregateID := value.NewAggregateID(uuid.NewString())
	if err := es.Append(t.Context(), command.AppendEvent{
		TenantID:    value.NewTenantID(),
		Options:     command.AppendOptions{},
		PartitionID: value.NewPartitionID(0),
		Events: []command.EventOptions{{
			AggregateName: value.NewAggregateName("example"),
			// ID:            value.NewEventID(),
			AggregateID: aggregateID,
			Version:     value.NewVersion(1),
		}},
	}); err != nil {
		t.Fatal(err)
	}

	if err := es.Append(t.Context(), command.AppendEvent{
		TenantID:    value.NewTenantID(),
		Options:     command.AppendOptions{},
		PartitionID: value.NewPartitionID(0),
		Events: []command.EventOptions{{
			AggregateName: value.NewAggregateName("example"),
			// ID:            value.NewEventID(),
			AggregateID: aggregateID,
			Version:     value.NewVersion(2),
		}},
	}); err != nil {
		t.Fatal(err)
	}

	if err := es.Append(t.Context(), command.AppendEvent{
		TenantID:    value.NewTenantID(),
		Options:     command.AppendOptions{},
		PartitionID: value.NewPartitionID(0),
		Events: []command.EventOptions{{
			AggregateName: value.NewAggregateName("example"),
			// ID:            value.NewEventID(),
			AggregateID: aggregateID,
			Version:     value.NewVersion(3),
		}},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestEventService_Append__Large(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventService_Append__Large)
}

func testEventService_Append__Large(t *testing.T) {
	m, _ := createManager(t, 1024)
	es := NewEventService(m)

	aggregateID := value.NewAggregateID(uuid.NewString())
	if err := es.Append(t.Context(), command.AppendEvent{
		TenantID:    value.NewTenantID(),
		Options:     command.AppendOptions{},
		PartitionID: value.NewPartitionID(0),
		Events: []command.EventOptions{{
			AggregateName: value.NewAggregateName("example"),
			// ID:            value.NewEventID(),
			AggregateID: aggregateID,
			Version:     value.NewVersion(1),
		}},
	}); err != nil {
		t.Fatal(err)
	}

	if err := es.Append(t.Context(), command.AppendEvent{
		TenantID:    value.NewTenantID(),
		Options:     command.AppendOptions{},
		PartitionID: value.NewPartitionID(0),
		Events: []command.EventOptions{{
			AggregateName: value.NewAggregateName("example"),
			// ID:            value.NewEventID(),
			AggregateID: aggregateID,
			Version:     value.NewVersion(2),
		}},
	}); err != nil {
		t.Fatal(err)
	}

	if err := es.Append(t.Context(), command.AppendEvent{
		TenantID:    value.NewTenantID(),
		Options:     command.AppendOptions{},
		PartitionID: value.NewPartitionID(0),
		Events: []command.EventOptions{{
			AggregateName: value.NewAggregateName("example"),
			// ID:            value.NewEventID(),
			AggregateID: aggregateID,
			Version:     value.NewVersion(3),
		}},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestEventService_Append__Invalid(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventService_Append__Invalid)
}

func testEventService_Append__Invalid(t *testing.T) {
	es := NewEventService(nil)

	aggregateID := value.NewAggregateID(uuid.NewString())
	want := command.EventOptions{
		AggregateName: value.NewAggregateName("example"),
		// ID:            value.NewEventID(),
		AggregateID: aggregateID,
		Version:     value.NewVersion(0),
	}

	if err := es.Append(t.Context(), command.AppendEvent{
		TenantID:    value.NewTenantID(),
		Options:     command.AppendOptions{},
		PartitionID: value.NewPartitionID(0),
		Events:      []command.EventOptions{want}},
	); !errors.Is(err, value.ErrAggregateVersionInvalid) {
		t.Error("should return error")
	}
}

func createManager(t *testing.T, pageSize int64) (*transaction.Manager, *log.Log) {
	t.Helper()

	fsys, dir := createFiles()
	l := createLog(t, fsys, dir, pageSize)
	s := storagetesting.CreateStorageWithFileSystem(t, 2, 1024, fsys, dir)

	m := transaction.NewManager(l, s, 1)
	m.Run(t.Context())

	t.Cleanup(func() {
		m.Shutdown(t.Context(), time.Now())
	})

	return m, l
}

func createLog(t *testing.T, fsys filecache.FileSystem, dir string, pageSize int64) *log.Log {
	t.Helper()

	l := log.New(pageSize, 10, pageSize, 100, fsys, dir)
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
