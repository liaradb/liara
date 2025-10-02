package log

import (
	"bytes"
	"io"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/cardboardrobots/assert"
	"github.com/liaradb/liaradb/file/mock"
	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/log/record"
)

func TestPageWriter(t *testing.T) {
	t.Parallel()

	fsys := mock.NewFileSystem(nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	pid, tlid, rem, pw := createPage()

	if err := pw.Write(f); err != nil {
		t.Fatal(err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	sr := NewSegmentReader(256, 2, f)
	ph, err := sr.Read()
	if err != nil {
		t.Fatal(err)
	}

	testPageHeader(t, ph, pid, tlid, rem)
}

func TestPageWriter_Append(t *testing.T) {
	t.Parallel()

	fsys := mock.NewFileSystem(nil)
	f, _ := fsys.OpenFile(path.Join(t.TempDir(), "logfile"))
	pid, tlid, rem, pw := createPage()

	rc, data, err := createRecord()
	if err != nil {
		t.Fatal(err)
	}

	crc := page.NewCRC(data)

	if err := pw.append(crc, data); err != nil {
		t.Fatal(err)
	}

	if err := pw.append(crc, data); err != nil {
		t.Fatal(err)
	}

	if err := pw.Write(f); err != nil {
		t.Fatal(err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	sr := NewSegmentReader(256, 2, f)
	ph, err := sr.Read()
	if err != nil {
		t.Fatal(err)
	}

	testPageHeader(t, ph, pid, tlid, rem)

	count := 0
	for r, err := range sr.Records() {
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

func createPage() (page.PageID, page.TimeLineID, page.RecordLength, *PageWriter) {
	pid := page.PageID(1)
	tlid := page.TimeLineID(2)
	rem := page.RecordLength(3)

	pw := createEmptyPage()
	pw.init(pid, tlid, rem)

	return pid, tlid, rem, pw
}

func createEmptyPage() *PageWriter {
	return newPageWriter(256)
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

func testPageHeader(
	t *testing.T,
	ph *page.PageHeader,
	pid page.PageID,
	tlid page.TimeLineID,
	rem page.RecordLength,
) {
	t.Helper()
	assert.Getter(t, ph.ID, pid, "ID")
	assert.Getter(t, ph.TimeLineID, tlid, "TimeLineID")
	assert.Getter(t, ph.LengthRemaining, rem, "LengthRemaining")
}
