package page

import (
	"bytes"
	"io"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/liaradb/liaradb/file/filetesting"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/node"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/util/testutil"
)

func TestWriter(t *testing.T) {
	t.Parallel()

	fsys := filetesting.NewMockFileSystem(t, nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	pid, tlid, rem, pr, wr := createWriter()

	if err := wr.Write(f); err != nil {
		t.Fatal(err)
	}

	_, err := pr.Iterate(io.NewSectionReader(f, 256, 256))
	if err != nil {
		t.Fatal(err)
	}

	// TODO: This is private
	testPage(t, pr.page, pid, tlid, rem)
}

func TestWriter_Append(t *testing.T) {
	t.Parallel()

	fsys := filetesting.NewMockFileSystem(t, nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	pid, tlid, rem, pr, pw := createWriter()

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

	it, err := pr.Iterate(io.NewSectionReader(f, 256, 256))
	if err != nil {
		t.Fatal(err)
	}

	// TODO: This is private
	testPage(t, pr.page, pid, tlid, rem)

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

func createWriter() (action.PageID, action.TimeLineID, record.Length, *Reader, *Writer) {
	pid := action.PageID(1)
	tlid := action.TimeLineID(2)
	rem := record.NewLength(3)

	n := node.New(make([]byte, 256))
	pr := NewReader(n)
	pw := NewWriter(256, n)
	pw.Init(pid, tlid, rem)

	return pid, tlid, rem, pr, pw
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
	p Page,
	pid action.PageID,
	tlid action.TimeLineID,
	rem record.Length,
) {
	t.Helper()
	testutil.Getter(t, p.ID, pid, "ID")
	testutil.Getter(t, p.TimeLineID, tlid, "TimeLineID")
	testutil.Getter(t, p.LengthRemaining, rem, "LengthRemaining")
}
