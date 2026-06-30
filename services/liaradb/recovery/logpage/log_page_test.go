package logpage

import "testing"

func TestLogPage_Handler(t *testing.T) {
	p := New(256, 8)

	c := 0
	p.SetHandler(func() {
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

func TestLogPage_Reset(t *testing.T) {
	p := New(256, 8)

	c := 0
	p.SetHandler(func() {
		c++
	})

	p.Reset()

	p.Complete()
	if c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}
}
