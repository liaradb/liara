package log

import (
	"testing"
	"time"

	"github.com/liaradb/liaradb/log/record"
)

func TestLog_Default(t *testing.T) {
	t.Parallel()

	wr := createWriter(t)

	testPosition(t, wr, 0, 0)
}

func TestLog_Append(t *testing.T) {
	t.Parallel()

	wr := createWriter(t)
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}
	var rec = record.New(1, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse)

	if lsn, err := wr.Append(rec); err != nil {
		t.Error(err)
	} else if lsn != 1 {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, wr, 0, 1)
}

func TestLog_Flush(t *testing.T) {
	t.Parallel()

	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}
	var rec = record.New(1, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse)

	t.Run("should flush", func(t *testing.T) {
		t.Parallel()

		wr := createWriter(t)

		lsn1, err := wr.Append(rec)
		if err != nil {
			t.Error(err)
		}

		_, err = wr.Append(rec)
		if err != nil {
			t.Error(err)
		}

		if err := wr.Flush(lsn1); err != nil {
			t.Error(err)
		}

		testPosition(t, wr, 1, 2)
	})

	t.Run("should not flush beyond HighWater", func(t *testing.T) {
		t.Parallel()

		wr := createWriter(t)

		_, err := wr.Append(rec)
		if err != nil {
			t.Error(err)
		}

		_, err = wr.Append(rec)
		if err != nil {
			t.Error(err)
		}

		if err := wr.Flush(10); err != nil {
			t.Error(err)
		}

		testPosition(t, wr, 2, 2)
	})

	t.Run("should write to multiple pages", func(t *testing.T) {
		t.Parallel()

		wr := createWriter(t)

		count := 10

		for range count - 1 {
			_, err := wr.Append(rec)
			if err != nil {
				t.Fatal(err)
			}
		}

		lsn2, err := wr.Append(rec)
		if err != nil {
			t.Fatal(err)
		}

		if err := wr.Flush(lsn2); err != nil {
			t.Fatal(err)
		}

		if p := wr.PageID(); p != 2 {
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

		wr := createWriter(t)

		lsn1, err := wr.Append(rec)
		if err != nil {
			t.Error(err)
		}

		if err := wr.Flush(lsn1); err != nil {
			t.Error(err)
		}

		lsn2, err := wr.Append(rec)
		if err != nil {
			t.Error(err)
		}

		if err := wr.Flush(lsn2); err != nil {
			t.Error(err)
		}

		testPosition(t, wr, 2, 2)
	})
}

func createWriter(t *testing.T) *Log {
	t.Helper()

	fsys, dir := createFiles(t)
	l := NewLog(256, 3, fsys, dir)
	if err := l.Open(); err != nil {
		t.Fatal(err)
	}

	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	return l
}

func testPosition(t *testing.T, l *Log, lw, hw record.LogSequenceNumber) {
	if h := l.HighWater(); h != hw {
		t.Errorf("incorrect high water: %v, expected: %v", h, hw)
	}

	if l := l.LowWater(); l != lw {
		t.Errorf("incorrect low water: %v, expected: %v", l, lw)
	}
}
