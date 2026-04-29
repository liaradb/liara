package page

import (
	"testing"

	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

func TestPool(t *testing.T) {
	pp := NewPool(1024)

	p0 := pp.Get(action.PageID(1), action.TimeLineID(2), record.NewLength(3))
	pp.Return(p0)

	if i := p0.ID(); i != action.PageID(1) {
		t.Errorf("incorrect id: %v, expected: %v", i, action.PageID(1))
	}

	p1 := pp.Get(action.PageID(10), action.TimeLineID(20), record.NewLength(30))
	pp.Return(p1)

	if i := p1.ID(); i != action.PageID(10) {
		t.Errorf("incorrect id: %v, expected: %v", i, action.PageID(10))
	}
}
