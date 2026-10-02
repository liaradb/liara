package base

import (
	"io"
	"testing"

	"github.com/liaradb/liaradb/encoder/buffer"
)

func TestUint64_String(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip  bool
		value Uint64
		want  string
	}{
		"should handle 0": {
			value: NewUint64(0),
			want:  "0000000000000000",
		},
		"should handle 1": {
			value: NewUint64(1),
			want:  "0000000000000001",
		},
		"should handle 2": {
			value: NewUint64(2),
			want:  "0000000000000002",
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			if s := c.value.String(); s != c.want {
				t.Errorf("%v: incorrect string: %v, expected: %v", message, s, c.want)
			}
		})
	}
}

func TestUInt64__Remainder(t *testing.T) {
	t.Parallel()

	b := NewUint64(2)

	data := make([]byte, 16)
	data0 := b.WriteData(data)

	if l := len(data0); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	b0 := Uint64(0)
	data1 := b0.ReadData(data)

	if l := len(data1); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	if v := b0.Value(); v != 2 {
		t.Errorf("incorrect value: %v, expected: %v", v, 2)
	}

	if s := b.Size(); s != 8 {
		t.Errorf("incorrect size: %v, expected: %v", s, 8)
	}
}

func TestUint64_Read_Write(t *testing.T) {
	t.Parallel()

	buf := buffer.New(8)

	b := NewUint64(2)
	if err := b.Write(buf); err != nil {
		t.Fatal(err)
	}

	if _, err := buf.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	var b0 Uint64
	if err := b0.Read(buf); err != nil {
		t.Fatal(err)
	}

	if b != b0 {
		t.Errorf("incorrect value: %v, expected: %v", b, b0)
	}
}

func TestUInt64(t *testing.T) {
	t.Parallel()
	b0 := NewUint64(2)

	var b1 Uint64
	_ = b1.ReadData(b0.Bytes())

	if v := b1.Value(); v != 2 {
		t.Errorf("incorrect value: %v, expected: %v", v, 2)
	}
}

func TestUint64_Signed(t *testing.T) {
	t.Parallel()

	b := NewInt64(-10)

	if v := b.Signed(); v != -10 {
		t.Errorf("incorrect value: %v, expected: %v", v, -10)
	}
}
