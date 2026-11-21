package transaction

import (
	"context"
	"iter"
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/transaction/record"
	"github.com/liaradb/liaradb/util/testing/filetesting"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestManager(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testManager)
}

func testManager(t *testing.T) {
	m, _ := createManager(t)
	ctx := t.Context()

	tx0, err := m.Next(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var tid0 = record.NewTransactionID(1)
	if i := tx0.ID(); i != tid0 {
		t.Errorf("id does not match: %v, expected: %v", i, tid0)
	}

	tx1, err := m.Next(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var tid1 = record.NewTransactionID(2)
	if i := tx1.ID(); i != tid1 {
		t.Errorf("id does not match: %v, expected: %v", i, tid1)
	}
}

func createManager(t *testing.T) (*Manager, *log.Log) {
	t.Helper()

	fsys, dir := createFiles()
	l := createLog(t, fsys, dir)
	s := storagetesting.CreateStorageWithFileSystem(t, 2, 1024, fsys, dir)
	m := NewManager(l, s, 1)
	m.Run(t.Context(), func(err error) {})

	return m, l
}

func createManagerIterator(t *testing.T) (*Manager, *log.Log, func() iter.Seq[*record.Record]) {
	t.Helper()

	fsys, dir := createFiles()
	l, cancel := createLogCancel(t, fsys, dir)

	s := storagetesting.CreateStorageWithFileSystem(t, 2, 1024, fsys, dir)
	m := NewManager(l, s, 1)

	ctx, cancel2 := context.WithCancel(t.Context())
	m.Run(ctx, func(err error) {})

	return m, l, func() iter.Seq[*record.Record] {
		if err := l.Flush(t.Context()); err != nil {
			t.Fatal(err)
		}

		cancel2()
		cancel()
		return iterateLog(t, fsys, dir)
	}
}

func createLog(t *testing.T, fsys filecache.FileSystem, dir string) *log.Log {
	t.Helper()

	l, cancel := createLogCancel(t, fsys, dir)

	t.Cleanup(cancel)

	return l
}

func createLogCancel(t *testing.T, fsys filecache.FileSystem, dir string) (*log.Log, func()) {
	t.Helper()

	ctx, cancel := context.WithCancel(t.Context())

	l := log.New(256, 3, 256, 100, fsys, dir)
	if err := l.Run(ctx, func(err error) {}); err != nil {
		t.Fatal(err)
	}

	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	return l, func() {
		cancel()
		synctest.Wait()
		if err := l.Close(); err != nil {
			t.Error(err)
		}
	}
}

func iterateLog(t *testing.T, fsys filecache.FileSystem, dir string) iter.Seq[*record.Record] {
	t.Helper()

	l := log.New(256, 3, 256, 100, fsys, dir)

	it, err := l.Recover()
	if err != nil {
		t.Fatal(err)
	}

	return it
}

func createFiles() (filecache.FileSystem, string) {
	// return &disk.FileSystem{}, t.TempDir()
	return filetesting.New(nil), "."
}

func TestManager_Active(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testManager_ActiveTestManager_Active)
}

func testManager_ActiveTestManager_Active(t *testing.T) {
	m, _ := createManager(t)
	ctx := t.Context()

	tx0, err := m.Next(ctx)
	if err != nil {
		t.Fatal(err)
	}

	active := m.Active()

	if l := len(active); l != 1 {
		t.Errorf("incorrect length: %v, expected: %v", l, 1)
	}

	if !slices.Contains(active, tx0.ID()) {
		t.Errorf("should include: %v", tx0.ID())
	}

	if err := tx0.Run(ctx, func() error {
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	synctest.Wait()

	active = m.Active()

	if l := len(active); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}
}
