package base

import (
	"io"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/encoder/buffer"
)

func TestString(t *testing.T) {
	t.Parallel()

	size := 32
	n := String("name")
	data := make([]byte, size)
	_ = n.WriteData(data, size)

	var r String
	data0 := r.ReadData(data, size)

	if l := len(data0); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}

	if r.String() != n.String() {
		t.Errorf("incorrect result: %v, expected: %v", r.String(), n.String())
	}

	if b := r.Bytes(); !slices.Equal(b, []byte("name")) {
		t.Errorf("incorrect bytes: %v, expected: %v", b, []byte("name"))
	}

	if l := r.Length(); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if s := r.Size(); s != 4+4 {
		t.Errorf("incorrect size: %v, expected: %v", s, 4+4)
	}
}

func TestString_Read_Write(t *testing.T) {
	t.Parallel()

	buf := buffer.New(8)

	b := String("name")
	if err := b.Write(buf); err != nil {
		t.Fatal(err)
	}

	if _, err := buf.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	var b0 String
	if err := b0.Read(buf); err != nil {
		t.Fatal(err)
	}

	if b != b0 {
		t.Errorf("incorrect value: %v, expected: %v", b, b0)
	}
}
