package graph

import (
	"testing"

	"github.com/liaradb/liaradb/storage/link"
)

func TestEdgeProperty(t *testing.T) {
	t.Parallel()

	var (
		id   = link.NewRecordID(link.BlockID(1), link.SlotID(2))
		next = link.NewRecordID(link.BlockID(3), link.SlotID(4))
	)

	ep := NewEdgeProperty(id, next)

	if v := ep.ID(); v != id {
		t.Errorf("incorrect id: %v, expected: %v", v, id)
	}

	if v := ep.Next(); v != next {
		t.Errorf("incorrect next: %v, expected: %v", v, next)
	}

	if v := ep.Size(); v != EdgePropertySize {
		t.Errorf("incorrect size: %v, expected: %v", v, EdgePropertySize)
	}
}
