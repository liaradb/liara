package btree

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/collection/btree/page"
)

func TestKeyNode(t *testing.T) {
	t.Parallel()

	t.Run("should insert", func(t *testing.T) {
		t.Parallel()

		bp := page.New(newMockBuffer(256))
		kn := newKeyNode(bp)

		data := testKeyNodeInsertData(t, kn)
		testKeyNodeChildren(t, kn, data)
	})

	t.Run("should iterate in order", func(t *testing.T) {
		t.Parallel()

		bp := page.New(newMockBuffer(256))
		kn := newKeyNode(bp)

		data := testKeyNodeInsertData(t, kn)
		testKeyNodeChildren(t, kn, data)
	})

	t.Run("should search items", func(t *testing.T) {
		t.Parallel()

		bp := page.New(newMockBuffer(256))
		kn := newKeyNode(bp)

		data := testKeyNodeInsertData(t, kn)
		testKeyNodeSearch(t, kn, data)
	})
}

// Insert in mixed order
func testKeyNodeInsertData(t *testing.T, kn *KeyNode) []KeyEntry {
	data := []KeyEntry{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	}

	order := []int{0, 2, 1}
	for _, i := range order {
		e := data[i]
		if _, _, ok := kn.Insert(e.key, e.block); !ok {
			t.Error("should insert")
		}
	}

	return data
}

// Verify items are in order
func testKeyNodeChildren(t *testing.T, kn *KeyNode, data []KeyEntry) {
	result := make([]KeyEntry, 0, len(data))
	for c := range kn.Children() {
		result = append(result, c)
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}

// Verify items are all searchable
func testKeyNodeSearch(t *testing.T, kn *KeyNode, data []KeyEntry) {
	for _, e := range data {
		if block := kn.Search(e.key); block != e.block {
			t.Errorf("incorrect record id: %v, expected: %v", block, e.block)
		}
	}
}
