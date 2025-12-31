package storagetesting

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/filetesting"
	"github.com/liaradb/liaradb/storage"
)

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
