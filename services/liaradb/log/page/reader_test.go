package page

import (
	"io"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/filetesting"
	"github.com/liaradb/liaradb/log/record"
)

func TestReader_Iterate(t *testing.T) {
	t.Parallel()

	f, rd, sw := createReaderWriter(t)

	var count = record.NewLogSequenceNumber(3)
	records, _ := createRecords(count)

	for _, rc := range records {
		d, err := recordToBytes(rc)
		if err != nil {
			t.Fatal(err)
		}

		if err := sw.Append(record.NewBoundary(d), d); err != nil {
			t.Fatal(err)
		}
	}

	if err := sw.Flush(f); err != nil {
		t.Error(err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	var c record.LogSequenceNumber
	it, err := rd.Iterate(f)
	if err != nil {
		t.Fatal(err)
	}

	for rc, err := range it {
		c = c.Increment()
		if err != nil {
			t.Fatal(err)
		}

		rec := records[c.Value()-1]

		if !reflect.DeepEqual(rc, rec) {
			t.Fatalf("incorrect value:\n%#v, expected:\n%#v", rc, rec)
		}
	}

	if c != count {
		t.Errorf("incorrect count: %v, expected: %v", c, count)
	}
}

func TestReader_Reverse(t *testing.T) {
	f, rd, sw := createReaderWriter(t)

	var count = record.NewLogSequenceNumber(3)
	records, _ := createRecords(count)

	for _, rc := range records {
		d, err := recordToBytes(rc)
		if err != nil {
			t.Fatal(err)
		}

		if err := sw.Append(record.NewBoundary(d), d); err != nil {
			t.Fatal(err)
		}
	}

	if err := sw.Flush(f); err != nil {
		t.Error(err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	var c record.LogSequenceNumber
	it, err := rd.Reverse(f)
	if err != nil {
		t.Fatal(err)
	}

	for rc, err := range it {
		c = c.Increment()
		if err != nil {
			t.Fatal(err)
		}

		rec := records[count.Value()-c.Value()]

		if !reflect.DeepEqual(rc, rec) {
			t.Fatalf("incorrect value:\n%#v, expected:\n%#v", rc, rec)
		}
	}

	if c != count {
		t.Errorf("incorrect count: %v, expected: %v", c, count)
	}
}

func createReaderWriter(t *testing.T) (file.File, *Reader, *Writer) {
	t.Helper()

	fsys := filetesting.NewMockFileSystem(t, nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))

	sw := NewWriter(256)
	return f, NewReader(256), sw
}

func createRecords(count record.LogSequenceNumber) ([]*record.Record, record.LogSequenceNumber) {
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	records := make([]*record.Record, 0, count.Value())
	for i := range count.Value() {
		records = append(records, record.New(record.NewLogSequenceNumber(i), record.NewTransactionID(2), time.UnixMicro(1234567890), record.ActionInsert, data, reverse))
	}
	return records, count.Decrement()
}
