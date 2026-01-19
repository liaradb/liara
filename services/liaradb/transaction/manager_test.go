package transaction

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/filetesting"
	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/mempage"
	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/storage/storagetesting"
)

func TestManager(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testManager)
}

func testManager(t *testing.T) {
	m, _ := createManager(t)

	tx0 := m.Next()

	var tid0 = record.NewTransactionID(1)
	if i := tx0.ID(); i != tid0 {
		t.Errorf("id does not match: %v, expected: %v", i, tid0)
	}

	tx1 := m.Next()

	var tid1 = record.NewTransactionID(2)
	if i := tx1.ID(); i != tid1 {
		t.Errorf("id does not match: %v, expected: %v", i, tid1)
	}
}

func createManager(t *testing.T) (*Manager, *recovery.Log) {
	t.Helper()

	fsys, dir := createFiles(t)
	l := createLog(t, fsys, dir)
	s := storagetesting.CreateStorageWithFileSystem(t, 2, 1024, fsys)
	lt := createLockTable(t)
	return NewManager(l, s, lt), l
}

func createLog(t *testing.T, fsys file.FileSystem, dir string) *recovery.Log {
	t.Helper()

	l := recovery.NewLog(256, 3, fsys, dir, mempage.NewWithHeader(256, &page.Header{}))
	if err := l.Open(t.Context()); err != nil {
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
	lt := locktable.NewLockTable[action.ItemID](1)
	lt.Run(t.Context())
	t.Cleanup(lt.Close)
	return lt
}

func createFiles(t *testing.T) (file.FileSystem, string) {
	// return &disk.FileSystem{}, t.TempDir()
	return filetesting.NewMockFileSystem(t, nil), "."
}
