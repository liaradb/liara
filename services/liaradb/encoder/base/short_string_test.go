package base

import (
	"io"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/encoder/buffer"
)

func TestShortString(t *testing.T) {
	t.Parallel()

	want := "name"
	n := NewShortString(want)

	if n.String() != want {
		t.Errorf("incorrect string: %v, expected: %v", n.String(), want)
	}

	if b := n.Bytes(); !slices.Equal(b, []byte(want)) {
		t.Errorf("incorrect bytes: %v, expected: %v", b, []byte(want))
	}

	if l := n.Length(); l != len(want) {
		t.Errorf("incorrect length: %v, expected: %v", l, len(want))
	}

	if s := n.Size(); s != 256 {
		t.Errorf("incorrect size: %v, expected: %v", s, 256)
	}
}

func TestShortString_Read_Write(t *testing.T) {
	t.Parallel()

	buf := buffer.New(8)

	b := NewShortString("name")
	if err := b.Write(buf); err != nil {
		t.Fatal(err)
	}

	if _, err := buf.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	var b0 ShortString
	if err := b0.Read(buf); err != nil {
		t.Fatal(err)
	}

	if b != b0 {
		t.Errorf("incorrect value: %v, expected: %v", b, b0)
	}
}
