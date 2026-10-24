package record

import (
	"bufio"
	"bytes"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

func TestRecord(t *testing.T) {
	t.Parallel()

	lsn := logpage.NewLogSequenceNumber(1)
	txid := NewTransactionID(2)
	rid := link.NewRecordID(link.NewBlockID("", 1), 2)
	action := ActionInsert
	data := []byte("abcde")
	reverse := []byte("fghij")

	rc := New(txid, rid, action, data, reverse)
	rc.SetLogSequenceNumber(lsn)

	if i := rc.LogSequenceNumber(); i != lsn {
		t.Errorf("incorrect log sequence number: %v, expected: %v", i, lsn)
	}

	if i := rc.TransactionID(); i != txid {
		t.Errorf("incorrect transaction id: %v, expected: %v", i, txid)
	}

	if i := rc.RecordID(); i != rid {
		t.Errorf("incorrect record locator: %v, expected: %v", i, rid)
	}

	if i := rc.Action(); i != action {
		t.Errorf("incorrect action: %v, expected: %v", i, action)
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

	lsn := logpage.NewLogSequenceNumber(1)
	txid := NewTransactionID(2)
	rid := link.NewRecordID(link.NewBlockID("", 1), 2)
	action := ActionInsert
	data := []byte("abcde")
	reverse := []byte("fghij")

	rc := New(txid, rid, action, data, reverse)
	rc.SetLogSequenceNumber(lsn)

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

	if i := rc2.TransactionID(); i != txid {
		t.Errorf("incorrect transaction id: %v, expected: %v", i, txid)
	}

	if i := rc2.RecordID(); i != rid {
		t.Errorf("incorrect record locator: %v, expected: %v", i, rid)
	}

	if i := rc2.Action(); i != action {
		t.Errorf("incorrect action: %v, expected: %v", i, action)
	}

	if i := rc2.Data(); !slices.Equal(i, data) {
		t.Errorf("incorrect data: %v, expected: %v", i, data)
	}

	if i := rc2.Reverse(); !slices.Equal(i, reverse) {
		t.Errorf("incorrect reverse: %v, expected: %v", i, reverse)
	}

	if i := rc2.IsCheckpoint(); i != (action == ActionCheckpoint) {
		t.Errorf("incorrect is checkpoint: %v, expected: %v", i, action == ActionCheckpoint)
	}
}

func TestRecord_Compare(t *testing.T) {
	t.Parallel()

	txid := NewTransactionID(2)
	rid := link.NewRecordID(link.NewBlockID("", 1), 2)
	action := ActionInsert
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
			a:     New(txid, rid, action, data, reverse),
			b:     New(txid, rid, action, data, reverse),
			equal: true,
		},
		"should not equal different values": {
			a:     New(txid, rid, action, data, reverse),
			b:     New(NewTransactionID(3), rid, action, data, reverse),
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
