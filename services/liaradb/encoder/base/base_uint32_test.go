package base

import "testing"

func TestBaseUint32_String(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip  bool
		value BaseUint32
		want  string
	}{
		"should handle 0": {
			value: NewBaseUint32(0),
			want:  "00000000",
		},
		"should handle 1": {
			value: NewBaseUint32(1),
			want:  "00000001",
		},
		"should handle 2": {
			value: NewBaseUint32(2),
			want:  "00000002",
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

func TestBaseUIn32__Remainder(t *testing.T) {
	t.Parallel()

	b := NewBaseUint32(2)

	data := make([]byte, 8)
	data0 := b.WriteData(data)

	if l := len(data0); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	b0 := BaseUint32(0)
	data1 := b0.ReadData(data)

	if l := len(data1); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if v := b0.Value(); v != 2 {
		t.Errorf("incorrect value: %v, expected: %v", v, 2)
	}

	if s := b.Size(); s != 4 {
		t.Errorf("incorrect size: %v, expected: %v", s, 4)
	}
}
