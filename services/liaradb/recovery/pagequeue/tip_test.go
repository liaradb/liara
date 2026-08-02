package pagequeue

import (
	"testing"

	"github.com/liaradb/liaradb/encoder/span"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pagepool"
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
	want := 128
	s := tip.Span(want)

	if l := s.Length(); l != want {
		t.Errorf("incorrect length: %v, expected: %v", l, want)
	}

	complete := 0
	pages, ok := tip.Commit(logpage.NewLogSequenceNumber(0), func() {
		complete++
	})
	if !ok {
		t.Error("should commit")
	}

	if l := len(pages); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	last := pages[len(pages)-1]
	last.Complete()
	if complete != 1 {
		t.Errorf("incorrect complete count: %v, expected: %v", complete, 1)
	}
}

func TestTip__MultiplePerPage(t *testing.T) {
	t.Parallel()
	pool := pagepool.New(pageSize, span.FragmentHeaderSize)

	current := logpage.New(pageSize, span.FragmentHeaderSize)

	var last *logpage.LogPage
	completeA := 0
	completeB := 0
	{
		tip := NewTip(pool, current)
		want := 8
		s := tip.Span(want)

		if l := s.Length(); l != want {
			t.Errorf("incorrect length: %v, expected: %v", l, want)
		}

		pages, ok := tip.Commit(logpage.NewLogSequenceNumber(0), func() {
			completeA++
		})
		if !ok {
			t.Error("should commit")
		}

		if l := len(pages); l != 1 {
			t.Errorf("incorrect length: %v, expected: %v", l, 1)
		}
	}
	{
		tip := NewTip(pool, current)
		want := 8
		s := tip.Span(want)

		if l := s.Length(); l != want {
			t.Errorf("incorrect length: %v, expected: %v", l, want)
		}

		pages, ok := tip.Commit(logpage.NewLogSequenceNumber(0), func() {
			completeB++
		})
		if !ok {
			t.Error("should commit")
		}

		if l := len(pages); l != 1 {
			t.Errorf("incorrect length: %v, expected: %v", l, 1)
		}

		last = pages[len(pages)-1]
	}

	last.Complete()
	if completeA != 1 {
		t.Errorf("incorrect complete A count: %v, expected: %v", completeA, 1)
	}
	if completeB != 1 {
		t.Errorf("incorrect complete B count: %v, expected: %v", completeB, 1)
	}
}
