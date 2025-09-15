package log

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/cardboardrobots/assert"
)

func TestLogPage(t *testing.T) {
	lpid := LogPageID(1)
	tlid := TimeLineID(2)

	lp := newLogPage(256)
	lp.init(lpid, tlid)

	r, w := assert.NewReaderWriter()

	if err := lp.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	lp2 := &LogPage{}
	if err := lp2.Read(r); err != nil {
		t.Fatal(err)
	}

	assert.Getter(t, lp2.ID, lpid, "ID")
	assert.Getter(t, lp2.TimeLineID, tlid, "TimeLineID")
}

func TestLogPage_Append(t *testing.T) {
	r, w := assert.NewReaderWriter()

	lpid := LogPageID(1)
	tlid := TimeLineID(2)
	lp := newLogPage(256)
	lp.init(lpid, tlid)

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

	lp2 := newLogPage(256)
	if err := lp2.Read(r); err != nil {
		t.Fatal(err)
	}

	assert.Getter(t, lp2.ID, lpid, "ID")
	assert.Getter(t, lp2.TimeLineID, tlid, "TimeLineID")
	// TODO: This is not using the public API
	assert.EqualsArray(t, lp.data, lp2.data, "data")

	count := 0
	for r, err := range lp2.Records() {
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
