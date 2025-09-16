package log

import (
	"path"
	"testing"
	"time"

	"github.com/liaradb/liaradb/mock"
)

var data = []byte{0, 1, 2, 3, 4, 5}
var reverse = []byte{6, 7, 8, 9, 10, 11}
var record = newLogRecord(1, 2, time.UnixMicro(1234567890), data, reverse)

func TestLogWriter_Default(t *testing.T) {
	t.Parallel()

	l := createLogWriter(t)

	testPosition(t, l, 0, 0)
}

func TestLogWriter_Append(t *testing.T) {
	t.Parallel()

	l := createLogWriter(t)

	if lsn, err := l.Append(record); err != nil {
		t.Error(err)
	} else if lsn != 1 {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, l, 0, 1)
}

func TestLogWriter_Flush(t *testing.T) {
	t.Parallel()

	t.Run("should flush", func(t *testing.T) {
		t.Parallel()

		l := createLogWriter(t)

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

		l := createLogWriter(t)

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

func createLogWriter(t *testing.T) *LogWriter {
	t.Helper()

	f := mock.NewMockFile(path.Join(t.TempDir(), "logfile"))
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))

	return NewLogWriter(256, f)
}

func testPosition(t *testing.T, l *LogWriter, lw, hw LogSequenceNumber) {
	if h := l.HighWater(); h != hw {
		t.Errorf("incorrect value: %v, expected: %v", h, hw)
	}

	if l := l.LowWater(); l != lw {
		t.Errorf("incorrect value: %v, expected: %v", l, lw)
	}
}
