package filetesting

import (
	"testing"

	"github.com/liaradb/liaradb/filecache"
)

func TestNewDiskFileCache(t *testing.T) {
	t.Parallel()

	ic := NewDiskFileCache(t)
	if ic == nil {
		t.Error("should return value")
	}

	if fsys, ok := ic.FSYS().(*filecache.DiskFileSystem); !ok {
		t.Error("incorrect type")
	} else if fsys == nil {
		t.Error("should return value")
	}
}
