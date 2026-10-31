package link

import "testing"

func TestRecordID_Defaults(t *testing.T) {
	t.Parallel()

	id := RecordID{}

	if b := id.Block(); b != 0 {
		t.Errorf("incorrect block: %v, expected: %v", b, 0)
	}

	if s := id.SlotID(); s != 0 {
		t.Errorf("incorrect slot: %v, expected: %v", s, 0)
	}

	if s := id.Size(); s != 10 {
		t.Errorf("incorrect size: %v, expected: %v", s, 10)
	}
}

func TestRecordID_New(t *testing.T) {
	t.Parallel()

	id := NewRecordID(1, 2)

	if b := id.Block(); b != 1 {
		t.Errorf("incorrect block: %v, expected: %v", b, 1)
	}

	if s := id.SlotID(); s != 2 {
		t.Errorf("incorrect slot id: %v, expected: %v", s, 2)
	}
}

func TestRecordID_WriteRead(t *testing.T) {
	t.Parallel()

	id := NewRecordID(1, 2)

	data := make([]byte, RecordIDSize)
	id.WriteData(data)

	id0 := RecordID{}
	id0.ReadData(data)

	if b := id0.Block(); b != 1 {
		t.Errorf("incorrect block: %v, expected: %v", b, 1)
	}

	if s := id0.SlotID(); s != 2 {
		t.Errorf("incorrect slot id: %v, expected: %v", s, 2)
	}
}
