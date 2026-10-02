package base

import (
	"io"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/encoder/buffer"
)

func TestBytes(t *testing.T) {
	t.Parallel()

	size := 32
	n := NewBytes([]byte("name"))
	data := make([]byte, size)
	_ = n.WriteData(data, size)

	var r Bytes
	data0 := r.ReadData(data, size)

	if l := len(data0); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}

	if !r.Compare(n) {
		t.Error("should be equal")
	}

	if r.String() != n.String() {
		t.Errorf("incorrect result: %v, expected: %v", r.String(), n.String())
	}

	if b := r.Value(); !slices.Equal(b, []byte("name")) {
		t.Errorf("incorrect bytes: %v, expected: %v", b, []byte("name"))
	}

	if l := r.Length(); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if s := r.Size(); s != 4+4 {
		t.Errorf("incorrect size: %v, expected: %v", s, 4+4)
	}
}

func TestBytes_Read_Write(t *testing.T) {
	t.Parallel()

	buf := buffer.New(8)

	b := NewBytes([]byte{1, 2, 3, 4})
	if err := b.Write(buf); err != nil {
		t.Fatal(err)
	}

	if _, err := buf.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	var b0 Bytes
	if err := b0.Read(buf); err != nil {
		t.Fatal(err)
	}

	if !b0.Compare(b) {
		t.Errorf("incorrect value: %v, expected: %v", b0.Value(), b.Value())
	}
}
