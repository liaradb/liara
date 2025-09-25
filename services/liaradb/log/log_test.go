package log

import (
	"reflect"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/liaradb/liaradb/file/mock"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

func TestLog_EmptyReader(t *testing.T) {
	t.Parallel()

	fsys := createFiles(0, 0)

	l := NewLog(256, 2, fsys, ".")
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

	l := NewLog(256, 2, fsys, ".")
	if err := l.Open(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := l.Close(); err != nil {
			t.Error(err)
		}
	}()

	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	records, _ := createRecords(100)
	var lsn record.LogSequenceNumber
	var err error
	for _, rec := range records {
		lsn, err = l.Append(rec)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err = l.Flush(lsn); err != nil {
		t.Fatal(err)
	}

	i := 0
	for rc, err := range l.Iterate(0) {
		if err != nil {
			t.Fatal(err)
		}

		rec := records[i]

		if !reflect.DeepEqual(rc, rec) {
			t.Error("records do not match")
		}
		i++
	}
	if i != 100 {
		t.Errorf("incorrect count: %v, expected: %v", i, 100)
	}
}

func TestLog_Reverse(t *testing.T) {
	t.Parallel()

	fsys := createFiles(0, 0)

	l := NewLog(256, 2, fsys, ".")
	if err := l.Open(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := l.Close(); err != nil {
			t.Error(err)
		}
	}()

	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	records, _ := createRecords(100)
	var lsn record.LogSequenceNumber
	var err error
	for _, rec := range records {
		lsn, err = l.Append(rec)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err = l.Flush(lsn); err != nil {
		t.Fatal(err)
	}

	slices.Reverse(records)
	i := 0
	for rc, err := range l.Reverse() {
		if err != nil {
			t.Fatal(err)
		}

		rec := records[i]

		if !reflect.DeepEqual(rc, rec) {
			t.Error("records do not match")
		}
		i++
	}
	if i != 100 {
		t.Errorf("incorrect count: %v, expected: %v", i, 100)
	}
}

func createFiles(start segment.SegmentID, count segment.SegmentID) *mock.FileSystem {
	fsys := &mock.FileSystem{MapFS: fstest.MapFS{}}
	for i := range count {
		fsys.MapFS[segment.NewSegmentName(start+i, 0).String()] = &fstest.MapFile{}
	}
	return fsys
}
