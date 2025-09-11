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

func GetterArray[T comparable](t *testing.T, a func() []T, b []T, name string) bool {
	t.Helper()

	value := a()
	equals := reflect.DeepEqual(value, b)
	if !equals {
		t.Errorf("%v is incorrect.  Expected: %v, Recieved: %v", name, b, value)
	}

	return equals
}
