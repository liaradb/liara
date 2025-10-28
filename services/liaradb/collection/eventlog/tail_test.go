package eventlog

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/file/filetesting"
	"github.com/liaradb/liaradb/storage"
)

func TestTail(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTail)
}

func testTail(t *testing.T) {
	s := createStorage(t, 2, 16)
	tl := NewTail(s)
	if err := tl.Append(t.Context()); err != nil {
		t.Error(err)
	}
}

func createStorage(t *testing.T, max int, bs int64) *storage.Storage {
	fsys := filetesting.NewMockFileSystem(t, nil)
	s := storage.New(fsys, max, bs, t.TempDir())

	if err := s.Run(t.Context()); err != nil {
		t.Fatal(err)
	}

	return s
}
