package log

import (
	"path"
	"reflect"
	"testing"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/mock"
	"github.com/liaradb/liaradb/log/record"
)

func TestPageReader_Iterate(t *testing.T) {
	t.Parallel()

	_, lr, lw := createPageReaderWriter(t)

	var count record.LogSequenceNumber = 3
	records, lsn := createRecords(count)

	for _, rc := range records {
		_, err := lw.Append(rc)
		if err != nil {
			t.Error(err)
		}
	}

	if err := lw.Flush(lsn); err != nil {
		t.Error(err)
	}

	var c record.LogSequenceNumber
	for rc, err := range lr.Iterate() {
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
	f, sr, lw := createPageReaderWriter(t)

	var count record.LogSequenceNumber = 3
	records, lsn := createRecords(count)

	for _, rc := range records {
		_, err := lw.Append(rc)
		if err != nil {
			t.Error(err)
		}
	}

	if err := lw.Flush(lsn); err != nil {
		t.Error(err)
	}

	stat, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	var c record.LogSequenceNumber
	for rc, err := range sr.Reverse(stat.Size()) {
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

func createPageReaderWriter(t *testing.T) (file.File, *PageReader, *LogWriter) {
	t.Helper()

	fsys := mock.NewFileSystem(nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))

	lw := NewLogWriter(256, 3, f)
	_ = lw.Initialize()
	return f, NewPageReader(256, f), lw
}
