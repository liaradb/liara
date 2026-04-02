package storagetesting

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/util/testing/filetesting"
)

type Storage struct {
	Storage *storage.Storage
	FSys    file.FileSystem
}

func SyncTest(t *testing.T, max int, bs int64, f func(*testing.T, Storage)) {
	t.Helper()
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		fsys := filetesting.NewMockFileSystem(t, nil)
		s := CreateStorageWithFileSystem(t, max, bs, fsys)

		t.Cleanup(func() {
			if p := s.CountPinned(); p != 0 {
				t.Errorf("incorrect pin count: %v, expected: %v", p, 0)
			}
		})

		f(t, Storage{
			Storage: s,
			FSys:    fsys,
		})

		synctest.Wait()
	})
}

func CreateStorage(t *testing.T, max int, bs int64) *storage.Storage {
	t.Helper()

	fsys := filetesting.NewMockFileSystem(t, nil)
	return CreateStorageWithFileSystem(t, max, bs, fsys)
}

func CreateStorageWithFileSystem(t *testing.T, max int, bs int64, fsys file.FileSystem) *storage.Storage {
	t.Helper()

	s := storage.New(fsys, max, bs, t.TempDir())
	if err := s.Run(t.Context()); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		synctest.Wait()

		if p := s.CountPinned(); p != 0 {
			t.Errorf("incorrect pin count: %v, expected: %v", p, 0)
		}
	})

	return s
}
