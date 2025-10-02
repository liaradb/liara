package record

import (
	"io"
	"testing"
	"time"

	"github.com/cardboardrobots/assert"
)

func TestRecord(t *testing.T) {
	t.Parallel()

	lsn := LogSequenceNumber(1)
	tid := TransactionID(2)
	now := time.UnixMicro(1234567890)
	action := ActionInsert
	data := []byte("abcde")
	reverse := []byte("fghij")

	rc := NewRecord(lsn, tid, now, action, data, reverse)

	assert.Getter(t, rc.LogSequenceNumber, lsn, "LogSequenceNumber")
	assert.Getter(t, rc.TransactionID, tid, "TransactionID")
	assert.Getter(t, rc.Time, now, "Time")
	assert.Getter(t, rc.Action, action, "Action")
	assert.GetterArray(t, rc.Data, data, "Data")
	assert.GetterArray(t, rc.Reverse, reverse, "Reverse")
}

func TestRecord_Write(t *testing.T) {
	t.Parallel()

	lsn := LogSequenceNumber(1)
	tid := TransactionID(2)
	now := time.UnixMicro(1234567890)
	action := ActionInsert
	data := []byte("abcde")
	reverse := []byte("fghij")

	rc := NewRecord(lsn, tid, now, action, data, reverse)

	r, w := assert.NewReaderWriter()

	if err := rc.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Size() - w.Available()
	if l := rc.Size(); l != size {
		t.Errorf("incorrect length: %v, expected: %v", l, size)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	rc2 := &Record{}
	if err := rc2.Read(r); err != nil {
		t.Fatal(err)
	}

	assert.Getter(t, rc2.LogSequenceNumber, lsn, "LogSequenceNumber")
	assert.Getter(t, rc2.TransactionID, tid, "TransactionID")
	assert.Getter(t, rc.Time, now, "Time")
	assert.Getter(t, rc.Action, action, "Action")
	assert.GetterArray(t, rc2.Data, data, "Data")
	assert.GetterArray(t, rc2.Reverse, reverse, "Reverse")
}

func TestRecord_Time(t *testing.T) {
	t.Parallel()

	r, w := assert.NewReaderWriter()

	rc := Record{
		time: time.UnixMicro(1234567890)}
	if err := rc.writeTime(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var rc2 Record
	if err := rc2.readTime(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if !rc.time.Equal(rc2.time) {
		t.Errorf("incorrect value: %v, expected: %v", rc.time, rc2.time)
	}
}
