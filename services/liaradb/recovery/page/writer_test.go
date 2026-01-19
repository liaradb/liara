package page

import (
	"bytes"
	"io"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/liaradb/liaradb/file/filetesting"
	"github.com/liaradb/liaradb/recovery/mempage"
	"github.com/liaradb/liaradb/recovery/record"
)

func TestWriter(t *testing.T) {
	t.Parallel()

	fsys := filetesting.NewMockFileSystem(t, nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	pid, tlid, rem, wr := createWriter()

	if err := wr.Write(f); err != nil {
		t.Fatal(err)
	}

	pr := NewReader(mempage.NewWithHeader(256, &Header{}))
	_, err := pr.Iterate(io.NewSectionReader(f, 256, 256))
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

	if err := pw.Append(data); err != nil {
		t.Fatal(err)
	}

	if err := pw.Append(data); err != nil {
		t.Fatal(err)
	}

	if err := pw.Write(f); err != nil {
		t.Fatal(err)
	}

	pr := NewReader(mempage.NewWithHeader(256, &Header{}))
	it, err := pr.Iterate(io.NewSectionReader(f, 256, 256))
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
	rem := record.NewLength(3)

	pw := NewWriter(256, mempage.NewWithHeader(256, &Header{}))
	pw.Init(pid, tlid, rem)

	return pid, tlid, rem, pw
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
