package filetesting

import (
	"testing"

	"github.com/liaradb/liaradb/filecache"
)

func NewDiskFileCache(t *testing.T) *filecache.Cache {
	t.Helper()

	fsys := filecache.New()
	t.Cleanup(func() {
		if err := fsys.Close(); err != nil {
			t.Error(err)
		}
	})
	return fsys
}
