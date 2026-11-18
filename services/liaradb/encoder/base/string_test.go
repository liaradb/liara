package base

import (
	"slices"
	"testing"
)

func TestString(t *testing.T) {
	t.Parallel()

	s0 := String("name")

	r, w := newReaderWriter()

	if err := s0.Write(w); err != nil {
		t.Fatal(err)
	}

	var s1 String
	if err := s1.Read(r); err != nil {
		t.Fatal(err)
	}

	if s1.String() != s0.String() {
		t.Errorf("incorrect result: %v, expected: %v", s1.String(), s0.String())
	}

	if b := s1.Bytes(); !slices.Equal(b, []byte("name")) {
		t.Errorf("incorrect bytes: %v, expected: %v", b, []byte("name"))
	}

	if l := s1.Length(); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if s := s1.Size(); s != 4+4 {
		t.Errorf("incorrect size: %v, expected: %v", s, 4+4)
	}
}

func TestString_Read_Write(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	b := String("name")
	if err := b.Write(w); err != nil {
		t.Fatal(err)
	}

	var b0 String
	if err := b0.Read(r); err != nil {
		t.Fatal(err)
	}

	if b != b0 {
		t.Errorf("incorrect value: %v, expected: %v", b, b0)
	}
}
