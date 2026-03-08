package value

import (
	"slices"
	"testing"
)

func TestData(t *testing.T) {
	value := "value"
	data := NewData([]byte(value))

	if s := data.String(); s != value {
		t.Errorf("incorrect string: %v, expected: %v", s, value)
	}

	if v := data.Value(); !slices.Equal(v, []byte(value)) {
		t.Errorf("incorrect value: %v, expected: %v", v, []byte(value))
	}

	if s := data.Size(); s != len(value) {
		t.Errorf("incorrect size: %v, expected: %v", s, len(value))
	}
}

func TestData_Compare(t *testing.T) {
	for message, c := range map[string]struct {
		skip  bool
		a     Data
		b     Data
		equal bool
	}{
		"should equal zero": {
			a:     Data{},
			b:     Data{},
			equal: true,
		},
		"should not equal zero and data": {
			a:     Data{},
			b:     NewData([]byte{1, 2, 3}),
			equal: false,
		},
		"should equal with data": {
			a:     NewData([]byte{1, 2, 3}),
			b:     NewData([]byte{1, 2, 3}),
			equal: true,
		},
		"should not equal with different data": {
			a:     NewData([]byte{1, 2, 3}),
			b:     NewData([]byte{3, 2, 1}),
			equal: false,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			if c.a.Compare(&c.b) != c.equal {
				if c.equal {
					t.Error("should equal")
				} else {
					t.Error("should not equal")
				}
			}
		})
	}
}
