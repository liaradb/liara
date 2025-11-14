package rawpage

import (
	"io"
	"slices"
	"testing"
)

func TestRawPage(t *testing.T) {
	p := New(make([]byte, 256))
	v0 := []byte{1, 2, 3, 4, 5}
	v1 := []byte{6, 7, 8, 9, 10}

	i, b0, ok := p.Append(16)
	if !ok {
		t.Error("should get a buffer")
	} else if i != 0 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	if _, err := b0.Write(v0); err != nil {
		t.Error(err)
	}

	i, b1, ok := p.Append(16)
	if !ok {
		t.Error("should get a buffer")
	} else if i != 1 {
		t.Errorf("incorrect index: %v, expected: %v", i, 1)
	}

	if _, err := b1.Write(v1); err != nil {
		t.Error(err)
	}

	if _, err := b0.Seek(0, io.SeekStart); err != nil {
		t.Error(err)
	}

	r0 := make([]byte, 5)
	if _, err := b0.Read(r0); err != nil {
		t.Error(err)
	}

	if !slices.Equal(r0, v0) {
		t.Errorf("incorrect result: %v, expected: %v", r0, v0)
	}

	if _, err := b1.Seek(0, io.SeekStart); err != nil {
		t.Error(err)
	}

	r1 := make([]byte, 5)
	if _, err := b1.Read(r1); err != nil {
		t.Error(err)
	}

	if !slices.Equal(r1, v1) {
		t.Errorf("incorrect result: %v, expected: %v", r1, v1)
	}
}
