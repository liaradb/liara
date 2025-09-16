package log

import (
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/liaradb/liaradb/mock"
)

var data = []byte{0, 1, 2, 3, 4, 5}
var reverse = []byte{6, 7, 8, 9, 10, 11}
var record = newLogRecord(1, 2, time.UnixMicro(1234567890), data, reverse)

func TestLog_Default(t *testing.T) {
	t.Parallel()

	l := createLog(t)

	testPosition(t, l, 0, 0)
}

func TestLog_Append(t *testing.T) {
	t.Parallel()

	l := createLog(t)

	if lsn, err := l.Append(record); err != nil {
		t.Error(err)
	} else if lsn != 1 {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, l, 0, 1)
}

func TestLog_Flush(t *testing.T) {
	t.Parallel()

	t.Run("should flush", func(t *testing.T) {
		t.Parallel()

		l := createLog(t)

		lsn1, err := l.Append(record)
		if err != nil {
			t.Error(err)
		}

		_, err = l.Append(record)
		if err != nil {
			t.Error(err)
		}

		if err := l.Flush(lsn1); err != nil {
			t.Error(err)
		}

		testPosition(t, l, 1, 2)
	})

	t.Run("should not flush beyond HighWater", func(t *testing.T) {
		t.Parallel()

		l := createLog(t)

		_, err := l.Append(record)
		if err != nil {
			t.Error(err)
		}

		_, err = l.Append(record)
		if err != nil {
			t.Error(err)
		}

		if err := l.Flush(10); err != nil {
			t.Error(err)
		}

		testPosition(t, l, 2, 2)
	})
}

func TestLog_Iterate(t *testing.T) {
	t.Parallel()

	l := createLog(t)

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
	for r, err := range l.Iterate() {
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

func createLog(t *testing.T) *Log {
	t.Helper()

	f := mock.NewMockFile(path.Join(t.TempDir(), "logfile"))
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))
	l := &Log{
		pageSize: 256,
	}
	l.Open(f)

	return l
}

func testPosition(t *testing.T, l *Log, lw, hw LogSequenceNumber) {
	if h := l.HighWater(); h != hw {
		t.Errorf("incorrect value: %v, expected: %v", h, hw)
	}

	if l := l.LowWater(); l != lw {
		t.Errorf("incorrect value: %v, expected: %v", l, lw)
	}
}
