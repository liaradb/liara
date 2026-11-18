package base

import "testing"

func TestUint32_String(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip  bool
		value Uint32
		want  string
	}{
		"should handle 0": {
			value: NewUint32(0),
			want:  "00000000",
		},
		"should handle 1": {
			value: NewUint32(1),
			want:  "00000001",
		},
		"should handle 2": {
			value: NewUint32(2),
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

func TestUIn32__Remainder(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	b0 := NewUint32(2)
	if err := b0.Write(w); err != nil {
		t.Fatal(err)
	}

	b1 := Uint32(0)
	if err := b1.Read(r); err != nil {
		t.Fatal(err)
	}

	if v := b1.Value(); v != 2 {
		t.Errorf("incorrect value: %v, expected: %v", v, 2)
	}

	if s := b0.Size(); s != 4 {
		t.Errorf("incorrect size: %v, expected: %v", s, 4)
	}
}

func TestUint32_Read_Write(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	b := NewUint32(2)
	if err := b.Write(w); err != nil {
		t.Fatal(err)
	}

	var b0 Uint32
	if err := b0.Read(r); err != nil {
		t.Fatal(err)
	}

	if b != b0 {
		t.Errorf("incorrect value: %v, expected: %v", b, b0)
	}
}

func TestUint32_Signed(t *testing.T) {
	t.Parallel()

	b := NewInt32(-10)

	if v := b.Signed(); v != -10 {
		t.Errorf("incorrect value: %v, expected: %v", v, -10)
	}
}
