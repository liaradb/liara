package transaction

import (
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery"
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

	tid := value.NewTenantID()

	tx0, err := m.Next(ctx, tid)
	if err != nil {
		t.Fatal(err)
	}

	var tid0 = record.NewTransactionID(1)
	if i := tx0.ID(); i != tid0 {
		t.Errorf("id does not match: %v, expected: %v", i, tid0)
	}

	tx1, err := m.Next(ctx, tid)
	if err != nil {
		t.Fatal(err)
	}

	var tid1 = record.NewTransactionID(2)
	if i := tx1.ID(); i != tid1 {
		t.Errorf("id does not match: %v, expected: %v", i, tid1)
	}
}

func createManager(t *testing.T) (*Manager, *recovery.Log) {
	t.Helper()

	fsys, dir := createFiles()
	l := createLog(t, fsys, dir)
	s := storagetesting.CreateStorageWithFileSystem(t, 2, 1024, fsys, dir)
	m := NewManager(l, s, 1)
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

	tid := value.NewTenantID()

	tx0, err := m.Next(ctx, tid)
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

	if err := Run(ctx, tx0, time.Now(), func() error {
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
