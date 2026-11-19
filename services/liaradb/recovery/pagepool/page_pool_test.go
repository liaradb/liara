package pagepool

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/recovery/logpage"
)

func TestPagePool(t *testing.T) {
	pp := New(128, 4)

	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	empty := make([]byte, logpage.HeaderSize)

	p0 := pp.Get()
	d0 := p0.Data()
	want := slices.Clip(d0)

	copy(d0, data)

	pp.Put(p0)

	p1 := pp.Get()

	d1 := p1.Data()
	if !slices.Equal(d1, want) {
		t.Errorf("incorrect header: %v, expected: %v", d1, empty)
	}
	if !slices.Equal(d0, want) {
		t.Errorf("incorrect header: %v, expected: %v", d0, empty)
	}

	pp.Put(p1)
}
