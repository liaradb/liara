package record

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/liaradb/liaradb/domain/value"
)

func TestSpan_Write(t *testing.T) {
	t.Parallel()

	lsn := NewLogSequenceNumber(1)
	tid := value.NewTenantID()
	txid := NewTransactionID(2)
	now := NewTime(time.UnixMicro(1234567890))
	action := ActionInsert
	collection := CollectionEvent
	data := []byte("abcdef")
	reverse := []byte("fghij")

	rc := New(lsn, tid, txid, now, action, collection, data, reverse)

	size := float64(rc.Size()) / 2
	a, b := int(math.Floor(size)), int(math.Ceil(size))

	s := NewSpan(
		NewFragment(make([]byte, a)),
		NewFragment(make([]byte, b)),
	)

	if err := rc.Write(s); err != nil {
		t.Fatal(err)
	}

	if err := s.SeekStart(); err != nil {
		t.Fatal(err)
	}

	rc2 := &Record{}
	if err := rc2.Read(s); err != nil {
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
