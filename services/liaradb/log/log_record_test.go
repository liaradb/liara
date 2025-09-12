package log

import (
	"reflect"
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
	GetterArray(t, lr.Data, data, "Data")
	GetterArray(t, lr.Reverse, reverse, "Reverse")
}

func TestLogRecord_Write(t *testing.T) {
	lsn := LogSequenceNumber(1)
	tid := TransactionID(2)
	data := []byte("abcde")
	reverse := []byte("fghij")

	lr := newLogRecord(lsn, tid, data, reverse)

	r, w := createReaderWriter()

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
	GetterArray(t, lr2.Data, data, "Data")
	GetterArray(t, lr2.Reverse, reverse, "Reverse")
}

func GetterArray[T comparable](t *testing.T, a func() []T, b []T, name string) bool {
	t.Helper()

	value := a()
	equals := reflect.DeepEqual(value, b)
	if !equals {
		t.Errorf("%v is incorrect.  Expected: %v, Recieved: %v", name, b, value)
	}

	return equals
}
