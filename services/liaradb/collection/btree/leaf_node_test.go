package btree

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/collection/btree/page"
)

func TestLeafNode_Child(t *testing.T) {
	t.Parallel()

	p := page.New(make([]byte, 256))
	ln := NewLeafNode(p)

	data := []LeafEntry{
		newLeafEntry(
			Key("abcde"),
			NewRecordID(1, 2)),
		newLeafEntry(
			Key("fghij"),
			NewRecordID(3, 4)),
	}

	if i, ok := ln.Append(data[0].key, data[0].recordID); !ok {
		t.Error("should append")
	} else if i != 0 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	if i, ok := ln.Append(data[1].key, data[1].recordID); !ok {
		t.Error("should append")
	} else if i != 1 {
		t.Errorf("incorrect index: %v, expected: %v", i, 1)
	}

	result := make([]LeafEntry, 0, len(data))
	for i := range len(data) {
		c, ok := ln.Child(int16(i))
		if !ok {
			t.Fatal("should get child")
		}

		result = append(result, c)
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}

func TestLeafNode_Children(t *testing.T) {
	t.Parallel()

	p := page.New(make([]byte, 256))
	ln := NewLeafNode(p)

	data := []LeafEntry{
		newLeafEntry(
			Key("abcde"),
			NewRecordID(1, 2)),
		newLeafEntry(
			Key("fghij"),
			NewRecordID(3, 4)),
	}

	if i, ok := ln.Append(data[0].key, data[0].recordID); !ok {
		t.Error("should append")
	} else if i != 0 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	if i, ok := ln.Append(data[1].key, data[1].recordID); !ok {
		t.Error("should append")
	} else if i != 1 {
		t.Errorf("incorrect index: %v, expected: %v", i, 1)
	}

	result := make([]LeafEntry, 0, len(data))
	for c := range ln.Children() {
		result = append(result, c)
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}

func TestLeafNode_Insert(t *testing.T) {
	t.Parallel()

	bp := page.New(make([]byte, 256))
	ln := NewLeafNode(bp)

	data := []LeafEntry{
		newLeafEntry(
			Key("a"),
			NewRecordID(1, 2)),
		newLeafEntry(
			Key("b"),
			NewRecordID(3, 4)),
		newLeafEntry(
			Key("c"),
			NewRecordID(5, 6)),
	}

	// Insert in mixed order
	order := []int{0, 2, 1}
	for _, i := range order {
		e := data[i]
		if _, _, ok := ln.Insert(e.key, e.recordID); !ok {
			t.Error("should insert")
		}
	}

	// Verify items are in order
	result := make([]LeafEntry, 0, len(data))
	for c := range ln.Children() {
		result = append(result, c)
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}

	// Verify items are all searchable
	for _, e := range data {
		if rid, ok := ln.Search(e.key); !ok {
			t.Error("should find record id")
		} else if rid != e.recordID {
			t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
		}
	}
}
