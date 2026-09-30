package page

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/util/testing/should"
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

	p.Commit(8)

	header, data = p.Slot(0)

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
	p.Commit(18)

	header, data := p.Next(2)
	if l := len(header); l != 0 {
		t.Errorf("incorrect header: %v, expected: %v", l, 0)
	}

	if l := len(data); l != 0 {
		t.Errorf("incorrect data: %v, expected: %v", l, 0)
	}
}

func TestPage_NextMustFit(t *testing.T) {
	t.Parallel()

	p := New(32, 4, 4)

	header, data, _ := p.NextMustFit(8)

	if l := len(header); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if l := len(data); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	p.Commit(8)

	header, data = p.Slot(0)

	if l := len(header); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if l := len(data); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}
}

func TestPage_Insert(t *testing.T) {
	t.Parallel()

	p := New(128, 4, 4)

	{
		header, data := p.Next(8)

		if l := len(header); l != 4 {
			t.Errorf("incorrect length: %v, expected: %v", l, 4)
		}

		if l := len(data); l != 8 {
			t.Errorf("incorrect length: %v, expected: %v", l, 8)
		}

		header[0] = 1
		data[0] = 2

		p.Commit(8)
	}

	{
		header, data := p.Next(8)

		if l := len(header); l != 4 {
			t.Errorf("incorrect length: %v, expected: %v", l, 4)
		}

		if l := len(data); l != 8 {
			t.Errorf("incorrect length: %v, expected: %v", l, 8)
		}

		header[0] = 3
		data[0] = 4

		p.Insert(8, 0)
	}

	{
		header, data := p.Slot(0)

		if l := len(header); l != 4 {
			t.Errorf("incorrect length: %v, expected: %v", l, 4)
		}

		if l := len(data); l != 8 {
			t.Errorf("incorrect length: %v, expected: %v", l, 8)
		}

		if h := header[0]; h != 3 {
			t.Errorf("incorrect header: %v, expected: %v", h, 3)
		}

		if d := data[0]; d != 4 {
			t.Errorf("incorrect data: %v, expected: %v", d, 4)
		}
	}

	{
		header, data := p.Slot(1)

		if l := len(header); l != 4 {
			t.Errorf("incorrect length: %v, expected: %v", l, 4)
		}

		if l := len(data); l != 8 {
			t.Errorf("incorrect length: %v, expected: %v", l, 8)
		}

		if h := header[0]; h != 1 {
			t.Errorf("incorrect header: %v, expected: %v", h, 1)
		}

		if d := data[0]; d != 2 {
			t.Errorf("incorrect data: %v, expected: %v", d, 2)
		}
	}
}

func TestPage_Slot(t *testing.T) {
	t.Parallel()

	p := New(32, 4, 4)

	should.Panic(t, "slot should not exist", func() {
		_, _ = p.Slot(0)
	})

	c := 0
	for range p.Slots() {
		c++
	}

	if c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}

	_, _ = p.Next(8)
	p.Commit(8)

	_, _ = p.Next(8)
	p.Commit(8)

	_, _ = p.Slot(0)
	_, _ = p.Slot(1)

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
