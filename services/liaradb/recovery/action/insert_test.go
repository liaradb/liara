package action

import (
	"slices"
	"testing"
)

func TestInsert(t *testing.T) {
	cid := CollectionID("a")
	iid := ItemID("b")
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	i := NewInsert(cid, iid, data)

	if v := i.CollectionID(); v != cid {
		t.Errorf("incorrect collection id: %v, expected: %v", v, cid)
	}

	if v := i.ItemID(); v != iid {
		t.Errorf("incorrect item id: %v, expected: %v", v, iid)
	}

	if v := i.Data(); !slices.Equal(v, data) {
		t.Errorf("incorrect data: %v, expected: %v", v, data)
	}
}
