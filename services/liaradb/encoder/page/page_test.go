package page

import (
	"slices"
	"testing"
)

func TestPage_New(t *testing.T) {
	t.Parallel()

	p := New(32, 4, 4)

	want := make([]byte, 32)
	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}
}

func TestPage_NewFromSlice(t *testing.T) {
	t.Parallel()

	want := []byte{1, 2, 3, 4}
	p := NewFromSlice(want, 4, 4)

	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}
}

func TestPage_Fill(t *testing.T) {
	t.Parallel()

	p := New(4, 4, 4)

	want := []byte{1, 2, 3, 4}
	p.Fill(want)

	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}

	p.Clear()

	want = make([]byte, 4)
	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}
}

func TestPage_Next(t *testing.T) {
	t.Parallel()

	p := New(32, 4, 4)

	header, data := p.Next(8)

	if l := len(header); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if l := len(data); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	i, ok := p.Commit(8)
	if !ok {
		t.Error("should commit")
	}

	header, data, ok = p.Slot(i)
	if !ok {
		t.Error("should get slot")
	}

	if l := len(header); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if l := len(data); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}
}
