package log

import (
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/liaradb/liaradb/file/mock"
	"github.com/liaradb/liaradb/log/record"
)

func TestLogReader_Iterate(t *testing.T) {
	t.Parallel()

	lr, lw := createLogReaderWriter(t)

	var count record.LogSequenceNumber = 10
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

func TestLogReader_Reverse(t *testing.T) {
	lr, lw := createLogReaderWriter(t)

	var count record.LogSequenceNumber = 10
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
	for rc, err := range lr.Reverse() {
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

func createLogReaderWriter(t *testing.T) (*LogReader, *LogWriter) {
	t.Helper()

	fsys := mock.NewFileSystem(nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))

	return NewLogReader(256, f), NewLogWriter(256, 2, f)
}

func createRecords(count record.LogSequenceNumber) ([]*record.Record, record.LogSequenceNumber) {
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	records := make([]*record.Record, 0, count)
	for i := range count {
		records = append(records, record.NewRecord(i, 2, time.UnixMicro(1234567890), data, reverse))
	}
	return records, count - 1
}
