package pagequeue

import (
	"slices"
	"testing"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/recovery/record"
)

func TestPageQueue(t *testing.T) {
	t.Parallel()

	pq := New(256)

	lsn := record.NewLogSequenceNumber(1)
	tid := value.NewTenantID()
	txid := record.NewTransactionID(2)
	now := record.NewTime(time.UnixMicro(1234567890))
	action := record.ActionInsert
	collection := record.CollectionEvent
	data := []byte("abcdef")
	reverse := []byte("fghij")

	rc := record.New(lsn, tid, txid, now, action, collection, data, reverse)

	if err := pq.Append(rc); err != nil {
		t.Fatal(err)
	}

	if c := pq.Count(); c != 1 {
		t.Fatalf("incorrect count: %v, expected: %v", c, 1)
	}

	var item []byte
	for i := range pq.current.Items() {
		item = i
		break
	}

	b := buffer.NewFromSlice(item)

	rc2 := &record.Record{}
	if err := rc2.Read(b); err != nil {
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

	if i := rc.IsCheckpoint(); i != (action == record.ActionCheckpoint) {
		t.Errorf("incorrect is checkpoint: %v, expected: %v", i, action == record.ActionCheckpoint)
	}
}

func TestPageQueue__Next(t *testing.T) {
	t.Parallel()

	pq := New(256)

	lsn := record.NewLogSequenceNumber(1)
	tid := value.NewTenantID()
	txid := record.NewTransactionID(2)
	now := record.NewTime(time.UnixMicro(1234567890))
	action := record.ActionInsert
	collection := record.CollectionEvent
	data := []byte("abcdef")
	reverse := []byte("fghij")

	rc := record.New(lsn, tid, txid, now, action, collection, data, reverse)

	if err := pq.Append(rc); err != nil {
		t.Fatal(err)
	}

	if err := pq.Append(rc); err != nil {
		t.Fatal(err)
	}

	if err := pq.Append(rc); err != nil {
		t.Fatal(err)
	}

	if c := pq.Count(); c != 2 {
		t.Fatalf("incorrect count: %v, expected: %v", c, 1)
	}
}
