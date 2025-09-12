package log

import (
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestLogRecord(t *testing.T) {
	lsn := LogSequenceNumber(1)
	tid := TransactionID(2)
	data := []byte("abcde")
	reverse := []byte("fghij")

	lr := newLogRecord(lsn, tid, data, reverse)

	assert.Getter(t, lr.LogSequenceNumber, lsn, "LogSequenceNumber")
	assert.Getter(t, lr.TransactionID, tid, "TransactionID")
	assert.GetterArray(t, lr.Data, data, "Data")
	assert.GetterArray(t, lr.Reverse, reverse, "Reverse")
}

func TestLogRecord_Write(t *testing.T) {
	lsn := LogSequenceNumber(1)
	tid := TransactionID(2)
	data := []byte("abcde")
	reverse := []byte("fghij")

	lr := newLogRecord(lsn, tid, data, reverse)

	r, w := assert.NewReaderWriter()

	if err := lr.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	lr2 := &LogRecord{}
	if err := lr2.Read(r); err != nil {
		t.Fatal(err)
	}

	assert.Getter(t, lr2.LogSequenceNumber, lsn, "LogSequenceNumber")
	assert.Getter(t, lr2.TransactionID, tid, "TransactionID")
	assert.GetterArray(t, lr2.Data, data, "Data")
	assert.GetterArray(t, lr2.Reverse, reverse, "Reverse")
}
