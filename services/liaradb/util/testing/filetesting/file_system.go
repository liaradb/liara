package filetesting

import (
	"testing"
	"testing/fstest"
	"time"

	"github.com/liaradb/liaradb/file/disk"
	"github.com/liaradb/liaradb/file/mock"
)

func NewDiskFileSystem(t *testing.T) *disk.FileSystem {
	fsys := &disk.FileSystem{}
	t.Cleanup(func() {
		if err := fsys.Close(); err != nil {
			t.Error(err)
		}
	})
	return fsys
}

func NewMockFileSystem(t *testing.T, fsys fstest.MapFS) *mock.FileSystem {
	return mock.NewFileSystem(fsys)
}

func NewMockFileSystemDelay(t *testing.T, fsys fstest.MapFS, delay time.Duration) *mock.FileSystem {
	return mock.NewFileSystemDelay(fsys, delay)
}
