package page

import (
	"bytes"
	"io"
	"path"
	"reflect"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/filetesting"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/storage/storagetesting"
	"github.com/liaradb/liaradb/util/testutil"
)

const (
	testHeaderSize = 2 + headerSize
)

func TestPage(t *testing.T) {
	t.Parallel()

	fsys := filetesting.NewMockFileSystem(t, nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	pid, tlid, rem, p := createWriter()

	if err := p.Write(f); err != nil {
		t.Fatal(err)
	}

	_, err := p.Iterate(io.NewSectionReader(f, 256, 256))
	if err != nil {
		t.Fatal(err)
	}

	testPage(t, p, pid, tlid, rem)
}

func TestPage_Append(t *testing.T) {
	t.Parallel()

	fsys := filetesting.NewMockFileSystem(t, nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	pid, tlid, rem, p := createWriter()

	rc, data, err := createRecord()
	if err != nil {
		t.Fatal(err)
	}

	if ok := p.Append(data); !ok {
		t.Fatal("should append record")
	}

	if ok := p.Append(data); !ok {
		t.Fatal("should append record")
	}

	if err := p.Write(f); err != nil {
		t.Fatal(err)
	}

	it, err := p.Iterate(io.NewSectionReader(f, 256, 256))
	if err != nil {
		t.Fatal(err)
	}

	testPage(t, p, pid, tlid, rem)

	count := 0
	for r, err := range it {
		count++
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(r, rc) {
			t.Error("data does not match")
		}
	}

	if count != 2 {
		t.Errorf("incorrect count: %v, expected: %v", count, 2)
	}
}

func TestPage_Iterate(t *testing.T) {
	t.Parallel()

	f, p := createReaderWriter(t)

	var count = record.NewLogSequenceNumber(3)
	records, _ := createRecords(count)

	for _, rc := range records {
		d, err := recordToBytes(rc)
		if err != nil {
			t.Fatal(err)
		}

		if ok := p.Append(d); !ok {
			t.Fatal("should append record")
		}
	}

	if err := p.Write(f); err != nil {
		t.Error(err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	var c record.LogSequenceNumber
	it, err := p.Iterate(f)
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

func TestPage_Reverse(t *testing.T) {
	t.Parallel()

	f, p := createReaderWriter(t)

	var count = record.NewLogSequenceNumber(3)
	records, _ := createRecords(count)

	for _, rc := range records {
		d, err := recordToBytes(rc)
		if err != nil {
			t.Fatal(err)
		}

		if ok := p.Append(d); !ok {
			t.Fatal("should append record")
		}
	}

	if err := p.Write(f); err != nil {
		t.Error(err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	var c record.LogSequenceNumber
	it, err := p.Reverse(f)
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

func TestNode_Append(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Append)
}

func testNode_Append(t *testing.T) {
	const (
		size int16 = 256
		s0         = size - itemSize - testHeaderSize
		s1         = s0 - itemSize - 16
		s2         = s1 - itemSize - 16
	)

	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	p := NewFromSlice(b.Raw())
	v0 := []byte{1, 2, 3, 4, 5}
	v1 := []byte{6, 7, 8, 9, 10}

	if s := p.Space(); s != s0 {
		t.Errorf("incorrect space: %v, expected: %v", s, s0)
	}

	if ok := p.Append(v0); !ok {
		t.Fatal("should append record")
	}

	// if s := p.Space(); s != s1 {
	// 	t.Errorf("incorrect space: %v, expected: %v", s, s1)
	// }

	if ok := p.Append(v1); !ok {
		t.Fatal("should append record")
	}

	// if s := p.Space(); s != s2 {
	// 	t.Errorf("incorrect space: %v, expected: %v", s, s2)
	// }

	// if _, err := raw.NewBufferFromSlice(b0).Seek(0, io.SeekStart); err != nil {
	// 	t.Error(err)
	// }

	result := make([][]byte, 0)
	for i := range p.Items() {
		result = append(result, i)
	}
	// // r0 := make([]byte, 5)
	// // if _, err := raw.NewBufferFromSlice(b0).Read(r0); err != nil {
	// // 	t.Error(err)
	// // }

	r0 := result[0]
	if !slices.Equal(r0, v0) {
		t.Errorf("incorrect result: %v, expected: %v", r0, v0)
	}

	// if _, err := raw.NewBufferFromSlice(b1).Seek(0, io.SeekStart); err != nil {
	// 	t.Error(err)
	// }

	// r1 := make([]byte, 5)
	// if _, err := raw.NewBufferFromSlice(b1).Read(r1); err != nil {
	// 	t.Error(err)
	// }

	r1 := result[1]
	if !slices.Equal(r1, v1) {
		t.Errorf("incorrect result: %v, expected: %v", r1, v1)
	}
}

func TestNode_Space(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Space)
}

func testNode_Space(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 16+itemSize+testHeaderSize)
	b := createBuffer(t, s)
	defer b.Release()

	p := NewFromSlice(b.Raw())

	if s := p.Space(); s != 16 {
		t.Fatalf("incorrect space: %v, expected: %v", s, 16)
	}

	if ok := p.Append(make([]byte, 16)); !ok {
		t.Fatal("should append record")
	}

	if s := p.Space(); s != 0 {
		t.Fatalf("incorrect space: %v, expected: %v", s, 0)
	}

	if ok := p.Append(make([]byte, 16)); ok {
		t.Fatal("should not get a buffer")
	}
}

func TestNode_Child(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Child)
}

func testNode_Child(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	b.Release()

	p := NewFromSlice(b.Raw())
	values := [][]byte{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10}}

	if ok := p.Append(values[0]); !ok {
		t.Fatal("should append record")
	}

	if ok := p.Append(values[1]); !ok {
		t.Fatal("should append record")
	}

	result := make([][]byte, 0, 2)
	for i := range 2 {
		c, ok := p.Child(int16(i))
		if !ok {
			t.Fatal("should get a buffer")
		}

		v := make([]byte, 5)
		if _, err := buffer.NewFromSlice(c).Read(v); err != nil {
			t.Fatal(err)
		}

		result = append(result, v)
	}

	if !slices.EqualFunc(result, values, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, values)
	}
}

func TestNode_Items(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Items)
}

func testNode_Items(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	p := NewFromSlice(b.Raw())
	values := [][]byte{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
		{9, 10}}

	for _, v := range values {
		if ok := p.Append(v); !ok {
			t.Fatal("should append record")
		}
	}

	result := make([][]byte, 0, len(values))
	for c := range p.Items() {
		v := make([]byte, 2)
		if _, err := buffer.NewFromSlice(c).Read(v); err != nil {
			t.Fatal(err)
		}

		result = append(result, v)
	}

	if !slices.EqualFunc(result, values, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, values)
	}
}

func TestNode_ChildrenRange(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_ChildrenRange)
}

func testNode_ChildrenRange(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	p := NewFromSlice(b.Raw())
	values := [][]byte{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
		{9, 10}}

	for _, v := range values {
		if ok := p.Append(v); !ok {
			t.Fatal("should append record")
		}
	}

	result := make([][]byte, 0, len(values))
	for c := range p.ChildrenRange(1, 4) {
		v := make([]byte, 2)
		if _, err := buffer.NewFromSlice(c).Read(v); err != nil {
			t.Fatal(err)
		}

		result = append(result, v)
	}

	want := [][]byte{
		{3, 4},
		{5, 6},
		{7, 8}}

	if !slices.EqualFunc(result, want, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func createWriter() (action.PageID, action.TimeLineID, record.Length, *Page) {
	pid := action.PageID(1)
	tlid := action.TimeLineID(2)
	rem := record.NewLength(3)

	p := New(256)
	p.Init(pid, tlid, rem)

	return pid, tlid, rem, p
}

func createRecord() (*record.Record, []byte, error) {
	lsn := record.NewLogSequenceNumber(1)
	tid := value.NewTenantID()
	txid := record.NewTransactionID(2)
	now := time.UnixMicro(1234567890)
	data := []byte("abcde")
	reverse := []byte("fghij")

	rc := record.New(lsn, tid, txid, now, record.ActionInsert, data, reverse)
	data, err := recordToBytes(rc)
	return rc, data, err
}

func recordToBytes(rc *record.Record) ([]byte, error) {
	recordBuf := bytes.NewBuffer(nil)
	if err := rc.Write(recordBuf); err != nil {
		return nil, err
	}

	return recordBuf.Bytes(), nil
}

func testPage(
	t *testing.T,
	p *Page,
	pid action.PageID,
	tlid action.TimeLineID,
	rem record.Length,
) {
	t.Helper()
	testutil.Getter(t, p.ID, pid, "ID")
	testutil.Getter(t, p.TimeLineID, tlid, "TimeLineID")
	testutil.Getter(t, p.LengthRemaining, rem, "LengthRemaining")
}

func createReaderWriter(t *testing.T) (file.File, *Page) {
	t.Helper()

	fsys := filetesting.NewMockFileSystem(t, nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))

	return f, New(256)
}

func createRecords(count record.LogSequenceNumber) ([]*record.Record, record.LogSequenceNumber) {
	tid := value.NewTenantID()
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	records := make([]*record.Record, 0, count.Value())
	for i := range count.Value() {
		records = append(records, record.New(record.NewLogSequenceNumber(i), tid, record.NewTransactionID(2), time.UnixMicro(1234567890), record.ActionInsert, data, reverse))
	}
	return records, count.Decrement()
}

func createBuffer(t *testing.T, s *storage.Storage) *storage.Buffer {
	b, err := s.Request(t.Context(), link.BlockID{})
	if err != nil {
		t.Fatal(err)
	}

	return b
}
