package value

import (
	"io"
	"slices"
	"testing"
)

func TestGlobalVersion(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var gv GlobalVersion = NewGlobalVersion(1)
	if err := gv.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := gv.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var gv2 GlobalVersion
	if err := gv2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	s1, s2 := gv.String(), gv2.String()
	if s1 != s2 {
		t.Errorf("incorrect value: %v, expected: %v", s2, s1)
	}
}

func TestGlobalVersion_Increment(t *testing.T) {
	t.Parallel()

	gv0 := NewGlobalVersion(1)
	gv1 := gv0.Increment()
	gv2 := gv1.Increment()

	if v := gv0.Value(); v != 1 {
		t.Errorf("incorrect value: %v, expected: %v", v, 1)
	}
	if v := gv1.Value(); v != 2 {
		t.Errorf("incorrect value: %v, expected: %v", v, 2)
	}
	if v := gv2.Value(); v != 3 {
		t.Errorf("incorrect value: %v, expected: %v", v, 3)
	}
}

func TestGlobalVersion_Bytes(t *testing.T) {
	t.Parallel()

	gv0 := NewGlobalVersion(1)
	gv1 := gv0.Increment()

	want0 := [8]byte{0, 0, 0, 0, 0, 0, 0, 1}
	want1 := [8]byte{0, 0, 0, 0, 0, 0, 0, 2}

	if b := gv0.Bytes(); !slices.Equal(b, want0[:]) {
		t.Errorf("incorrect bytes: %v, expected: %v", b, want0[:])
	}

	if b := gv1.Bytes(); !slices.Equal(b, want1[:]) {
		t.Errorf("incorrect bytes: %v, expected: %v", b, want1[:])
	}
}
