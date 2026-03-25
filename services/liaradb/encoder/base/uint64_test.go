package base

import "testing"

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

func TestUIn64__Remainder(t *testing.T) {
	t.Parallel()

	b := NewUint64(2)

	data := make([]byte, 16)
	data0, ok := b.WriteData(data)
	if !ok {
		t.Error("unable to write")
	}

	if l := len(data0); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	b0 := Uint64(0)
	data1, ok := b0.ReadData(data)
	if !ok {
		t.Error("unable to read")
	}

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
