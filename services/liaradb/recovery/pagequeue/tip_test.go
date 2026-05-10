package pagequeue

import (
	"testing"

	"github.com/liaradb/liaradb/encoder/page"
)

const (
	pageSize       = 64
	headerSize     = 22
	slotHeaderSize = 4
)

func TestTip(t *testing.T) {
	t.Parallel()

	current := page.New(pageSize, headerSize, slotHeaderSize)
	tip := NewTip(pageSize, headerSize, slotHeaderSize, current)
	var want int16 = 128
	s := tip.Span(want)

	// TODO: Fix type
	l := int16(s.Length())
	if l != want {
		t.Errorf("incorrect length: %v, expected: %v", l, want)
	}

	pages := tip.Pages()
	if l := len(pages); l != 3 {
		t.Errorf("incorrect length: %v, expected: %v", l, 3)
	}
}
