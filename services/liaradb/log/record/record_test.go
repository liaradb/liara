package record

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/cardboardrobots/assert"
)

func TestRecord(t *testing.T) {
	t.Parallel()

	lsn := LogSequenceNumber(1)
	tid := NewTransactionID(2)
	now := time.UnixMicro(1234567890)
	action := ActionInsert
	data := []byte("abcde")
	reverse := []byte("fghij")

	rc := New(lsn, tid, now, action, data, reverse)

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
	tid := NewTransactionID(2)
	now := time.UnixMicro(1234567890)
	action := ActionInsert
	data := []byte("abcde")
	reverse := []byte("fghij")

	rc := New(lsn, tid, now, action, data, reverse)

	r, w := newReaderWriter()

	if err := rc.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := rc.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
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

func newReaderWriter() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
