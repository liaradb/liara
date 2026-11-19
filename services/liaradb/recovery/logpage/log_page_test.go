package logpage

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/encoder/buffer"
)

func TestLogPage_AddHandler(t *testing.T) {
	t.Parallel()

	p := New(256, 8)

	c := 0
	p.AddHandler(func() {
		c++
	})

	p.Complete()
	if c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}

	p.Complete()
	if c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}
}

func TestLogPage_AddHandler__Multiple(t *testing.T) {
	t.Parallel()

	p := New(256, 8)

	a := 0
	p.AddHandler(func() {
		a++
	})

	b := 0
	p.AddHandler(func() {
		b++
	})

	p.Complete()
	if a != 1 {
		t.Errorf("incorrect count: %v, expected: %v", a, 1)
	}

	if b != 1 {
		t.Errorf("incorrect count: %v, expected: %v", b, 1)
	}

	p.Complete()
	if a != 1 {
		t.Errorf("incorrect count: %v, expected: %v", a, 1)
	}

	if b != 1 {
		t.Errorf("incorrect count: %v, expected: %v", b, 1)
	}
}

func TestLogPage_Clear(t *testing.T) {
	t.Parallel()

	p := New(258, 8)

	if tlid := p.TimeLineID(); tlid != 0 {
		t.Errorf("incorrect timeline id: %v, expected: %v", tlid, 0)
	}

	p.SetTimeLineID(1)
	if tlid := p.TimeLineID(); tlid != 1 {
		t.Errorf("incorrect timeline id: %v, expected: %v", tlid, 1)
	}

	p.Clear()
	if tlid := p.TimeLineID(); tlid != 0 {
		t.Errorf("incorrect timeline id: %v, expected: %v", tlid, 0)
	}
}

func TestLogPage_Reset(t *testing.T) {
	t.Parallel()

	p := New(256, 8)

	c := 0
	p.AddHandler(func() {
		c++
	})

	p.Reset()

	p.Complete()
	if c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}
}

func TestLogPage_Replace(t *testing.T) {
	t.Parallel()

	p := New(256, 8)

	c := 0
	p.AddHandler(func() {
		c++
	})

	data := make([]byte, 256)

	want := []byte{1, 2, 3, 4}
	copy(data[:4], want)

	if slices.Equal(p.Data()[:4], want) {
		t.Error("should not equal new data")
	}

	b := buffer.NewFromSlice(data)
	if err := p.Replace(b); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(p.Data()[:4], want) {
		t.Errorf("incorrect data: %v, expected: %v", p.Data()[:4], data)
	}

	p.Complete()
	if c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}
}

func TestLogPage_Replace__Error(t *testing.T) {
	t.Parallel()

	p := New(256, 8)

	b := buffer.New(0)
	if err := p.Replace(b); err == nil {
		t.Error("should return error")
	}
}

func TestLogPage_Shadow(t *testing.T) {
	t.Parallel()

	p := New(256, 8)

	c := 0
	p.AddHandler(func() {
		c++
	})

	data := []byte{1, 2, 3, 4}
	copy(p.Data(), data)

	if !slices.Equal(p.Data()[:4], data) {
		t.Errorf("incorrect data: %v, expected: %v", p.Data()[:4], data)
	}

	s := New(256, 8)
	s.Shadow(p)

	if !slices.Equal(s.Data()[:4], data) {
		t.Errorf("incorrect data: %v, expected: %v", s.Data()[:4], data)
	}

	p.Complete()
	if c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}

	s.Complete()
	if c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}
}

func TestLogPage_Next_Commit(t *testing.T) {
	t.Parallel()

	p := New(40, 4)

	header, data := p.Next(8)

	if l := len(header); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}
	if l := len(data); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	header[0] = 12
	data[1] = 34

	p.Commit(8)

	c := 0
	for header, data = range p.Slots() {
		c++
		break
	}

	if c != 1 {
		t.Fatalf("incorrect slots: %v, expected: %v", c, 1)
	}
	if l := len(header); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}
	if l := len(data); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}
	if header[0] != 12 {
		t.Errorf("incorrect header: %v, expected: %v", header[0], 12)
	}
	if data[1] != 34 {
		t.Errorf("incorrect data: %v, expected: %v", data[0], 34)
	}
}

func TestLogPage_Slots(t *testing.T) {
	t.Parallel()

	p := New(128, 4)

	type row struct {
		h []byte
		d []byte
	}

	want := []row{{
		h: []byte("head"),
		d: []byte("bodydata"),
	}, {
		h: []byte("word"),
		d: []byte("twowords"),
	}}

	for _, item := range want {
		h, d := p.Next(int16(len(item.d)))
		copy(h, item.h)
		copy(d, item.d)
		p.Commit(int16(len(item.d)))
	}

	result := make([]row, 0, len(want))

	for h, d := range p.Slots() {
		result = append(result, row{h, d})
	}

	if !slices.EqualFunc(result, want, func(a, b row) bool {
		return slices.Equal(a.h, b.h) && slices.Equal(a.d, b.d)
	}) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}

	slices.Reverse(want)
	reverse := make([]row, 0, len(want))

	for h, d := range p.SlotsReverse() {
		reverse = append(reverse, row{h, d})
	}

	if !slices.EqualFunc(reverse, want, func(a, b row) bool {
		return slices.Equal(a.h, b.h) && slices.Equal(a.d, b.d)
	}) {
		t.Errorf("incorrect result: %v, expected: %v", reverse, want)
	}
}
