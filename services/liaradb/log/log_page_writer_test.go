package log

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/cardboardrobots/assert"
)

func TestLogPageWriter(t *testing.T) {
	r, w := assert.NewReaderWriter()
	lpid, tlid, rem, lp := createPage()

	if err := lp.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	lpr := newLogPageReader(256)
	p, err := lpr.Read(r)
	if err != nil {
		t.Fatal(err)
	}

	testLogPageHeader(t, p, lpid, tlid, rem)
}

func TestLogPageWriter_Append(t *testing.T) {
	r, w := assert.NewReaderWriter()
	lpid, tlid, rem, lp := createPage()

	lr, data, err := createRecord()
	if err != nil {
		t.Fatal(err)
	}

	crc := NewCRC(data)

	if err := lp.append(crc, data); err != nil {
		t.Fatal(err)
	}

	if err := lp.append(crc, data); err != nil {
		t.Fatal(err)
	}

	if err := lp.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	lpr := newLogPageReader(256)
	p, err := lpr.Read(r)
	if err != nil {
		t.Fatal(err)
	}

	testLogPageHeader(t, p, lpid, tlid, rem)

	count := 0
	for r, err := range lpr.Records() {
		count++
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(r, lr) {
			t.Error("data does not match")
		}
	}

	if count != 2 {
		t.Errorf("incorrect count: %v, expected: %v", count, 2)
	}
}

func createPage() (LogPageID, TimeLineID, LogRecordLength, *LogPageWriter) {
	lpid := LogPageID(1)
	tlid := TimeLineID(2)
	rem := LogRecordLength(3)

	lp := createEmptyPage()
	lp.init(lpid, tlid, rem)

	return lpid, tlid, rem, lp
}

func createEmptyPage() *LogPageWriter {
	return newLogPageWriter(256)
}

func createRecord() (*LogRecord, []byte, error) {
	lsn := LogSequenceNumber(1)
	tid := TransactionID(2)
	now := time.UnixMicro(1234567890)
	data := []byte("abcde")
	reverse := []byte("fghij")

	lr := newLogRecord(lsn, tid, now, data, reverse)
	data, err := recordToBytes(lr)
	return lr, data, err
}

func recordToBytes(lr *LogRecord) ([]byte, error) {
	recordBuf := bytes.NewBuffer(nil)
	if err := lr.Write(recordBuf); err != nil {
		return nil, err
	}

	return recordBuf.Bytes(), nil
}
