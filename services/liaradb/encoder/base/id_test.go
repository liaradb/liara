package base

import (
	"io"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/liaradb/liaradb/encoder/buffer"
)

func TestID_NewIDFromUUID(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	b := NewIDFromUUID(id)
	b.Bytes()

	if !slices.Equal(b.Bytes(), id[:]) {
		t.Errorf("incorrect value: %v, expected: %v", b.Bytes(), id[:])
	}
}

func TestID_NewIDFromString(t *testing.T) {
	t.Parallel()

	id := uuid.New()
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
}

func TestID__Remainder(t *testing.T) {
	t.Parallel()

	b := NewID()

	data := make([]byte, 20)
	data0 := b.WriteData(data)

	if l := len(data0); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	b0 := ID{}
	data1 := b0.ReadData(data)

	if l := len(data1); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if b0 != b {
		t.Errorf("incorrect value: %v, expected: %v", b0, b)
	}

	if b.String() != b0.String() {
		t.Errorf("incorrect result: %v, expected: %v", b.String(), b0.String())
	}

	if bytes := b.Bytes(); !slices.Equal(bytes, data[:16]) {
		t.Errorf("incorrect bytes: %v, expected: %v", bytes, data[:16])
	}

	if s := b.Size(); s != 16 {
		t.Errorf("incorrect size: %v, expected: %v", s, 16)
	}
}

func TestID_Read_Write(t *testing.T) {
	t.Parallel()

	buf := buffer.New(16)

	b := NewID()
	if err := b.Write(buf); err != nil {
		t.Fatal(err)
	}

	if _, err := buf.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	var b0 ID
	if err := b0.Read(buf); err != nil {
		t.Fatal(err)
	}

	if b != b0 {
		t.Errorf("incorrect value: %v, expected: %v", b, b0)
	}
}
