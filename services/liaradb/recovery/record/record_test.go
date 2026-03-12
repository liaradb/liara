package record

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/util/testing/testutil"
)

func TestRecord(t *testing.T) {
	t.Parallel()

	lsn := NewLogSequenceNumber(1)
	tid := value.NewTenantID()
	txid := NewTransactionID(2)
	now := time.UnixMicro(1234567890)
	action := ActionInsert
	data := []byte("abcde")
	reverse := []byte("fghij")

	rc := New(lsn, tid, txid, now, action, data, reverse)

	testutil.Getter(t, rc.LogSequenceNumber, lsn, "LogSequenceNumber")
	testutil.Getter(t, rc.TenantID, tid, "TenantID")
	testutil.Getter(t, rc.TransactionID, txid, "TransactionID")
	testutil.Getter(t, rc.Time, now, "Time")
	testutil.Getter(t, rc.Action, action, "Action")
	testutil.GetterArray(t, rc.Data, data, "Data")
	testutil.GetterArray(t, rc.Reverse, reverse, "Reverse")
}

func TestRecord_Write(t *testing.T) {
	t.Parallel()

	lsn := NewLogSequenceNumber(1)
	tid := value.NewTenantID()
	txid := NewTransactionID(2)
	now := time.UnixMicro(1234567890)
	action := ActionInsert
	data := []byte("abcde")
	reverse := []byte("fghij")

	rc := New(lsn, tid, txid, now, action, data, reverse)

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

	testutil.Getter(t, rc2.LogSequenceNumber, lsn, "LogSequenceNumber")
	testutil.Getter(t, rc2.TenantID, tid, "TenantID")
	testutil.Getter(t, rc2.TransactionID, txid, "TransactionID")
	testutil.Getter(t, rc.Time, now, "Time")
	testutil.Getter(t, rc.Action, action, "Action")
	testutil.GetterArray(t, rc2.Data, data, "Data")
	testutil.GetterArray(t, rc2.Reverse, reverse, "Reverse")
}

func newReaderWriter() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
