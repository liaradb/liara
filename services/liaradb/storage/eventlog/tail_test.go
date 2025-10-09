package eventlog

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/filetesting"
	"github.com/liaradb/liaradb/storage"
)

func TestTail(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTail)
}

func testTail(t *testing.T) {
	s := createStorage(t)
	tl := NewTail(s)
	if err := tl.Append(t.Context()); err != nil {
		t.Error(err)
	}
}

func createStorage(t *testing.T) *storage.Storage {
	fsys := filetesting.NewMockFileSystem(t, nil)
	s := storage.NewStorage(fsys, 2, 16)

	s.Run(t.Context())
	return s
}
