package segment

import (
	"io"
	"path"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/util/testing/filetesting"
)

func TestWriter_Append(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testWriter_Append)
}

func testWriter_Append(t *testing.T) {
	sw := createWriter(t)
	tid := value.NewTenantID()
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}
	var rec = record.New(
		record.NewLogSequenceNumber(1),
		tid,
		record.NewTransactionID(2),
		record.NewTime(time.UnixMicro(1234567890)),
		record.ActionInsert,
		record.CollectionEvent,
		data,
		reverse)

	if _, err := sw.Append(rec); err != nil {
		t.Error(err)
	}
}

func TestWriter_SeekTail(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testWriter_SeekTail)
}

func testWriter_SeekTail(t *testing.T) {
	f := createFile(t)
	sw := createWriterFromFile(t, f)
	tid := value.NewTenantID()
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}
	var rec = record.New(
		record.NewLogSequenceNumber(1),
		tid,
		record.NewTransactionID(2),
		record.NewTime(time.UnixMicro(1234567890)),
		record.ActionInsert,
		record.CollectionEvent,
		data,
		reverse)

	if _, err := sw.Append(rec); err != nil {
		t.Fatal(err)
	}

	if err := sw.Flush(); err != nil {
		t.Error(err)
	}

	// Seek start
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	sw2 := createWriterFromFile(t, f)

	if _, err := sw2.Append(rec); err != nil {
		t.Fatal(err)
	}

	if err := sw2.Flush(); err != nil {
		t.Error(err)
	}

	// Seek start
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	lr := NewReader(256)

	want := []*record.Record{rec, rec}
	result := make([]*record.Record, 0)
	for r, err := range lr.Iterate(f) {
		if err != nil {
			t.Fatal(err)
		}
		result = append(result, r)
	}

	if !slices.EqualFunc(result, want, (*record.Record).Compare) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func TestWriter_Flush(t *testing.T) {
	t.Parallel()

	tid := value.NewTenantID()
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}
	var rec = record.New(
		record.NewLogSequenceNumber(1),
		tid,
		record.NewTransactionID(2),
		record.NewTime(time.UnixMicro(1234567890)),
		record.ActionInsert,
		record.CollectionEvent,
		data,
		reverse)

	t.Run("should flush", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			sw := createWriter(t)

			if _, err := sw.Append(rec); err != nil {
				t.Error(err)
			}

			if _, err := sw.Append(rec); err != nil {
				t.Error(err)
			}

			if err := sw.Flush(); err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("should not flush beyond HighWater", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			sw := createWriter(t)

			if _, err := sw.Append(rec); err != nil {
				t.Error(err)
			}

			if _, err := sw.Append(rec); err != nil {
				t.Error(err)
			}

			if err := sw.Flush(); err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("should write to multiple pages", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			l := createWriter(t)

			count := 10

			for range count - 1 {
				if _, err := l.Append(rec); err != nil {
					t.Fatal(err)
				}
			}

			if err := l.Flush(); err != nil {
				t.Fatal(err)
			}

			if p := l.PageID(); p != 2 {
				t.Errorf("incorrect value: %v, expected: %v", p, 2)
			}
		})
	})

	t.Run("should return error if appending beyond maximum", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			l := createWriterSmall(t)

			if _, err := l.Append(rec); err != raw.ErrInsufficientSpace {
				t.Error("should return error")
			}
		})
	})

	t.Run("should write after flushing", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			l := createWriter(t)

			if _, err := l.Append(rec); err != nil {
				t.Error(err)
			}

			if err := l.Flush(); err != nil {
				t.Error(err)
			}

			if _, err := l.Append(rec); err != nil {
				t.Error(err)
			}

			if err := l.Flush(); err != nil {
				t.Error(err)
			}
		})
	})
}

func createWriter(t *testing.T) *Writer {
	t.Helper()

	f := filetesting.NewMockFile(path.Join(t.TempDir(), "logfile"), 0, time.Time{})
	f.Open()
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))

	sw := NewWriter(270, 3, 270)
	if err := sw.SeekTail(0, f); err != nil {
		t.Fatal(err)
	}

	return sw
}

func createWriterSmall(t *testing.T) *Writer {
	t.Helper()

	f := filetesting.NewMockFile(path.Join(t.TempDir(), "logfile"), 0, time.Time{})
	f.Open()

	sw := NewWriter(32, 1, 32)
	sw.SeekTail(0, f)
	return sw
}

func createFile(t *testing.T) filecache.File {
	t.Helper()

	f := filetesting.NewMockFile(path.Join(t.TempDir(), "logfile"), 0, time.Time{})
	f.Open()
	return f
}

func createWriterFromFile(t *testing.T, f filecache.File) *Writer {
	t.Helper()

	sw := NewWriter(256, 3, 256)

	i, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if err := sw.SeekTail(i.Size(), f); err != nil {
		t.Fatal(err)
	}

	return sw
}
