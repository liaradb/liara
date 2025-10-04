package page

import (
	"bytes"
	"io"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/liaradb/liaradb/filetesting"
	"github.com/liaradb/liaradb/log/record"
)

func TestWriter(t *testing.T) {
	t.Parallel()

	fsys := filetesting.NewMockFileSystem(t, nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	pid, tlid, rem, wr := createWriter()

	if err := wr.Write(f); err != nil {
		t.Fatal(err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	pr := NewReader(256)
	_, err := pr.Iterate(f)
	if err != nil {
		t.Fatal(err)
	}

	testHeader(t, pr.Header(), pid, tlid, rem)
}

func TestWriter_Append(t *testing.T) {
	t.Parallel()

	fsys := filetesting.NewMockFileSystem(t, nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	pid, tlid, rem, pw := createWriter()

	rc, data, err := createRecord()
	if err != nil {
		t.Fatal(err)
	}

	rb := record.NewBoundary(data)

	if err := pw.Append(rb, data); err != nil {
		t.Fatal(err)
	}

	if err := pw.Append(rb, data); err != nil {
		t.Fatal(err)
	}

	if err := pw.Write(f); err != nil {
		t.Fatal(err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	pr := NewReader(256)
	it, err := pr.Iterate(f)
	if err != nil {
		t.Fatal(err)
	}

	testHeader(t, pr.Header(), pid, tlid, rem)

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

func createWriter() (PageID, TimeLineID, record.Length, *Writer) {
	pid := PageID(1)
	tlid := TimeLineID(2)
	rem := record.Length(3)

	pw := NewWriter(256)
	pw.Init(pid, tlid, rem)

	return pid, tlid, rem, pw
}

func createRecord() (*record.Record, []byte, error) {
	lsn := record.LogSequenceNumber(1)
	tid := record.TransactionID(2)
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
