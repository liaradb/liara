package log

import (
	"io"
	"path"
	"reflect"
	"testing"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/mock"
	"github.com/liaradb/liaradb/log/record"
)

func TestPageReader_Iterate(t *testing.T) {
	t.Parallel()

	f, pr, sw := createPageReaderWriter(t)

	var count record.LogSequenceNumber = 3
	records, _ := createRecords(count)

	for _, rc := range records {
		if err := sw.Append(rc); err != nil {
			t.Error(err)
		}
	}

	if err := sw.Flush(); err != nil {
		t.Error(err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	var c record.LogSequenceNumber
	it, err := pr.Iterate(f)
	if err != nil {
		t.Fatal(err)
	}

	for rc, err := range it {
		c++
		if err != nil {
			t.Fatal(err)
		}

		rec := records[c-1]

		if !reflect.DeepEqual(rc, rec) {
			t.Fatalf("incorrect value:\n%#v, expected:\n%#v", rc, rec)
		}
	}

	if c != count {
		t.Errorf("incorrect count: %v, expected: %v", c, count)
	}
}

func TestPageReader_Reverse(t *testing.T) {
	f, pr, sw := createPageReaderWriter(t)

	var count record.LogSequenceNumber = 3
	records, _ := createRecords(count)

	for _, rc := range records {
		if err := sw.Append(rc); err != nil {
			t.Error(err)
		}
	}

	if err := sw.Flush(); err != nil {
		t.Error(err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	var c record.LogSequenceNumber
	it, err := pr.Reverse(f)
	if err != nil {
		t.Fatal(err)
	}

	for rc, err := range it {
		c++
		if err != nil {
			t.Fatal(err)
		}

		rec := records[count-c]

		if !reflect.DeepEqual(rc, rec) {
			t.Fatalf("incorrect value:\n%#v, expected:\n%#v", rc, rec)
		}
	}

	if c != count {
		t.Errorf("incorrect count: %v, expected: %v", c, count)
	}
}

func createPageReaderWriter(t *testing.T) (file.File, *PageReader, *SegmentWriter) {
	t.Helper()

	fsys := mock.NewFileSystem(nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))

	sw := NewSegmentWriter(256, 3, f)
	_ = sw.Initialize()
	return f, NewPageReader(256), sw
}
