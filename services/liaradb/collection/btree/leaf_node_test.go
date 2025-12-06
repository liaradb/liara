package btree

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/collection/btree/page"
)

func TestLeafNode_Child(t *testing.T) {
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
		c, err := ln.Child(int16(i))
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, c)
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}

func TestLeafNode_Children(t *testing.T) {
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
	for c, err := range ln.Children() {
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, c)
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}

func TestLeafNode_Insert(t *testing.T) {
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

	if _, ok := ln.Insert(data[0].key, data[0].recordID); !ok {
		t.Error("should insert")
	}

	if _, ok := ln.Insert(data[2].key, data[2].recordID); !ok {
		t.Error("should insert")
	}

	if _, ok := ln.Insert(data[1].key, data[1].recordID); !ok {
		t.Error("should insert")
	}

	result := make([]LeafEntry, 0, len(data))

	for c, err := range ln.Children() {
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, c)
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}
