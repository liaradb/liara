package tablename

import "testing"

func TestIsBTree(t *testing.T) {
	for message, c := range map[string]struct {
		skip   bool
		name   string
		result bool
	}{
		"should handle default": {},
		"should handle bt": {
			name:   "file.bt",
			result: true,
		},
		"should handle other extension": {
			name: "file.col",
		},
		"should handle no extension": {
			name: "file",
		},
		"should handle only extension": {
			name: "bt",
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			t.Run(message, func(t *testing.T) {
				if result := IsBTree(c.name); result != c.result {
					t.Errorf("incorrect result: %v, expected: %v", result, c.result)
				}
			})
		})
	}
}

func TestIsCollection(t *testing.T) {
	for message, c := range map[string]struct {
		skip   bool
		name   string
		result bool
	}{
		"should handle default": {},
		"should handle index": {
			name:   "file.col",
			result: true,
		},
		"should handle other extension": {
			name: "file.bt",
		},
		"should handle no extension": {
			name: "file",
		},
		"should handle only extension": {
			name: "col",
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			t.Run(message, func(t *testing.T) {
				if result := IsCollection(c.name); result != c.result {
					t.Errorf("incorrect result: %v, expected: %v", result, c.result)
				}
			})
		})
	}
}
