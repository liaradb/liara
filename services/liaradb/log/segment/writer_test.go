package segment

import (
	"path"
	"testing"
	"time"

	"github.com/liaradb/liaradb/file/mock"
	"github.com/liaradb/liaradb/log/record"
)

func TestWriter_Append(t *testing.T) {
	t.Parallel()

	sw := createWriter(t)
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}
	var rec = record.New(record.NewLogSequenceNumber(1), record.NewTransactionID(2), time.UnixMicro(1234567890), record.ActionInsert, data, reverse)

	if err := sw.Append(rec); err != nil {
		t.Error(err)
	}
}

func TestWriter_Flush(t *testing.T) {
	t.Parallel()

	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}
	var rec = record.New(record.NewLogSequenceNumber(1), record.NewTransactionID(2), time.UnixMicro(1234567890), record.ActionInsert, data, reverse)

	t.Run("should flush", func(t *testing.T) {
		t.Parallel()

		sw := createWriter(t)

		if err := sw.Append(rec); err != nil {
			t.Error(err)
		}

		if err := sw.Append(rec); err != nil {
			t.Error(err)
		}

		if err := sw.Flush(); err != nil {
			t.Error(err)
		}
	})

	t.Run("should not flush beyond HighWater", func(t *testing.T) {
		t.Parallel()

		sw := createWriter(t)

		if err := sw.Append(rec); err != nil {
			t.Error(err)
		}

		if err := sw.Append(rec); err != nil {
			t.Error(err)
		}

		if err := sw.Flush(); err != nil {
			t.Error(err)
		}
	})

	t.Run("should write to multiple pages", func(t *testing.T) {
		t.Parallel()

		l := createWriter(t)

		count := 10

		for range count - 1 {
			if err := l.Append(rec); err != nil {
				t.Fatal(err)
			}
		}

		if err := l.Append(rec); err != nil {
			t.Fatal(err)
		}

		if err := l.Flush(); err != nil {
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

		l := createWriter(t)

		if err := l.Append(rec); err != nil {
			t.Error(err)
		}

		if err := l.Flush(); err != nil {
			t.Error(err)
		}

		if err := l.Append(rec); err != nil {
			t.Error(err)
		}

		if err := l.Flush(); err != nil {
			t.Error(err)
		}
	})
}

func createWriter(t *testing.T) *Writer {
	t.Helper()

	f := mock.NewMockFile(path.Join(t.TempDir(), "logfile"))
	f.Open()
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))

	sw := NewWriter(256, 3)
	_ = sw.Initialize(f)
	return sw
}
