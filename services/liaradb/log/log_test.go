package log

import (
	"reflect"
	"testing"
	"testing/fstest"
	"time"

	"github.com/liaradb/liaradb/file/mock"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

func TestLog_EmptyReader(t *testing.T) {
	t.Parallel()

	fsys := createFiles(0, 0)

	l := NewLog(256, fsys, ".")
	if err := l.Open(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := l.Close(); err != nil {
			t.Error(err)
		}
	}()

	for _, err := range l.Iterate(0) {
		if err != segment.ErrNoSegmentFile {
			t.Error("should have no files")
		}
	}
}

func TestLog_Iterate(t *testing.T) {
	t.Parallel()

	fsys := createFiles(0, 0)

	l := NewLog(256, fsys, ".")
	if err := l.Open(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := l.Close(); err != nil {
			t.Error(err)
		}
	}()

	rec := record.NewRecord(1, 2, time.UnixMicro(1234567890), []byte{1, 2, 3, 4, 5, 6}, []byte{7, 8, 9, 10, 11, 12})
	_, err := l.Append(rec)
	if err != nil {
		t.Fatal(err)
	}

	for rc, err := range l.Iterate(0) {
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(rc, rec) {
			t.Error("records do not match")
		}
	}
}

func createFiles(start segment.SegmentID, count segment.SegmentID) *mock.FileSystem {
	fsys := &mock.FileSystem{MapFS: fstest.MapFS{}}
	for i := range count {
		fsys.MapFS[segment.NewSegmentName(start+i, 0).String()] = &fstest.MapFile{}
	}
	return fsys
}
