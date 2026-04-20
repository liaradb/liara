package segment

import (
	"path"
	"reflect"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/util/testing/filetesting"
)

func TestReader_Iterate(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testReader_Iterate)
}

func testReader_Iterate(t *testing.T) {
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
	synctest.Test(t, testReader_Reverse)
}

func testReader_Reverse(t *testing.T) {
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

	fsys := filetesting.New(nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))

	sw := NewWriter(270, 4)
	sw.SeekTail(0, f)
	return f, NewReader(270), sw
}

func createRecords(count record.LogSequenceNumber) ([]*record.Record, record.LogSequenceNumber) {
	tid := value.NewTenantID()
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	records := make([]*record.Record, 0, count.Value())
	for i := range count.Value() {
		records = append(records, record.New(
			record.NewLogSequenceNumber(i),
			tid,
			record.NewTransactionID(2),
			record.NewTime(time.UnixMicro(1234567890)),
			record.ActionInsert,
			record.CollectionEvent,
			data,
			reverse))
	}
	return records, count.Decrement()
}
