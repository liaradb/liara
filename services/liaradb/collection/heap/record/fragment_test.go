package record

import (
	"math"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/txtesting"
)

func TestFragment(t *testing.T) {
	t.Parallel()

	tx := txtesting.Transaction{}
	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	tr0 := &testRecord{
		id:   base.Uint64(2),
		data: slices.Clone(want),
	}

	size := float64(tr0.Size()) / 2
	a, b := int(math.Floor(size)), int(math.Ceil(size))

	r := New()
	f0 := r.Append(&testBufferPage{
		fileBlockID: link.NewFileBlockID("", 1),
	}, link.SlotID(10), make([]byte, fragmentHeaderSize), make([]byte, a))
	f1 := r.Append(&testBufferPage{
		fileBlockID: link.NewFileBlockID("", 2),
	}, link.SlotID(20), make([]byte, fragmentHeaderSize), make([]byte, b))
	r.InitIndexes()

	r.CommitAppend(t.Context(), &tx)

	if nbid, bid := f0.NextBlockID(), f1.RecordID().BlockID(); nbid != bid {
		t.Errorf("incorrect next block id: %v, expected: %v", nbid, bid)
	}

	if nsid, sid := f0.NextSlotID(), f1.RecordID().SlotID(); nsid != sid {
		t.Errorf("incorrect next slot id: %v, expected: %v", nsid, sid)
	}
}

func TestFragment_Release(t *testing.T) {
	t.Parallel()

	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	tr0 := &testRecord{
		id:   base.Uint64(2),
		data: slices.Clone(want),
	}

	size := float64(tr0.Size()) / 2
	a, b := int(math.Floor(size)), int(math.Ceil(size))

	bp0 := testBufferPage{}
	bp1 := testBufferPage{}

	r := New()
	_ = r.Append(&bp0, link.SlotID(0), make([]byte, fragmentHeaderSize), make([]byte, a))
	_ = r.Append(&bp1, link.SlotID(0), make([]byte, fragmentHeaderSize), make([]byte, b))

	if bp0.isReleased() {
		t.Error("should not be released")
	}

	if bp1.isReleased() {
		t.Error("should not be released")
	}

	r.Release()

	if !bp0.isReleased() {
		t.Error("should be released")
	}

	if !bp1.isReleased() {
		t.Error("should be released")
	}
}

func TestFragment_Delete(t *testing.T) {
	t.Parallel()

	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	tr0 := &testRecord{
		id:   base.Uint64(2),
		data: slices.Clone(want),
	}

	size := float64(tr0.Size()) / 2
	a, b := int(math.Floor(size)), int(math.Ceil(size))

	bp0 := testBufferPage{}
	bp1 := testBufferPage{}

	r := New()
	_ = r.Append(&bp0, link.SlotID(0), make([]byte, fragmentHeaderSize), make([]byte, a))
	_ = r.Append(&bp1, link.SlotID(0), make([]byte, fragmentHeaderSize), make([]byte, b))

	if bp0.isDeleted() {
		t.Error("should not be deleted")
	}

	if bp1.isDeleted() {
		t.Error("should not be deleted")
	}

	r.Delete()

	if !bp0.isDeleted() {
		t.Error("should be deleted")
	}

	if !bp1.isDeleted() {
		t.Error("should be deleted")
	}
}
