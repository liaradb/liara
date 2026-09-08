package link

import (
	"io"
	"testing"
)

func TestSlotID(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var id0 SlotID = 123
	if err := id0.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := id0.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var id1 SlotID
	if err := id1.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if id0 != id1 {
		t.Errorf("incorrect value: %v, expected: %v", id1, id0)
	}

	if s := id0.String(); s != "123" {
		t.Errorf("incorrect string: %v, expected: %v", s, "123")
	}
}

func TestSlotID_ReadDataWriteData(t *testing.T) {
	t.Parallel()

	id := SlotID(1)

	data := make([]byte, 6)
	data0, ok := id.WriteData(data)
	if !ok {
		t.Error("unable to write")
	}

	if l := len(data0); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	id0 := SlotID(0)
	data1, ok := id0.ReadData(data)
	if !ok {
		t.Error("unable to read")
	}

	if l := len(data1); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if id0 != id {
		t.Errorf("incorrect value: %v, expected: %v", id0, id)
	}
}
