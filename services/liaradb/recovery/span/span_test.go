package span

import (
	"errors"
	"math"
	"reflect"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/encoder/page"
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

	s := Span{}
	s.Append(make([]byte, FragmentHeaderSize), make([]byte, a))
	s.Append(make([]byte, FragmentHeaderSize), make([]byte, b))
	s.InitIndexes()

	if err := tr0.Write(s); err != nil {
		t.Fatal(err)
	}

	s.Commit()
	s.SeekStart()

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

func TestSpan_Read__Invalid(t *testing.T) {
	t.Parallel()

	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	tr0 := &testRecord{
		id:   base.Uint64(2),
		data: slices.Clone(want),
	}

	size := float64(tr0.Size()) / 2
	a, b := int(math.Floor(size)), int(math.Ceil(size))

	s := Span{}
	s.Append(make([]byte, FragmentHeaderSize), make([]byte, a))
	f := s.Append(make([]byte, FragmentHeaderSize), make([]byte, b))
	s.InitIndexes()

	if err := tr0.Write(s); err != nil {
		t.Fatal(err)
	}

	s.Commit()
	s.SeekStart()

	f.buffer.Bytes()[8] = 0

	tr1 := &testRecord{
		data: make([]byte, 11),
	}

	if err := tr1.Read(s); !errors.Is(err, page.ErrInvalidCRC) {
		t.Errorf("incorrect error: %v, expected: %v", err, page.ErrInvalidCRC)
	}
}

func TestSpan_Length(t *testing.T) {
	t.Parallel()

	s := Span{}

	if l, want := s.Length(), 0; l != want {
		t.Errorf("incorrect length: %v, expected: %v", l, want)
	}

	s.Append(make([]byte, FragmentHeaderSize), make([]byte, 10))

	if l, want := s.Length(), 10; l != want {
		t.Errorf("incorrect length: %v, expected: %v", l, want)
	}

	s.Append(make([]byte, FragmentHeaderSize), make([]byte, 100))

	if l, want := s.Length(), 110; l != want {
		t.Errorf("incorrect length: %v, expected: %v", l, want)
	}
}

func TestSpan_Count(t *testing.T) {
	t.Parallel()

	s := Span{}

	f0 := s.Append(make([]byte, FragmentHeaderSize), make([]byte, 1))
	s.InitIndexes()

	if i, want := f0.Count(), int16(1); i != want {
		t.Errorf("incorrect count: %v, expected: %v", i, want)
	}

	f1 := s.Append(make([]byte, FragmentHeaderSize), make([]byte, 1))
	s.InitIndexes()

	if i, want := f0.Count(), int16(2); i != want {
		t.Errorf("incorrect count: %v, expected: %v", i, want)
	}

	if i, want := f1.Count(), int16(2); i != want {
		t.Errorf("incorrect count: %v, expected: %v", i, want)
	}

	f2 := s.Append(make([]byte, FragmentHeaderSize), make([]byte, 1))
	s.InitIndexes()

	if i, want := f0.Count(), int16(3); i != want {
		t.Errorf("incorrect count: %v, expected: %v", i, want)
	}

	if i, want := f1.Count(), int16(3); i != want {
		t.Errorf("incorrect count: %v, expected: %v", i, want)
	}

	if i, want := f2.Count(), int16(3); i != want {
		t.Errorf("incorrect count: %v, expected: %v", i, want)
	}
}

func TestSpan_Reverse(t *testing.T) {
	t.Parallel()

	s := Span{}
	f0 := s.Append(make([]byte, FragmentHeaderSize), make([]byte, 1))
	f1 := s.Append(make([]byte, FragmentHeaderSize), make([]byte, 1))
	f2 := s.Append(make([]byte, FragmentHeaderSize), make([]byte, 1))

	s.InitIndexes()

	if i, want := f0.Index(), int16(2); i != want {
		t.Errorf("incorrect index: %v, expected: %v", i, want)
	}

	if i, want := f1.Index(), int16(1); i != want {
		t.Errorf("incorrect index: %v, expected: %v", i, want)
	}

	if i, want := f2.Index(), int16(0); i != want {
		t.Errorf("incorrect index: %v, expected: %v", i, want)
	}

	s.Reverse()
	s.InitIndexes()

	if i, want := f0.Index(), int16(0); i != want {
		t.Errorf("incorrect index: %v, expected: %v", i, want)
	}

	if i, want := f1.Index(), int16(1); i != want {
		t.Errorf("incorrect index: %v, expected: %v", i, want)
	}

	if i, want := f2.Index(), int16(2); i != want {
		t.Errorf("incorrect index: %v, expected: %v", i, want)
	}
}
