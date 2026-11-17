package graph

import (
	"testing"

	"github.com/liaradb/liaradb/storage/link"
)

func TestPoint(t *testing.T) {
	t.Parallel()

	var (
		id       = link.NewRecordID(link.BlockID(1), link.SlotID(2))
		edge     = link.NewRecordID(link.BlockID(3), link.SlotID(4))
		property = link.NewRecordID(link.BlockID(7), link.SlotID(8))
	)

	p := NewPoint(id, edge, property)

	if v := p.ID(); v != id {
		t.Errorf("incorrect id: %v, expected: %v", v, id)
	}

	if v := p.Edge(); v != edge {
		t.Errorf("incorrect edge: %v, expected: %v", v, edge)
	}

	if v := p.Property(); v != property {
		t.Errorf("incorrect property: %v, expected: %v", v, property)
	}

	if v := p.Size(); v != Pointsize {
		t.Errorf("incorrect size: %v, expected: %v", v, Pointsize)
	}
}
