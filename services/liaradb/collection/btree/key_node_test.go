package btree

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/collection/btree/page"
)

func TestKeyNode(t *testing.T) {
	bp := page.New(make([]byte, 256))
	kn := newKeyNode(bp)

	data := []KeyEntry{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	}

	// Insert in mixed order
	order := []int{0, 2, 1}
	for _, i := range order {
		e := data[i]
		if _, ok := kn.Insert(e.key, e.block); !ok {
			t.Error("should insert")
		}
	}

	// Verify items are in order
	result := make([]KeyEntry, 0, len(data))
	for c, err := range kn.Children() {
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, c)
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}

	// Verify items are all searchable
	for _, e := range data {
		if block, err := kn.Search(e.key); err != nil {
			t.Error(err)
		} else if block != e.block {
			t.Errorf("incorrect record id: %v, expected: %v", block, e.block)
		}
	}
}
