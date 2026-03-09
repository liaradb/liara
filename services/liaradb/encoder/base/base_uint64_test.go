package base

import "testing"

func TestBaseUint64_String(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip  bool
		value BaseUint64
		want  string
	}{
		"should handle 0": {
			value: NewBaseUint64(0),
			want:  "0000000000000000",
		},
		"should handle 1": {
			value: NewBaseUint64(1),
			want:  "0000000000000001",
		},
		"should handle 2": {
			value: NewBaseUint64(2),
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

func TestBaseUIn64__Remainder(t *testing.T) {
	t.Parallel()

	b := NewBaseUint64(2)

	data := make([]byte, 16)
	data0 := b.WriteData(data)

	if l := len(data0); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	b0 := BaseUint64(0)
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
