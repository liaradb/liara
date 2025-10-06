package transaction

import (
	"testing"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/filetesting"
	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/log"
	"github.com/liaradb/liaradb/log/action"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/storage"
)

func TestManager(t *testing.T) {
	m, _ := createManager(t)

	tx0 := m.Next()

	var tid0 record.TransactionID = 1
	if i := tx0.ID(); i != tid0 {
		t.Errorf("id does not match: %v, expected: %v", i, tid0)
	}

	tx1 := m.Next()

	var tid1 record.TransactionID = 2
	if i := tx1.ID(); i != tid1 {
		t.Errorf("id does not match: %v, expected: %v", i, tid1)
	}
}

func createManager(t *testing.T) (*Manager, *log.Log) {
	t.Helper()

	fsys, dir := createFiles(t)
	l := createLog(t, fsys, dir)
	s := createStorage(t, fsys)
	lt := createLockTable(t)
	c := createConcurrencyMgr(lt)
	return NewManager(l, s, c), l
}

func createLog(t *testing.T, fsys file.FileSystem, dir string) *log.Log {
	t.Helper()

	l := log.NewLog(256, 3, fsys, dir)
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

func createStorage(t *testing.T, fsys file.FileSystem) *storage.Storage {
	s := storage.NewStorage(fsys, 2, 1024)
	s.Run(t.Context())
	return s
}

func createLockTable(t *testing.T) *locktable.LockTable[action.ItemID] {
	lt := locktable.NewLockTable[action.ItemID](t.Context(), 1)
	t.Cleanup(lt.Close)
	return lt
}

func createConcurrencyMgr(lt *locktable.LockTable[action.ItemID]) *locktable.ConcurrencyMgr[action.ItemID] {
	return locktable.NewConcurrencyMgr(lt)
}

func createFiles(t *testing.T) (file.FileSystem, string) {
	// return &disk.FileSystem{}, t.TempDir()
	return filetesting.NewMockFileSystem(t, nil), "."
}
