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

func TestUInt64__Remainder(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	b0 := NewUint64(2)
	if err := b0.Write(w); err != nil {
		t.Fatal(err)
	}

	b1 := Uint64(0)
	if err := b1.Read(r); err != nil {
		t.Fatal(err)
	}

	if v := b1.Value(); v != 2 {
		t.Errorf("incorrect value: %v, expected: %v", v, 2)
	}

	if s := b0.Size(); s != 8 {
		t.Errorf("incorrect size: %v, expected: %v", s, 8)
	}
}

func TestUint64_Read_Write(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	b := NewUint64(2)
	if err := b.Write(w); err != nil {
		t.Fatal(err)
	}

	var b0 Uint64
	if err := b0.Read(r); err != nil {
		t.Fatal(err)
	}

	if b != b0 {
		t.Errorf("incorrect value: %v, expected: %v", b, b0)
	}
}

func TestUint64_Signed(t *testing.T) {
	t.Parallel()

	b := NewInt64(-10)

	if v := b.Signed(); v != -10 {
		t.Errorf("incorrect value: %v, expected: %v", v, -10)
	}
}
