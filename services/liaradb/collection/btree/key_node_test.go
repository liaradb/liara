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

	if _, ok := kn.Insert(data[0].key, data[0].block); !ok {
		t.Error("should insert")
	}

	if _, ok := kn.Insert(data[2].key, data[2].block); !ok {
		t.Error("should insert")
	}

	if _, ok := kn.Insert(data[1].key, data[1].block); !ok {
		t.Error("should insert")
	}

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
}
