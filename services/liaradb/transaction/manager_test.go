package transaction

import (
	"testing"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/filetesting"
	"github.com/liaradb/liaradb/log"
	"github.com/liaradb/liaradb/log/record"
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

	l := createLog(t)
	return NewManager(l), l
}

func createLog(t *testing.T) *log.Log {
	t.Helper()

	fsys, dir := createFiles(t)
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

func createFiles(t *testing.T) (file.FileSystem, string) {
	// return &disk.FileSystem{}, t.TempDir()
	return filetesting.NewMockFileSystem(t, nil), "."
}
