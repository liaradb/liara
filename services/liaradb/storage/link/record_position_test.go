package link

import (
	"io"
	"testing"
)

func TestRecordPosition(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var p RecordPosition = 123
	if err := p.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := p.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var p2 RecordPosition
	if err := p2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if p != p2 {
		t.Errorf("incorrect value: %v, expected: %v", p2, p)
	}

	if s := p.String(); s != "123" {
		t.Errorf("incorrect string: %v, expected: %v", s, "123")
	}
}

func TestRecordPosition_ReadDataWriteData(t *testing.T) {
	t.Parallel()

	rp := RecordPosition(1)

	data := make([]byte, 5)
	data0, ok := rp.WriteData(data)
	if !ok {
		t.Error("unable to write")
	}

	if l := len(data0); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	rp0 := RecordPosition(0)
	data1, ok := rp0.ReadData(data)
	if !ok {
		t.Error("unable to read")
	}

	if l := len(data1); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if rp0 != rp {
		t.Errorf("incorrect value: %v, expected: %v", rp0, rp)
	}
}
