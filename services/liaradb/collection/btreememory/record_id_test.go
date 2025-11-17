package btreememory

import "testing"

func TestRecordID_Defaults(t *testing.T) {
	id := RecordID{}

	if b := id.Block(); b != 0 {
		t.Errorf("incorrect block: %v, expected: %v", b, 0)
	}

	if p := id.Position(); p != 0 {
		t.Errorf("incorrect position: %v, expected: %v", p, 0)
	}
}

func TestRecordID_New(t *testing.T) {
	id := NewRecordID(1, 2)

	if b := id.Block(); b != 1 {
		t.Errorf("incorrect block: %v, expected: %v", b, 1)
	}

	if p := id.Position(); p != 2 {
		t.Errorf("incorrect position: %v, expected: %v", p, 2)
	}
}
