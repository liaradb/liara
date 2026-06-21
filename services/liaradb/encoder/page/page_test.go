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

	if ok := p.Commit(8); !ok {
		t.Error("should commit")
	}

	header, data, ok := p.Slot(0)
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

func TestPage_Next__Empty(t *testing.T) {
	t.Parallel()

	p := New(32, 4, 4)
	if ok := p.Commit(18); !ok {
		t.Error("should commit")
	}

	header, data := p.Next(2)
	if l := len(header); l != 0 {
		t.Errorf("incorrect header: %v, expected: %v", l, 0)
	}

	if l := len(data); l != 0 {
		t.Errorf("incorrect data: %v, expected: %v", l, 0)
	}
}

func TestPage_Slot(t *testing.T) {
	t.Parallel()

	p := New(32, 4, 4)
	if _, _, ok := p.Slot(0); ok {
		t.Error("slot should not exist")
	}

	c := 0
	for range p.Slots() {
		c++
	}

	if c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}

	_, _ = p.Next(8)
	if ok := p.Commit(8); !ok {
		t.Error("should commit")
	}

	_, _ = p.Next(8)
	if ok := p.Commit(8); !ok {
		t.Error("should commit")
	}

	if _, _, ok := p.Slot(0); !ok {
		t.Error("slot should exist")
	}

	if _, _, ok := p.Slot(1); !ok {
		t.Error("slot should exist")
	}

	c = 0
	for range p.Slots() {
		c++
	}

	if c != 2 {
		t.Errorf("incorrect count: %v, expected: %v", c, 2)
	}

	// Early return
	for range p.Slots() {
		break
	}
}
