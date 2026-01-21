package page

import (
	"bytes"
	"io"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/filetesting"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/node"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/util/testutil"
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

	// TODO: This is private
	testPage(t, p.page, pid, tlid, rem)
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

	// TODO: This is private
	testPage(t, p.page, pid, tlid, rem)

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

func createWriter() (action.PageID, action.TimeLineID, record.Length, *Page) {
	pid := action.PageID(1)
	tlid := action.TimeLineID(2)
	rem := record.NewLength(3)

	p := New(256, node.New(256))
	p.Init(pid, tlid, rem)

	return pid, tlid, rem, p
}

func createRecord() (*record.Record, []byte, error) {
	lsn := record.NewLogSequenceNumber(1)
	tid := record.NewTransactionID(2)
	now := time.UnixMicro(1234567890)
	data := []byte("abcde")
	reverse := []byte("fghij")

	rc := record.New(lsn, tid, now, record.ActionInsert, data, reverse)
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
	p *node.Node,
	pid action.PageID,
	tlid action.TimeLineID,
	rem record.Length,
) {
	t.Helper()
	testutil.Getter(t, p.ID, pid, "ID")
	testutil.Getter(t, p.TimeLineID, tlid, "TimeLineID")
	testutil.Getter(t, p.LengthRemaining, rem, "LengthRemaining")
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

func createReaderWriter(t *testing.T) (file.File, *Page) {
	t.Helper()

	fsys := filetesting.NewMockFileSystem(t, nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	// fs := &file.FileSystem{}
	// f, _ := fs.Open(path.Join(t.TempDir(), "logfile"))

	return f, New(256, node.New(256))
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
