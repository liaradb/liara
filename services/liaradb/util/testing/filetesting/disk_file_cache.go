package filetesting

import (
	"testing"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/disk"
)

func NewDiskFileCache(t *testing.T) *file.Cache {
	fsys := disk.New()
	t.Cleanup(func() {
		if err := fsys.Close(); err != nil {
			t.Error(err)
		}
	})
	return fsys
}
