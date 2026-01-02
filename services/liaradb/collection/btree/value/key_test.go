package value

import "testing"

func TestKey(t *testing.T) {
	for message, c := range map[string]struct {
		skip         bool
		a            []string
		b            []string
		equal        bool
		greater      bool
		greaterEqual bool
		less         bool
		lessEqual    bool
	}{
		"should equal": {
			a:            []string{"a", "1"},
			b:            []string{"a", "1"},
			equal:        true,
			greaterEqual: true,
			lessEqual:    true,
		},
		"should be less for first": {
			a:         []string{"a", "1"},
			b:         []string{"b", "1"},
			less:      true,
			lessEqual: true,
		},
		"should be less for second": {
			a:         []string{"a", "1"},
			b:         []string{"a", "2"},
			less:      true,
			lessEqual: true,
		},
		"should be greater for first": {
			a:            []string{"b", "1"},
			b:            []string{"a", "1"},
			greater:      true,
			greaterEqual: true,
		},
		"should be greater for second": {
			a:            []string{"a", "2"},
			b:            []string{"a", "1"},
			greater:      true,
			greaterEqual: true,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			a := NewKey2([]byte(c.a[0]), []byte(c.a[1]))
			b := NewKey2([]byte(c.b[0]), []byte(c.b[1]))

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
