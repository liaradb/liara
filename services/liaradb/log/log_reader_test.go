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

func createLogReaderWriter(t *testing.T) (*LogPageReader, *LogWriter) {
	t.Helper()

	f := mock.NewMockFile(path.Join(t.TempDir(), "logfile"))
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))

	return NewLogPageReader(256, f), NewLogWriter(256, f)
}
