package pagequeue

import (
	"testing"

	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/recovery/record"
)

func TestTip(t *testing.T) {
	t.Parallel()

	var pageSize int64 = 64
	current := page.New(pageSize)
	current.Init(
		action.NewPageIDFromSize(pageSize, 0),
		0,
		record.NewLength(0))

	tip := NewTip(current)
	var want int16 = 128
	s := tip.Span(want)

	// TODO: Fix type
	l := int16(s.Length())
	if l != want {
		t.Errorf("incorrect length: %v, expected: %v", l, want)
	}

	pages := tip.Pages()
	if l := len(pages); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if c := pages[0]; c != current {
		t.Errorf("incorrect first page: %v, expected: %v", c, current)
	}
}
