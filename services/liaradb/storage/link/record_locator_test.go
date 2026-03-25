package link

import "testing"

func TestRecordLocator_Defaults(t *testing.T) {
	t.Parallel()

	id := RecordLocator{}

	if b := id.Block(); b != 0 {
		t.Errorf("incorrect block: %v, expected: %v", b, 0)
	}

	if p := id.Position(); p != 0 {
		t.Errorf("incorrect position: %v, expected: %v", p, 0)
	}

	if s := id.Size(); s != 10 {
		t.Errorf("incorrect size: %v, expected: %v", s, 10)
	}
}

func TestRecordLocator_New(t *testing.T) {
	t.Parallel()

	id := NewRecordLocator(1, 2)

	if b := id.Block(); b != 1 {
		t.Errorf("incorrect block: %v, expected: %v", b, 1)
	}

	if p := id.Position(); p != 2 {
		t.Errorf("incorrect position: %v, expected: %v", p, 2)
	}
}

func TestRecordLocator_WriteRead(t *testing.T) {
	t.Parallel()

	id := NewRecordLocator(1, 2)

	data := make([]byte, RecordLocatorSize)
	id.Write(data)

	id0 := RecordLocator{}
	id0.Read(data)

	if b := id0.Block(); b != 1 {
		t.Errorf("incorrect block: %v, expected: %v", b, 1)
	}

	if p := id0.Position(); p != 2 {
		t.Errorf("incorrect position: %v, expected: %v", p, 2)
	}
}
