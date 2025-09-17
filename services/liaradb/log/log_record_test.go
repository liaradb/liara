package log

import (
	"io"
	"testing"
	"time"

	"github.com/cardboardrobots/assert"
)

func TestLogRecord(t *testing.T) {
	t.Parallel()

	lsn := LogSequenceNumber(1)
	tid := TransactionID(2)
	now := time.UnixMicro(1234567890)
	data := []byte("abcde")
	reverse := []byte("fghij")

	lr := newLogRecord(lsn, tid, now, data, reverse)

	assert.Getter(t, lr.LogSequenceNumber, lsn, "LogSequenceNumber")
	assert.Getter(t, lr.TransactionID, tid, "TransactionID")
	assert.Getter(t, lr.Time, now, "Time")
	assert.GetterArray(t, lr.Data, data, "Data")
	assert.GetterArray(t, lr.Reverse, reverse, "Reverse")
}

func TestLogRecord_Write(t *testing.T) {
	t.Parallel()

	lsn := LogSequenceNumber(1)
	tid := TransactionID(2)
	now := time.UnixMicro(1234567890)
	data := []byte("abcde")
	reverse := []byte("fghij")

	lr := newLogRecord(lsn, tid, now, data, reverse)

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
	assert.Getter(t, lr.Time, now, "Time")
	assert.GetterArray(t, lr2.Data, data, "Data")
	assert.GetterArray(t, lr2.Reverse, reverse, "Reverse")
}

func TestLogRecord_Time(t *testing.T) {
	t.Parallel()

	r, w := assert.NewReaderWriter()

	lr := LogRecord{
		time: time.UnixMicro(1234567890)}
	if err := lr.writeTime(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var lr2 LogRecord
	if err := lr2.readTime(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if !lr.time.Equal(lr2.time) {
		t.Errorf("incorrect value: %v, expected: %v", lr.time, lr2.time)
	}
}
