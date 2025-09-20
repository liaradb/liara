package log

import (
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/liaradb/liaradb/file/mock"
)

func TestLogReader_Iterate(t *testing.T) {
	t.Parallel()

	lr, l := createLogReaderWriter(t)

	var count LogSequenceNumber = 10
	records, lsn := createRecords(count)

	for _, r := range records {
		_, err := l.Append(r)
		if err != nil {
			t.Error(err)
		}
	}

	if err := l.Flush(lsn); err != nil {
		t.Error(err)
	}

	var c LogSequenceNumber
	for r, err := range lr.Iterate() {
		c++
		if err != nil {
			t.Fatal(err)
		}

		record := records[c-1]

		if !reflect.DeepEqual(r, record) {
			t.Fatalf("incorrect value:\n%#v, expected:\n%#v", r, record)
		}
	}

	if c != count {
		t.Errorf("incorrect count: %v, expected: %v", c, count)
	}
}

func TestLogReader_Reverse(t *testing.T) {
	lr, l := createLogReaderWriter(t)

	var count LogSequenceNumber = 10
	records, lsn := createRecords(count)

	for _, r := range records {
		_, err := l.Append(r)
		if err != nil {
			t.Error(err)
		}
	}

	if err := l.Flush(lsn); err != nil {
		t.Error(err)
	}

	var c LogSequenceNumber
	for r, err := range lr.Reverse() {
		c++
		if err != nil {
			t.Fatal(err)
		}

		record := records[count-c]

		if !reflect.DeepEqual(r, record) {
			t.Fatalf("incorrect value:\n%#v, expected:\n%#v", r, record)
		}
	}

	if c != count {
		t.Errorf("incorrect count: %v, expected: %v", c, count)
	}
}

func createLogReaderWriter(t *testing.T) (*LogReader, *LogWriter) {
	t.Helper()

	f := mock.NewMockFile(path.Join(t.TempDir(), "logfile"))
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))

	return NewLogReader(256, f), NewLogWriter(256, f)
}

func createRecords(count LogSequenceNumber) ([]*Record, LogSequenceNumber) {
	records := make([]*Record, 0, count)
	for i := range count {
		records = append(records, newRecord(i, 2, time.UnixMicro(1234567890), data, reverse))
	}
	return records, count - 1
}
