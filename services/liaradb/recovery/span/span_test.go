package span

import (
	"errors"
	"math"
	"reflect"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/logpage"
)

func TestSpan_Write(t *testing.T) {
	t.Parallel()

	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	tr0 := &testRecord{
		lsn:  logpage.NewLogSequenceNumber(2),
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

func TestSpan_Read__Invalid(t *testing.T) {
	t.Parallel()

	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	tr0 := &testRecord{
		lsn:  logpage.NewLogSequenceNumber(2),
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
