package transaction

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/filetesting"
	"github.com/liaradb/liaradb/log"
	"github.com/liaradb/liaradb/log/record"
)

func TestTransaction_Insert(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTransaction_Insert)
}

func testTransaction_Insert(t *testing.T) {
	m, l := createManager(t)

	tx := m.Next()

	var tid record.TransactionID = 1
	if i := tx.ID(); i != tid {
		t.Errorf("id does not match: %v, expected: %v", i, tid)
	}

	ctx := t.Context()
	if err := tx.Insert(ctx, time.UnixMicro(1234567890), nil); err != nil {
		t.Fatal(err)
	}

	if err := l.Flush(ctx, tx.LogSequenceNumber()); err != nil {
		t.Fatal(err)
	}

	c := 0
	for rc, err := range l.Iterate(0) {
		if err != nil {
			t.Fatal(err)
		}

		if lsn := tx.LogSequenceNumber(); lsn != rc.LogSequenceNumber() {
			t.Errorf("lsn does not match: %v, expected: %v", lsn, rc.LogSequenceNumber())
		}

		c++
	}

	if c != 1 {
		t.Errorf("incorrect record count: %v, expected: %v", c, 1)
	}
}

func TestTransaction_Commit(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTransaction_Commit)
}

func testTransaction_Commit(t *testing.T) {
	m, l := createManager(t)

	tx := m.Next()

	ctx := t.Context()
	if err := tx.Insert(ctx, time.UnixMicro(1234567890), nil); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(ctx, time.UnixMicro(1234567890)); err != nil {
		t.Fatal(err)
	}

	lsns := []record.LogSequenceNumber{1, 2}
	actions := []record.Action{record.ActionInsert, record.ActionCommit}

	c := 0
	for rc, err := range l.Iterate(0) {
		if err != nil {
			t.Fatal(err)
		}

		if lsn := rc.LogSequenceNumber(); lsn != lsns[c] {
			t.Errorf("lsn does not match: %v, expected: %v", lsn, lsns[c])
		}

		if a := rc.Action(); a != actions[c] {
			t.Errorf("action does not match: %v, expected: %v", a, actions[c])
		}

		c++
	}

	if c != 2 {
		t.Errorf("incorrect record count: %v, expected: %v", c, 2)
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
