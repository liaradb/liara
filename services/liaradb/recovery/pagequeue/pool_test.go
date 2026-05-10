package pagequeue

import (
	"slices"
	"testing"
)

func TestPool(t *testing.T) {
	pp := NewPool(1024, 8, 4)

	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	empty := make([]byte, 8)

	p0 := pp.Get()
	h0 := p0.Header()
	copy(h0, data)

	pp.Put(p0)

	p1 := pp.Get()

	h1 := p1.Header()
	if !slices.Equal(h1, empty) {
		t.Errorf("incorrect header: %v, expected: %v", h1, empty)
	}
	if !slices.Equal(h0, empty) {
		t.Errorf("incorrect header: %v, expected: %v", h0, empty)
	}

	pp.Put(p1)
}
