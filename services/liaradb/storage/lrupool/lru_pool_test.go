package lrupool

import "testing"

func TestNew(t *testing.T) {
	p := New()
	if p == nil {
		t.Fatal("should create MapQueue")
	}
}
