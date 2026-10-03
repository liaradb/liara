package page

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/encoder/scan"
	"github.com/liaradb/liaradb/util/testing/should"
)

func TestPage_New(t *testing.T) {
	t.Parallel()

	p := New(32, 4, 4)

	want := make([]byte, 32)
	_ = scan.SetInt32(want, int32(MagicPage))
	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}
}

func TestPage_NewFromSlice(t *testing.T) {
	t.Parallel()

	want := []byte{0, 0, 0, 0, 1, 2, 3, 4, 0, 0, 0, 0}
	p := NewFromSlice(want, 4, 4)

	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}
}

func TestPage_Fill(t *testing.T) {
	t.Parallel()

	p := New(8, 4, 4)

	want := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	p.Fill(want)

	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}

	p.Clear()

	want = make([]byte, 8)
	_ = scan.SetInt32(want, int32(MagicPage))
	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}
}

func TestPage_Header(t *testing.T) {
	t.Parallel()

	want := []byte{0, 0, 0, 0, 1, 2, 3, 4, 0, 0, 0, 0}
	p := NewFromSlice(want, 4, 4)

	if h := p.Header(); !slices.Equal(h, want[4:8]) {
		t.Errorf("incorrect data: %v, expected: %v", h, want[4:8])
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

func TestPage_Slots(t *testing.T) {
	t.Parallel()

	p := New(256, 4, 4)

	want := []int64{101, 102, 103, 104, 105}
	for _, v := range want {
		_, b := p.Next(8)
		_ = scan.SetInt64(b, v)
		p.Commit(8)
	}

	result := make([]int64, 0, 5)
	for _, b := range p.Slots() {
		v, _ := scan.Int64(b)
		result = append(result, v)
	}

	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func TestPage_SlotsRange(t *testing.T) {
	t.Parallel()

	p := New(256, 4, 4)

	want := []int64{101, 102, 103, 104, 105}
	for _, v := range want {
		_, b := p.Next(8)
		_ = scan.SetInt64(b, v)
		p.Commit(8)
	}

	result := make([]int64, 0, 5)
	for _, b := range p.SlotsRange(1, 4) {
		v, _ := scan.Int64(b)
		result = append(result, v)
	}

	if !slices.Equal(result, want[1:4]) {
		t.Errorf("incorrect result: %v, expected: %v", result, want[1:4])
	}
}
