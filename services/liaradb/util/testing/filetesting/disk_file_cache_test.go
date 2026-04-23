package filetesting

import (
	"testing"

	"github.com/liaradb/liaradb/file/disk"
)

func TestNewDiskFileCache(t *testing.T) {
	t.Parallel()

	fc := NewDiskFileCache(t)
	if fc == nil {
		t.Error("should return value")
	}

	if fsys, ok := fc.FSYS().(*disk.FileSystem); !ok {
		t.Error("incorrect type")
	} else if fsys == nil {
		t.Error("should return value")
	}
}
