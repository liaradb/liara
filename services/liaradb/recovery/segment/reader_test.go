package segment

import (
	"path"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/filetesting"
	"github.com/liaradb/liaradb/recovery/node"
	"github.com/liaradb/liaradb/recovery/record"
)

func TestReader_Iterate(t *testing.T) {
	t.Parallel()

	f, lr, sw := createReaderWriter(t)

	var count = record.NewLogSequenceNumber(10)
	records, _ := createRecords(count)

	for _, rc := range records {
		if err := sw.Append(rc); err != nil {
			t.Fatal(err)
		}
	}

	if err := sw.Flush(); err != nil {
		t.Error(err)
	}

	result := make([]*record.Record, 0)
	for rc, err := range lr.Iterate(f) {
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, rc)
	}

	if !slices.EqualFunc(result, records, func(a, b *record.Record) bool {
		return reflect.DeepEqual(a, b)
	}) {
		t.Fatalf("incorrect result:\n%v\nexpected:\n%v", result, records)
	}
}

func TestReader_Reverse(t *testing.T) {
	t.Parallel()

	f, sr, sw := createReaderWriter(t)

	var count = record.NewLogSequenceNumber(10)
	records, _ := createRecords(count)

	for _, rc := range records {
		if err := sw.Append(rc); err != nil {
			t.Fatal(err)
		}
	}

	if err := sw.Flush(); err != nil {
		t.Error(err)
	}

	stat, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	result := make([]*record.Record, 0)
	for rc, err := range sr.Reverse(stat.Size(), f) {
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, rc)
	}

	slices.Reverse(records)
	if !slices.EqualFunc(result, records, func(a, b *record.Record) bool {
		return reflect.DeepEqual(a, b)
	}) {
		t.Fatalf("incorrect result:\n%v\nexpected:\n%v", result, records)
	}
}

func createReaderWriter(t *testing.T) (file.File, *Reader, *Writer) {
	t.Helper()

	fsys := filetesting.NewMockFileSystem(t, nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))

	n := node.New(make([]byte, 256))
	sw := NewWriter(256, 4, n)
	sw.Initialize(f)
	return f, NewReader(256, n), sw
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
