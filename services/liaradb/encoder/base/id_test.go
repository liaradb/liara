package base

import (
	"slices"
	"testing"
	"uuid"

	"github.com/liaradb/liaradb/encoder/buffer"
)

func TestID_NewIDFromUUID(t *testing.T) {
	t.Parallel()

	id := uuid.NewV7()
	b := NewIDFromUUID(id)
	b.Bytes()

	if !slices.Equal(b.Bytes(), id[:]) {
		t.Errorf("incorrect value: %v, expected: %v", b.Bytes(), id[:])
	}

	if s := b.Size(); s != BaseIDSize {
		t.Errorf("incorrect size: %v, expected: %v", s, BaseIDSize)
	}
}

func TestID_NewIDFromString(t *testing.T) {
	t.Parallel()

	id := uuid.NewV7()
	b, err := NewIDFromString(id.String())
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(b.Bytes(), id[:]) {
		t.Errorf("incorrect value: %v, expected: %v", b.Bytes(), id[:])
	}

	if _, err := NewIDFromString("abcd"); err == nil {
		t.Error("should only parse valid UUID values")
	}

	if s := b.String(); s != id.String() {
		t.Errorf("incorrect string: %v, expected: %v", s, id.String())
	}
}

func TestID_Read_Write(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	b := NewID()
	if err := b.Write(w); err != nil {
		t.Fatal(err)
	}

	var b0 ID
	if err := b0.Read(r); err != nil {
		t.Fatal(err)
	}

	if b != b0 {
		t.Errorf("incorrect value: %v, expected: %v", b, b0)
	}
}

func TestID_Read_Write__Error(t *testing.T) {
	t.Parallel()

	buf := buffer.New(15)

	b := NewID()
	if err := b.Write(buf); err == nil {
		t.Fatal("should return error")
	}

	buf.SeekStart()

	if err := b.Read(buf); err == nil {
		t.Fatal("should return error")
	}
}
