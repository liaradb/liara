package log

import (
	"testing"
	"testing/fstest"
	"time"

	"github.com/liaradb/liaradb/file/mock"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

func TestWriter_Default(t *testing.T) {
	t.Parallel()

	l := createLogWriter(t)

	testPosition(t, l, 0, 0)
}

func TestWriter_Append(t *testing.T) {
	t.Parallel()

	lw := createLogWriter(t)
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}
	var rec = record.New(1, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse)

	if lsn, err := lw.AppendToSegment(rec); err != nil {
		t.Error(err)
	} else if lsn != 1 {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, lw, 0, 1)
}

func TestWriter_Flush(t *testing.T) {
	t.Parallel()

	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}
	var rec = record.New(1, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse)

	t.Run("should flush", func(t *testing.T) {
		t.Parallel()

		lw := createLogWriter(t)

		lsn1, err := lw.AppendToSegment(rec)
		if err != nil {
			t.Error(err)
		}

		_, err = lw.AppendToSegment(rec)
		if err != nil {
			t.Error(err)
		}

		if err := lw.Flush(lsn1); err != nil {
			t.Error(err)
		}

		testPosition(t, lw, 1, 2)
	})

	t.Run("should not flush beyond HighWater", func(t *testing.T) {
		t.Parallel()

		lw := createLogWriter(t)

		_, err := lw.AppendToSegment(rec)
		if err != nil {
			t.Error(err)
		}

		_, err = lw.AppendToSegment(rec)
		if err != nil {
			t.Error(err)
		}

		if err := lw.Flush(10); err != nil {
			t.Error(err)
		}

		testPosition(t, lw, 2, 2)
	})

	t.Run("should write to multiple pages", func(t *testing.T) {
		t.Parallel()

		lw := createLogWriter(t)

		count := 10

		for range count - 1 {
			_, err := lw.AppendToSegment(rec)
			if err != nil {
				t.Fatal(err)
			}
		}

		lsn2, err := lw.AppendToSegment(rec)
		if err != nil {
			t.Fatal(err)
		}

		if err := lw.Flush(lsn2); err != nil {
			t.Fatal(err)
		}

		if p := lw.PageID(); p != 2 {
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

		lw := createLogWriter(t)

		lsn1, err := lw.AppendToSegment(rec)
		if err != nil {
			t.Error(err)
		}

		if err := lw.Flush(lsn1); err != nil {
			t.Error(err)
		}

		lsn2, err := lw.AppendToSegment(rec)
		if err != nil {
			t.Error(err)
		}

		if err := lw.Flush(lsn2); err != nil {
			t.Error(err)
		}

		testPosition(t, lw, 2, 2)
	})
}

func createLogWriter(t *testing.T) *Writer {
	t.Helper()

	fsys := mock.NewFileSystem(fstest.MapFS{
		"log/": &fstest.MapFile{},
	})
	// fsys := &file.FileSystem{}

	sl := segment.NewList(fsys, "log")

	lw := NewWriter(256, 3, sl)
	if err := lw.Start(); err != nil {
		t.Fatal(err)
	}
	if err := lw.Initialize(); err != nil {
		t.Fatal(err)
	}
	return lw
}

func testPosition(t *testing.T, sw *Writer, lw, hw record.LogSequenceNumber) {
	if h := sw.HighWater(); h != hw {
		t.Errorf("incorrect high water: %v, expected: %v", h, hw)
	}

	if l := sw.LowWater(); l != lw {
		t.Errorf("incorrect low water: %v, expected: %v", l, lw)
	}
}
