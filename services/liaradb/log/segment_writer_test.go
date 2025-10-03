package log

import (
	"path"
	"testing"
	"time"

	"github.com/liaradb/liaradb/file/mock"
	"github.com/liaradb/liaradb/log/record"
)

func TestLogWriter_Default(t *testing.T) {
	t.Parallel()

	l := createLogWriter(t)

	testPosition(t, l, 0, 0)
}

func TestLogWriter_Append(t *testing.T) {
	t.Parallel()

	l := createLogWriter(t)
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}
	var rec = record.New(1, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse)

	if lsn, err := l.Append(rec); err != nil {
		t.Error(err)
	} else if lsn != 1 {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, l, 0, 1)
}

func TestLogWriter_Flush(t *testing.T) {
	t.Parallel()

	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}
	var rec = record.New(1, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse)

	t.Run("should flush", func(t *testing.T) {
		t.Parallel()

		l := createLogWriter(t)

		lsn1, err := l.Append(rec)
		if err != nil {
			t.Error(err)
		}

		_, err = l.Append(rec)
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

		_, err := l.Append(rec)
		if err != nil {
			t.Error(err)
		}

		_, err = l.Append(rec)
		if err != nil {
			t.Error(err)
		}

		if err := l.Flush(10); err != nil {
			t.Error(err)
		}

		testPosition(t, l, 2, 2)
	})

	t.Run("should write to multiple pages", func(t *testing.T) {
		t.Parallel()

		l := createLogWriter(t)

		count := 10

		for range count - 1 {
			_, err := l.Append(rec)
			if err != nil {
				t.Fatal(err)
			}
		}

		lsn2, err := l.Append(rec)
		if err != nil {
			t.Fatal(err)
		}

		if err := l.Flush(lsn2); err != nil {
			t.Fatal(err)
		}

		if p := l.PageID(); p != 2 {
			t.Errorf("incorrect value: %v, expected: %v", p, 2)
		}
	})

	t.Run("should return error if appending beyond maximum", func(t *testing.T) {
		t.Parallel()
		t.Skip()
		// TODO: Test this
	})

	t.Run("should write after flushing", func(t *testing.T) {
		t.Parallel()

		l := createLogWriter(t)

		lsn1, err := l.Append(rec)
		if err != nil {
			t.Error(err)
		}

		if err := l.Flush(lsn1); err != nil {
			t.Error(err)
		}

		lsn2, err := l.Append(rec)
		if err != nil {
			t.Error(err)
		}

		if err := l.Flush(lsn2); err != nil {
			t.Error(err)
		}

		testPosition(t, l, 2, 2)
	})
}

func createLogWriter(t *testing.T) *SegmentWriter {
	t.Helper()

	f := mock.NewMockFile(path.Join(t.TempDir(), "logfile"))
	f.Open()
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))

	sw := NewSegmentWriter(256, 3, f)
	_ = sw.Initialize()
	return sw
}

func testPosition(t *testing.T, sw *SegmentWriter, lw, hw record.LogSequenceNumber) {
	if h := sw.HighWater(); h != hw {
		t.Errorf("incorrect high water: %v, expected: %v", h, hw)
	}

	if l := sw.LowWater(); l != lw {
		t.Errorf("incorrect low water: %v, expected: %v", l, lw)
	}
}
