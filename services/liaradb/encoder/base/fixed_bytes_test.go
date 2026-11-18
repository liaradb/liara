package base

import (
	"slices"
	"testing"
)

func TestFixedBytesBytes(t *testing.T) {
	t.Parallel()

	want := []byte("name")
	size := len(want)
	fb0 := NewFixedBytes(want)

	r, w := newReaderWriter()

	if err := fb0.Write(w); err != nil {
		t.Fatal(err)
	}

	fb1 := NewFixedBytesFromSize(size)
	if err := fb1.Read(r); err != nil {
		t.Fatal(err)
	}

	if !fb1.Compare(fb0) {
		t.Error("should be equal")
	}

	if fb1.String() != fb0.String() {
		t.Errorf("incorrect result: %v, expected: %v", fb1.String(), fb0.String())
	}

	if b := fb1.Value(); !slices.Equal(b, []byte("name")) {
		t.Errorf("incorrect bytes: %v, expected: %v", b, []byte("name"))
	}

	if l := fb1.Length(); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if s := fb1.Size(); s != 4 {
		t.Errorf("incorrect size: %v, expected: %v", s, 4)
	}
}

func TestFixedBytes_Read_Write(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	b := NewFixedBytes([]byte{1, 2, 3, 4})
	if err := b.Write(w); err != nil {
		t.Fatal(err)
	}

	b0 := NewFixedBytesFromSize(4)
	if err := b0.Read(r); err != nil {
		t.Fatal(err)
	}

	if !b0.Compare(b) {
		t.Errorf("incorrect value: %v, expected: %v", b0.Value(), b.Value())
	}
}
