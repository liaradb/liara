package record

import (
	"bufio"
	"bytes"
	"slices"
	"testing"
	"time"

	"github.com/liaradb/liaradb/domain/value"
)

func TestRecord(t *testing.T) {
	t.Parallel()

	lsn := NewLogSequenceNumber(1)
	tid := value.NewTenantID()
	txid := NewTransactionID(2)
	now := NewTime(time.UnixMicro(1234567890))
	action := ActionInsert
	collection := CollectionEvent
	data := []byte("abcde")
	reverse := []byte("fghij")

	rc := New(lsn, tid, txid, now, action, collection, data, reverse)

	if i := rc.LogSequenceNumber(); i != lsn {
		t.Errorf("incorrect log sequence number: %v, expected: %v", i, lsn)
	}

	if i := rc.TenantID(); i != tid {
		t.Errorf("incorrect tenant id: %v, expected: %v", i, tid)
	}

	if i := rc.TransactionID(); i != txid {
		t.Errorf("incorrect transaction id: %v, expected: %v", i, txid)
	}

	if i := rc.Time(); i != now {
		t.Errorf("incorrect time: %v, expected: %v", i, now)
	}

	if i := rc.Action(); i != action {
		t.Errorf("incorrect action: %v, expected: %v", i, action)
	}

	if i := rc.Collection(); i != collection {
		t.Errorf("incorrect collection: %v, expected: %v", i, collection)
	}

	if i := rc.Data(); !slices.Equal(i, data) {
		t.Errorf("incorrect data: %v, expected: %v", i, data)
	}

	if i := rc.Reverse(); !slices.Equal(i, reverse) {
		t.Errorf("incorrect reverse: %v, expected: %v", i, reverse)
	}

	if i := rc.IsCheckpoint(); i != (action == ActionCheckpoint) {
		t.Errorf("incorrect is checkpoint: %v, expected: %v", i, action == ActionCheckpoint)
	}
}

func TestRecord_Write(t *testing.T) {
	t.Parallel()

	lsn := NewLogSequenceNumber(1)
	tid := value.NewTenantID()
	txid := NewTransactionID(2)
	now := NewTime(time.UnixMicro(1234567890))
	action := ActionInsert
	collection := CollectionEvent
	data := []byte("abcde")
	reverse := []byte("fghij")

	rc := New(lsn, tid, txid, now, action, collection, data, reverse)

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

	if i := rc2.LogSequenceNumber(); i != lsn {
		t.Errorf("incorrect log sequence number: %v, expected: %v", i, lsn)
	}

	if i := rc2.TenantID(); i != tid {
		t.Errorf("incorrect tenant id: %v, expected: %v", i, tid)
	}

	if i := rc2.TransactionID(); i != txid {
		t.Errorf("incorrect transaction id: %v, expected: %v", i, txid)
	}

	if i := rc.Time(); i != now {
		t.Errorf("incorrect time: %v, expected: %v", i, now)
	}

	if i := rc.Action(); i != action {
		t.Errorf("incorrect action: %v, expected: %v", i, action)
	}

	if i := rc.Collection(); i != collection {
		t.Errorf("incorrect collection: %v, expected: %v", i, collection)
	}

	if i := rc2.Data(); !slices.Equal(i, data) {
		t.Errorf("incorrect data: %v, expected: %v", i, data)
	}

	if i := rc2.Reverse(); !slices.Equal(i, reverse) {
		t.Errorf("incorrect reverse: %v, expected: %v", i, reverse)
	}

	if i := rc.IsCheckpoint(); i != (action == ActionCheckpoint) {
		t.Errorf("incorrect is checkpoint: %v, expected: %v", i, action == ActionCheckpoint)
	}
}

func TestRecord_Compare(t *testing.T) {
	t.Parallel()

	lsn := NewLogSequenceNumber(1)
	tid := value.NewTenantID()
	txid := NewTransactionID(2)
	now := NewTime(time.UnixMicro(1234567890))
	action := ActionInsert
	collection := CollectionEvent
	data := []byte("abcde")
	reverse := []byte("fghij")

	pointer := &Record{}

	for message, c := range map[string]struct {
		skip  bool
		a     *Record
		b     *Record
		equal bool
	}{
		"should equal zero": {
			a:     &Record{},
			b:     &Record{},
			equal: true,
		},
		"should equal pointer": {
			a:     pointer,
			b:     pointer,
			equal: true,
		},
		"should equal same values": {
			a:     New(lsn, tid, txid, now, action, collection, data, reverse),
			b:     New(lsn, tid, txid, now, action, collection, data, reverse),
			equal: true,
		},
		"should not equal different values": {
			a:     New(lsn, tid, txid, now, action, collection, data, reverse),
			b:     New(lsn, value.NewTenantID(), txid, now, action, collection, data, reverse),
			equal: false,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			if c.a.Compare(c.b) != c.equal {
				if c.equal {
					t.Error("should equal")
				} else {
					t.Error("should not equal")
				}
			}
		})
	}
}

func newReaderWriter() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
