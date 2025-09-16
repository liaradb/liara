package log

import (
	"path"
	"reflect"
	"testing"

	"github.com/liaradb/liaradb/mock"
)

func TestLogReader_Iterate(t *testing.T) {
	t.Parallel()

	lr, l := createLogReaderWriter(t)

	count := 10

	for range count - 1 {
		_, err := l.Append(record)
		if err != nil {
			t.Error(err)
		}
	}

	lsn2, err := l.Append(record)
	if err != nil {
		t.Error(err)
	}

	if err := l.Flush(lsn2); err != nil {
		t.Error(err)
	}

	if p := l.PageIndex(); p != 2 {
		t.Errorf("incorrect value: %v, expected: %v", p, 2)
	}

	c := 0
	for r, err := range lr.Iterate() {
		c++
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(r, record) {
			t.Fatalf("incorrect value:\n%#v, expected:\n%#v", r, record)
		}
	}

	if c != count {
		t.Errorf("incorrect count: %v, expected: %v", c, count)
	}
}

func createLogReaderWriter(t *testing.T) (*LogReader, *Log) {
	t.Helper()

	f := mock.NewMockFile(path.Join(t.TempDir(), "logfile"))
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))
	l := &Log{
		pageSize: 256,
	}
	l.Open(f)

	lr := &LogReader{
		pageSize: 256,
	}
	lr.Open(f)

	return lr, l
}
