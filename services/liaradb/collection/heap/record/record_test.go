package record

import (
	"errors"
	"math"
	"reflect"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/txtesting"
)

func TestRecord_Write(t *testing.T) {
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
	_ = r.Append(&testBufferPage{}, link.SlotID(0), make([]byte, fragmentHeaderSize), make([]byte, a))
	_ = r.Append(&testBufferPage{}, link.SlotID(0), make([]byte, fragmentHeaderSize), make([]byte, b))
	r.InitIndexes()

	if err := tr0.Write(r.Writer()); err != nil {
		t.Fatal(err)
	}

	if err := r.CommitAppend(t.Context(), &tx); err != nil {
		t.Fatal(err)
	}

	r.SeekStart()

	tr1 := &testRecord{
		data: make([]byte, 11),
	}

	if err := tr1.Read(r); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tr1, tr0) {
		t.Errorf("incorrect record: %v, expected: %v", tr1, tr0)
	}

	if c := len(tx.Records()); c != 2 {
		t.Errorf("incorrect record count: %v, expected: %v", c, 2)
	}
}

func TestRecord_Update(t *testing.T) {
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
	_ = r.Append(&testBufferPage{}, link.SlotID(0), make([]byte, fragmentHeaderSize), make([]byte, a))
	_ = r.Append(&testBufferPage{}, link.SlotID(0), make([]byte, fragmentHeaderSize), make([]byte, b))
	r.InitIndexes()

	if err := tr0.Write(r.Writer()); err != nil {
		t.Fatal(err)
	}

	if err := r.CommitUpdate(t.Context(), &tx); err != nil {
		t.Fatal(err)
	}

	r.SeekStart()

	tr1 := &testRecord{
		data: make([]byte, 11),
	}

	if err := tr1.Read(r); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tr1, tr0) {
		t.Errorf("incorrect record: %v, expected: %v", tr1, tr0)
	}

	if c := len(tx.Records()); c != 2 {
		t.Errorf("incorrect record count: %v, expected: %v", c, 2)
	}
}

func TestRecord_Length(t *testing.T) {
	t.Parallel()

	r := New()
	_ = r.Append(&testBufferPage{}, link.SlotID(0), make([]byte, fragmentHeaderSize), make([]byte, 100))
	_ = r.Append(&testBufferPage{}, link.SlotID(0), make([]byte, fragmentHeaderSize), make([]byte, 10000))
	r.InitIndexes()

	want := 100 + 10000
	if l := r.Length(); l != want {
		t.Errorf("incorrect length: %v, expected: %v", l, want)
	}
}

func TestRecord_InitIndexes(t *testing.T) {
	t.Parallel()

	t.Run("default", func(t *testing.T) {
		s := New()
		s.InitIndexes()
	})

	// TODO: How do we test other cases?
}

func TestRecord_RecordID(t *testing.T) {
	t.Parallel()

	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	tr0 := &testRecord{
		id:   base.Uint64(2),
		data: slices.Clone(want),
	}

	size := float64(tr0.Size()) / 2
	a, b := int(math.Floor(size)), int(math.Ceil(size))

	r := New()
	f0 := r.Append(&testBufferPage{
		fileBlockID: link.NewFileBlockID("", 2)},
		link.SlotID(0),
		make([]byte, fragmentHeaderSize), make([]byte, a))
	f1 := r.Append(&testBufferPage{
		fileBlockID: link.NewFileBlockID("", 1)},
		link.SlotID(0),
		make([]byte, fragmentHeaderSize), make([]byte, b))
	r.InitIndexes()

	rid0 := f0.RecordID()
	rid1 := f1.RecordID()

	// Expect it to be sorted
	if rid := r.RecordID(); rid != rid1 {
		t.Errorf("incorrect record id: %v, expected: %v", rid, rid1)
	} else if rid == rid0 {
		t.Errorf("incorrect record id: %v, expected: %v", rid, rid1)
	}
}

func TestRecord_Read__Invalid(t *testing.T) {
	t.Parallel()

	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	tr0 := &testRecord{
		id:   base.Uint64(2),
		data: slices.Clone(want),
	}

	size := float64(tr0.Size()) / 2
	a, b := int(math.Floor(size)), int(math.Ceil(size))

	r := New()
	_ = r.Append(&testBufferPage{}, link.SlotID(0), make([]byte, fragmentHeaderSize), make([]byte, a))
	_ = r.Append(&testBufferPage{}, link.SlotID(0), make([]byte, fragmentHeaderSize), make([]byte, b))
	r.InitIndexes()

	if err := tr0.Write(r.Writer()); err != nil {
		t.Fatal(err)
	}

	if err := r.CommitAppend(t.Context(), &txtesting.Transaction{}); err != nil {
		t.Fatal(err)
	}

	r.SeekStart()

	r.fragments[1].buffer.Bytes()[8] = 0

	tr1 := &testRecord{
		data: make([]byte, 11),
	}

	if err := tr1.Read(r); !errors.Is(err, page.ErrInvalidCRC) {
		t.Errorf("incorrect error: %v, expected: %v", err, page.ErrInvalidCRC)
	}
}

func TestRecord_Bytes(t *testing.T) {
	t.Parallel()

	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	r := New()
	_ = r.Append(&testBufferPage{}, link.SlotID(0), make([]byte, fragmentHeaderSize), make([]byte, 5))
	_ = r.Append(&testBufferPage{}, link.SlotID(0), make([]byte, fragmentHeaderSize), make([]byte, len(want)-5))
	r.InitIndexes()

	if _, err := r.Writer().Write(want); err != nil {
		t.Fatal(err)
	}

	if err := r.CommitAppend(t.Context(), &txtesting.Transaction{}); err != nil {
		t.Fatal(err)
	}

	r.SeekStart()

	if d, err := r.Bytes(); err != nil {
		t.Fatal(err)
	} else if !slices.Equal(d, want) {
		t.Errorf("incorrect result: %v, expected: %v", d, want)
	}
}

type testBufferPage struct {
	fileBlockID link.FileBlockID
	released    bool
	deleted     bool
}

func (t *testBufferPage) FileBlockID() link.FileBlockID { return t.fileBlockID }

func (t *testBufferPage) SetLogSequenceNumber(logpage.LogSequenceNumber) {
}

func (t *testBufferPage) CommitAppend(size int16) {
}

func (t *testBufferPage) CommitUpdate() {
}

func (t *testBufferPage) Release() {
	t.released = true
}

func (t *testBufferPage) isReleased() bool {
	return t.released
}

func (t *testBufferPage) Slot(link.SlotID) ([]byte, []byte) {
	return nil, nil
}

func (t *testBufferPage) Delete(link.SlotID) {
	t.deleted = true
}

func (t *testBufferPage) isDeleted() bool {
	return t.deleted
}
