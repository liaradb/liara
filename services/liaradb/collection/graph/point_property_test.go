package graph

import (
	"testing"

	"github.com/liaradb/liaradb/storage/link"
)

func TestPointProperty(t *testing.T) {
	t.Parallel()

	var (
		id   = link.NewRecordID(link.BlockID(1), link.SlotID(2))
		next = link.NewRecordID(link.BlockID(3), link.SlotID(4))
	)

	pp := NewPointProperty(id, next)

	if v := pp.ID(); v != id {
		t.Errorf("incorrect id: %v, expected: %v", v, id)
	}

	if v := pp.Next(); v != next {
		t.Errorf("incorrect next: %v, expected: %v", v, next)
	}

	if v := pp.Size(); v != PointPropertySize {
		t.Errorf("incorrect size: %v, expected: %v", v, PointPropertySize)
	}
}
