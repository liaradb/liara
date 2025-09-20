package log

import (
	"testing"
	"testing/fstest"

	"github.com/liaradb/liaradb/file/mock"
	"github.com/liaradb/liaradb/log/segment"
)

func TestLog_EmptyReader(t *testing.T) {
	t.Parallel()

	fsys := createFiles(0, 0)

	l := NewLog(256, fsys, ".")
	_, err := l.Reader(0)
	if err != segment.ErrNoSegmentFile {
		t.Error("should have no files")
	}
}

func createFiles(start segment.SegmentID, count segment.SegmentID) *mock.FileSystem {
	fsys := &mock.FileSystem{MapFS: fstest.MapFS{}}
	for i := range count {
		fsys.MapFS[segment.NewSegmentName(start+i, 0).String()] = &fstest.MapFile{}
	}
	return fsys
}
