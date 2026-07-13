package service

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"

	"github.com/google/uuid"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/transaction"
	"github.com/liaradb/liaradb/util/testing/filetesting"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestEventService_Append(t *testing.T) {
	t.Parallel()
	t.Skip()
	synctest.Test(t, testEventService_Append)
}

func testEventService_Append(t *testing.T) {
	m, l := createManager(t)
	l.Run(t.Context())
	l.StartWriter()
	m.Run(t.Context())
	es := NewEventService(m)

	aggregateID := value.NewAggregateID(uuid.NewString())
	if err := es.Append(
		t.Context(),
		value.NewTenantID(),
		AppendOptions{},
		value.NewPartitionID(0),
		AppendEvent{
			AggregateName: value.NewAggregateName("example"),
			// ID:            value.NewEventID(),
			AggregateID: aggregateID,
			Version:     value.NewVersion(1),
		},
	); err != nil {
		t.Fatal(err)
	}

	if err := es.Append(
		t.Context(),
		value.NewTenantID(),
		AppendOptions{},
		value.NewPartitionID(0),
		AppendEvent{
			AggregateName: value.NewAggregateName("example"),
			// ID:            value.NewEventID(),
			AggregateID: aggregateID,
			Version:     value.NewVersion(2),
		},
	); err != nil {
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
	want := AppendEvent{
		AggregateName: value.NewAggregateName("example"),
		// ID:            value.NewEventID(),
		AggregateID: aggregateID,
		Version:     value.NewVersion(0),
	}

	err := es.Append(context.Background(), value.NewTenantID(), AppendOptions{}, value.NewPartitionID(0), want)
	if !errors.Is(err, value.ErrAggregateVersionInvalid) {
		t.Error("should return error")
	}
}

func createManager(t *testing.T) (*transaction.Manager, *recovery.Log) {
	t.Helper()

	fsys, dir := createFiles()
	l := createLog(t, fsys, dir)
	s := storagetesting.CreateStorageWithFileSystem(t, 2, 1024, fsys, dir)
	lt := createLockTable(t)
	m := transaction.NewManager(l, s, lt)
	m.Run(t.Context())

	return m, l
}

func createLog(t *testing.T, fsys filecache.FileSystem, dir string) *recovery.Log {
	t.Helper()

	l := recovery.NewLog(256, 3, 256, 100, fsys, dir)
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

func createLockTable(t *testing.T) *locktable.LockTable[action.ItemID] {
	lt := locktable.New[action.ItemID](1)
	lt.Run(t.Context())
	t.Cleanup(lt.Close)
	return lt
}

func createFiles() (filecache.FileSystem, string) {
	// return &disk.FileSystem{}, t.TempDir()
	return filetesting.New(nil), "."
}
