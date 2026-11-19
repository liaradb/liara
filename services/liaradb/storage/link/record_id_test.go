package link

import "testing"

func TestRecordID_Defaults(t *testing.T) {
	t.Parallel()

	id := RecordID{}

	if b := id.BlockID(); b != 0 {
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

	if b := id.BlockID(); b != 1 {
		t.Errorf("incorrect block: %v, expected: %v", b, 1)
	}

	if s := id.SlotID(); s != 2 {
		t.Errorf("incorrect slot id: %v, expected: %v", s, 2)
	}
}

func TestRecordID_ReadWrite(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	rid0 := NewRecordID(2, 3)
	if err := rid0.Write(w); err != nil {
		t.Fatal(err)
	}

	var rid1 RecordID
	if err := rid1.Read(r); err != nil {
		t.Fatal(err)
	}

	if rid1 != rid0 {
		t.Errorf("incorrect result: %v, expected: %v", rid1, rid0)
	}
}

func TestRecordID_WriteDataReadData(t *testing.T) {
	t.Parallel()

	id := NewRecordID(1, 2)

	data := make([]byte, RecordIDSize)
	id.WriteData(data)

	id0 := RecordID{}
	id0.ReadData(data)

	if b := id0.BlockID(); b != 1 {
		t.Errorf("incorrect block: %v, expected: %v", b, 1)
	}

	if s := id0.SlotID(); s != 2 {
		t.Errorf("incorrect slot id: %v, expected: %v", s, 2)
	}
}
