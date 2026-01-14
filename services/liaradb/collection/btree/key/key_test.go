package key

import "testing"

func TestKey(t *testing.T) {
	for message, c := range map[string]struct {
		skip         bool
		aString      string
		aInt         int64
		bString      string
		bInt         int64
		equal        bool
		greater      bool
		greaterEqual bool
		less         bool
		lessEqual    bool
	}{
		"should equal": {
			aString:      "a",
			aInt:         1,
			bString:      "a",
			bInt:         1,
			equal:        true,
			greaterEqual: true,
			lessEqual:    true,
		},
		"should be less for first": {
			aString:   "a",
			aInt:      1,
			bString:   "b",
			bInt:      1,
			less:      true,
			lessEqual: true,
		},
		"should be less for second": {
			aString:   "a",
			aInt:      1,
			bString:   "a",
			bInt:      2,
			less:      true,
			lessEqual: true,
		},
		"should be greater for first": {
			aString:      "b",
			aInt:         1,
			bString:      "a",
			bInt:         1,
			greater:      true,
			greaterEqual: true,
		},
		"should be greater for second": {
			aString:      "a",
			aInt:         2,
			bString:      "a",
			bInt:         1,
			greater:      true,
			greaterEqual: true,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			a := NewKey2([]byte(c.aString), c.aInt)
			b := NewKey2([]byte(c.bString), c.bInt)

			if c.equal {
				if !a.Equal(b) {
					t.Error("should be equal")
				}
			} else {
				if a.Equal(b) {
					t.Error("should not be equal")
				}
			}

			if c.greater {
				if !a.Greater(b) {
					t.Error("should be greater")
				}
			} else {
				if a.Greater(b) {
					t.Error("should not be greater")
				}
			}

			if c.greaterEqual {
				if !a.GreaterEqual(b) {
					t.Error("should be greater than or equal")
				}
			} else {
				if a.GreaterEqual(b) {
					t.Error("should not be greater than or equal")
				}
			}

			if c.less {
				if !a.Less(b) {
					t.Error("should be less")
				}
			} else {
				if a.Less(b) {
					t.Error("should not be less")
				}
			}

			if c.lessEqual {
				if !a.LessEqual(b) {
					t.Error("should be less than or equal")
				}
			} else {
				if a.LessEqual(b) {
					t.Error("should not be less than or equal")
				}
			}
		})
	}
}

func TestKey_Size(t *testing.T) {
	for message, c := range map[string]struct {
		skip bool
		k    Key
		size int
	}{
		"should handle default": {
			k:    NewKey2(nil, 0),
			size: 8,
		},
		"should handle values": {
			k:    NewKey([]byte("abcdef")),
			size: 14,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			if s := c.k.Size(); s != c.size {
				t.Errorf("%v: incorrect size: %v, expected: %v", message, s, c.size)
			}
		})
	}
}

func TestKey_String(t *testing.T) {
	a := "a"
	b := 2
	k := NewKey2([]byte(a), int64(b))

	want := "a2"

	if s := k.String(); s != want {
		t.Errorf("incorrect string: %v, expected: %v", s, want)
	}
}
