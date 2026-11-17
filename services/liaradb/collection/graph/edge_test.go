package graph

import (
	"testing"

	"github.com/liaradb/liaradb/storage/link"
)

func TestEdge(t *testing.T) {
	t.Parallel()

	var (
		id        = link.NewRecordID(link.BlockID(1), link.SlotID(2))
		start     = link.NewRecordID(link.BlockID(3), link.SlotID(4))
		end       = link.NewRecordID(link.BlockID(5), link.SlotID(6))
		direction = ForwardDirection
		property  = link.NewRecordID(link.BlockID(7), link.SlotID(8))
	)

	e := NewEdge(id, start, end, direction, property)

	if v := e.ID(); v != id {
		t.Errorf("incorrect id: %v, expected: %v", v, id)
	}

	if v := e.Start(); v != start {
		t.Errorf("incorrect start: %v, expected: %v", v, start)
	}

	if v := e.End(); v != end {
		t.Errorf("incorrect end: %v, expected: %v", v, end)
	}

	if v := e.Direction(); v != direction {
		t.Errorf("incorrect direction: %v, expected: %v", v, direction)
	}

	if v := e.Property(); v != property {
		t.Errorf("incorrect property: %v, expected: %v", v, property)
	}

	if v := e.Size(); v != EdgeSize {
		t.Errorf("incorrect size: %v, expected: %v", v, EdgeSize)
	}
}
