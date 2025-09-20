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
	"github.com/liaradb/liaradb/log/record"
)

func TestPageWriter(t *testing.T) {
	t.Parallel()

	f := mock.NewMockFile(path.Join(t.TempDir(), "logfile"))
	pid, tlid, rem, pw := createPage()

	if err := pw.Write(f); err != nil {
		t.Fatal(err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	lr := NewLogReader(256, f)
	ph, err := lr.Read()
	if err != nil {
		t.Fatal(err)
	}

	testPageHeader(t, ph, pid, tlid, rem)
}

func TestPageWriter_Append(t *testing.T) {
	t.Parallel()

	f := mock.NewMockFile(path.Join(t.TempDir(), "logfile"))
	pid, tlid, rem, pw := createPage()

	rc, data, err := createRecord()
	if err != nil {
		t.Fatal(err)
	}

	crc := record.NewCRC(data)

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

	lr := NewLogReader(256, f)
	ph, err := lr.Read()
	if err != nil {
		t.Fatal(err)
	}

	testPageHeader(t, ph, pid, tlid, rem)

	count := 0
	for r, err := range lr.Records() {
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

func createPage() (record.PageID, record.TimeLineID, record.RecordLength, *PageWriter) {
	pid := record.PageID(1)
	tlid := record.TimeLineID(2)
	rem := record.RecordLength(3)

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

	rc := record.NewRecord(lsn, tid, now, data, reverse)
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
	ph *record.PageHeader,
	pid record.PageID,
	tlid record.TimeLineID,
	rem record.RecordLength,
) {
	t.Helper()
	assert.Getter(t, ph.ID, pid, "ID")
	assert.Getter(t, ph.TimeLineID, tlid, "TimeLineID")
	assert.Getter(t, ph.LengthRemaining, rem, "LengthRemaining")
}
