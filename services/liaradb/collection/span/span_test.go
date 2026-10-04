package span

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
	"github.com/liaradb/liaradb/util/testing/logtesting"
)

func TestSpan_Write(t *testing.T) {
	t.Parallel()

	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	tr0 := &testRecord{
		id:   base.Uint64(2),
		data: slices.Clone(want),
	}

	size := float64(tr0.Size()) / 2
	a, b := int(math.Floor(size)), int(math.Ceil(size))

	s := New(&logtesting.MockLog{})
	s.Append(&testBufferPage{}, link.SlotID(0), make([]byte, FragmentHeaderSize), make([]byte, a))
	s.Append(&testBufferPage{}, link.SlotID(0), make([]byte, FragmentHeaderSize), make([]byte, b))
	s.InitIndexes()

	if err := tr0.Write(s); err != nil {
		t.Fatal(err)
	}

	s.Commit()
	if err := s.SeekStart(); err != nil {
		t.Fatal(err)
	}

	tr1 := &testRecord{
		data: make([]byte, 11),
	}

	if err := tr1.Read(s); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tr1, tr0) {
		t.Errorf("incorrect record: %v, expected: %v", tr1, tr0)
	}
}

func TestSpan_Length(t *testing.T) {
	t.Parallel()

	s := New(&logtesting.MockLog{})
	s.Append(&testBufferPage{}, link.SlotID(0), make([]byte, FragmentHeaderSize), make([]byte, 100))
	s.Append(&testBufferPage{}, link.SlotID(0), make([]byte, FragmentHeaderSize), make([]byte, 10000))
	s.InitIndexes()

	want := 100 + 10000
	if l := s.Length(); l != want {
		t.Errorf("incorrect length: %v, expected: %v", l, want)
	}
}

func TestSpan_InitIndexes(t *testing.T) {
	t.Parallel()

	t.Run("default", func(t *testing.T) {
		s := New(&logtesting.MockLog{})
		s.InitIndexes()
	})

	// TODO: How do we test other cases?
}

func TestSpan_Read__Invalid(t *testing.T) {
	t.Parallel()

	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	tr0 := &testRecord{
		id:   base.Uint64(2),
		data: slices.Clone(want),
	}

	size := float64(tr0.Size()) / 2
	a, b := int(math.Floor(size)), int(math.Ceil(size))

	s := New(&logtesting.MockLog{})
	s.Append(&testBufferPage{}, link.SlotID(0), make([]byte, FragmentHeaderSize), make([]byte, a))
	s.Append(&testBufferPage{}, link.SlotID(0), make([]byte, FragmentHeaderSize), make([]byte, b))
	s.InitIndexes()

	if err := tr0.Write(s); err != nil {
		t.Fatal(err)
	}

	s.Commit()
	if err := s.SeekStart(); err != nil {
		t.Fatal(err)
	}

	s.fragments[1].buffer.Bytes()[8] = 0

	tr1 := &testRecord{
		data: make([]byte, 11),
	}

	if err := tr1.Read(s); !errors.Is(err, page.ErrInvalidCRC) {
		t.Errorf("incorrect error: %v, expected: %v", err, page.ErrInvalidCRC)
	}
}

func TestSpan_Bytes(t *testing.T) {
	t.Parallel()

	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	s := New(&logtesting.MockLog{})
	s.Append(&testBufferPage{}, link.SlotID(0), make([]byte, FragmentHeaderSize), make([]byte, 5))
	s.Append(&testBufferPage{}, link.SlotID(0), make([]byte, FragmentHeaderSize), make([]byte, len(want)-5))
	s.InitIndexes()

	if _, err := s.Write(want); err != nil {
		t.Fatal(err)
	}

	s.Commit()
	s.SeekStart()

	if d, err := s.Bytes(); err != nil {
		t.Fatal(err)
	} else if !slices.Equal(d, want) {
		t.Errorf("incorrect result: %v, expected: %v", d, want)
	}
}

type testBufferPage struct {
}

func (t *testBufferPage) BlockID() link.BlockID { return link.BlockID{} }
func (t *testBufferPage) SetLogSequenceNumber(logpage.LogSequenceNumber) {
}
func (t *testBufferPage) Commit(size int) {
}
