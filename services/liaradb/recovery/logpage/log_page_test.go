package logpage

import (
	"slices"
	"testing"
)

func TestLogPage_AddHandler(t *testing.T) {
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
	p := New(256, 8)

	a := 0
	p.AddHandler(func() {
		a++
	})

	b := 0
	p.AddHandler(func() {
		b++
	})

	p.AddHandler(nil)

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

func TestLogPage_Reset(t *testing.T) {
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

func TestLogPage_Shadow(t *testing.T) {
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
