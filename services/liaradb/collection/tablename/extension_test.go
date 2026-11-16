package tablename

import "testing"

func TestExtension_FileName(t *testing.T) {
	t.Parallel()

	t.Run("should create collection name", func(t *testing.T) {
		fn := CollectionExtension.FileName("name")
		want := "name.col"
		if n := fn.String(); n != "name.col" {
			t.Errorf("incorrect file name: %v, expected: %v", n, want)
		}
	})

	t.Run("should create B+Tree name", func(t *testing.T) {
		fn := BTreeExtension.FileName("name")
		want := "name.bt"
		if n := fn.String(); n != "name.bt" {
			t.Errorf("incorrect file name: %v, expected: %v", n, want)
		}
	})
}

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
