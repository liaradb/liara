package pagequeue

import (
	"testing"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/span"
)

const (
	pageSize       = 64
	largePageSize  = 256
	writeQueueSize = 100
)

func TestTip(t *testing.T) {
	t.Parallel()
	pool := pagepool.New(pageSize, span.FragmentHeaderSize)

	current := logpage.New(pageSize, span.FragmentHeaderSize)
	tip := NewTip(pool, current)
	var want int16 = 128
	s := tip.Span(want)

	// TODO: Fix type
	l := int16(s.Length())
	if l != want {
		t.Errorf("incorrect length: %v, expected: %v", l, want)
	}

	pages, ok := tip.Commit(record.NewLogSequenceNumber(0), nil)
	if !ok {
		t.Error("should commit")
	}

	if l := len(pages); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}
}
