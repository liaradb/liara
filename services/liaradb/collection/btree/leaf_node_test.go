package btree

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/collection/btree/page"
)

func TestLeafNode_Child(t *testing.T) {
	p := page.New(make([]byte, 256))
	ln := NewLeafNode(p)

	children := []LeafEntry{
		newLeafEntry(
			Key("abcde"),
			NewRecordID(1, 2)),
		newLeafEntry(
			Key("fghij"),
			NewRecordID(3, 4)),
	}

	if i, ok := ln.Append(children[0]); !ok {
		t.Error("should append")
	} else if i != 0 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	if i, ok := ln.Append(children[1]); !ok {
		t.Error("should append")
	} else if i != 1 {
		t.Errorf("incorrect index: %v, expected: %v", i, 1)
	}

	result := make([]LeafEntry, 0, len(children))
	for i := range len(children) {
		c, err := ln.Child(int16(i))
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, c)
	}

	if !slices.Equal(result, children) {
		t.Errorf("incorrect result: %v, expected: %v", result, children)
	}
}

func TestLeafNode_Children(t *testing.T) {
	p := page.New(make([]byte, 256))
	ln := NewLeafNode(p)

	children := []LeafEntry{
		newLeafEntry(
			Key("abcde"),
			NewRecordID(1, 2)),
		newLeafEntry(
			Key("fghij"),
			NewRecordID(3, 4)),
	}

	if i, ok := ln.Append(children[0]); !ok {
		t.Error("should append")
	} else if i != 0 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	if i, ok := ln.Append(children[1]); !ok {
		t.Error("should append")
	} else if i != 1 {
		t.Errorf("incorrect index: %v, expected: %v", i, 1)
	}

	result := make([]LeafEntry, 0, len(children))
	for c, err := range ln.Children() {
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, c)
	}

	if !slices.Equal(result, children) {
		t.Errorf("incorrect result: %v, expected: %v", result, children)
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

	if _, ok := ln.Insert(data[0]); !ok {
		t.Error("should insert")
	}

	if _, ok := ln.Insert(data[2]); !ok {
		t.Error("should insert")
	}

	if _, ok := ln.Insert(data[1]); !ok {
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
