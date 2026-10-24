package tablename

import "testing"

func TestIsIndex(t *testing.T) {
	for message, c := range map[string]struct {
		skip   bool
		name   string
		result bool
	}{
		"should handle default": {},
		"should handle index": {
			name:   "file.idx",
			result: true,
		},
		"should handle other extension": {
			name: "file.col",
		},
		"should handle no extension": {
			name: "file",
		},
		"should handle only extension": {
			name: "idx",
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			t.Run(message, func(t *testing.T) {
				if result := IsIndex(c.name); result != c.result {
					t.Errorf("incorrect result: %v, expected: %v", result, c.result)
				}
			})
		})
	}
}
