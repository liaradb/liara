package keynode

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
		kn := NewKeyNode(bp)

		data := testKeyNodeInsertData(t, kn)
		testKeyNodeChildren(t, kn, data)
	})

	t.Run("should iterate in order", func(t *testing.T) {
		t.Parallel()

		bp := page.New(newMockBuffer(256))
		kn := NewKeyNode(bp)

		data := testKeyNodeInsertData(t, kn)
		testKeyNodeChildren(t, kn, data)
	})

	t.Run("should search items", func(t *testing.T) {
		t.Parallel()

		bp := page.New(newMockBuffer(256))
		kn := NewKeyNode(bp)

		data := testKeyNodeInsertData(t, kn)
		testKeyNodeSearch(t, kn, data)
	})

	// TODO: This method is private
	t.Run("should search indexes", func(t *testing.T) {
		t.Parallel()

		bp := page.New(newMockBuffer(256))
		kn := NewKeyNode(bp)

		data := testKeyNodeInsertData(t, kn)
		for i, e := range data {
			index := kn.searchIndex(e.key)
			if index != int16(i) {
				t.Errorf("incorrect index: %v, expected: %v", index, int16(i))
			}
		}
		{
			before := int16(0)
			result := kn.searchIndex("a")
			if result != before {
				t.Errorf("incorrect before: %v, expected: %v", result, before)
			}
		}
		{
			after := int16(len(data))
			result := kn.searchIndex("e")
			if result != after {
				t.Errorf("incorrect after: %v, expected: %v", result, after)
			}
		}
	})

	t.Run("should search before items", func(t *testing.T) {
		t.Parallel()

		bp := page.New(newMockBuffer(256))
		kn := NewKeyNode(bp)

		data := testKeyNodeInsertData(t, kn)
		want := data[0].block
		result := kn.Search("a")
		if result != want {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}
	})

	t.Run("should search after items", func(t *testing.T) {
		t.Parallel()

		bp := page.New(newMockBuffer(256))
		kn := NewKeyNode(bp)

		data := testKeyNodeInsertData(t, kn)
		want := data[2].block
		result := kn.Search("e")
		if result != want {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}
	})
}

// Insert in mixed order
func testKeyNodeInsertData(t *testing.T, kn *KeyNode) []KeyEntry {
	data := []KeyEntry{
		{"b", 1},
		{"c", 2},
		{"d", 3},
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
