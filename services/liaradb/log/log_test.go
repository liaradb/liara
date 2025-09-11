package log

import (
	"path"
	"testing"

	"github.com/liaradb/liaradb/mock"
)

var data = []byte{0, 1, 2, 3, 4, 5}

func TestLog_Default(t *testing.T) {
	t.Parallel()

	l := createLog(t)

	testPosition(t, l, 0, 0)
}

func TestLog_Append(t *testing.T) {
	t.Parallel()

	l := createLog(t)

	if lsn, err := l.Append(data); err != nil {
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

		lsn1, err := l.Append(data)
		if err != nil {
			t.Error(err)
		}

		_, err = l.Append(data)
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

		_, err := l.Append(data)
		if err != nil {
			t.Error(err)
		}

		_, err = l.Append(data)
		if err != nil {
			t.Error(err)
		}

		if err := l.Flush(10); err != nil {
			t.Error(err)
		}

		testPosition(t, l, 2, 2)
	})
}

func createLog(t *testing.T) *Log {
	t.Helper()

	f := mock.NewMockFile(path.Join(t.TempDir(), "logfile"))

	l := &Log{}
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
