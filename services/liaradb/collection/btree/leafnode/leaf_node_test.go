package leafnode

import (
	"fmt"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/collection/btree/node"
	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/storagetesting"
)

func TestLeafNode_Child(t *testing.T) {
	t.Parallel()

	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)

	p := node.New(b)
	ln := New(p)

	data := []leafEntry{
		newLeafEntry(
			value.Key("abcde"),
			value.NewRecordID(1, 2)),
		newLeafEntry(
			value.Key("fghij"),
			value.NewRecordID(3, 4)),
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

	result := make([]leafEntry, 0, len(data))
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

	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)

	p := node.New(b)
	ln := New(p)

	data := []leafEntry{
		newLeafEntry(
			value.Key("abcde"),
			value.NewRecordID(1, 2)),
		newLeafEntry(
			value.Key("fghij"),
			value.NewRecordID(3, 4)),
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

	result := make([]leafEntry, 0, len(data))
	for key, rid := range ln.Children() {
		result = append(result, newLeafEntry(key, rid))
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}

func TestLeafNode_Insert(t *testing.T) {
	t.Parallel()

	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)

	bp := node.New(b)
	ln := New(bp)

	data := []leafEntry{
		newLeafEntry(
			value.Key("a"),
			value.NewRecordID(1, 2)),
		newLeafEntry(
			value.Key("b"),
			value.NewRecordID(3, 4)),
		newLeafEntry(
			value.Key("c"),
			value.NewRecordID(5, 6)),
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
	result := make([]leafEntry, 0, len(data))
	for key, rid := range ln.Children() {
		result = append(result, newLeafEntry(key, rid))
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}

	// Verify record ids
	{
		want := make([]value.RecordLocator, 0, len(data))
		for _, e := range data {
			want = append(want, e.RecordID())
		}

		result := make([]value.RecordLocator, 0, len(want))
		for c := range ln.RecordIDs() {
			result = append(result, c)
		}

		if !slices.Equal(result, want) {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}
	}

	// Verify items are all searchable
	for _, e := range data {
		if rid, ok := ln.Search(e.key); !ok {
			t.Error("should find record id")
		} else if rid != e.recordID {
			t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
		}
	}

	// Verify search range
	for i, e := range data {
		want := make([]value.RecordLocator, 0, len(data[i:]))
		for _, e := range data[i:] {
			want = append(want, e.RecordID())
		}

		result := make([]value.RecordLocator, 0, len(want))
		for c := range ln.SearchRange(e.Key()) {
			result = append(result, c)
		}

		if !slices.Equal(result, want) {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}
	}
}

func TestLeafNode_Insert__Split(t *testing.T) {
	t.Parallel()

	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)

	bp := node.New(b)
	ln := New(bp)

	data := []leafEntry{
		newLeafEntry(
			value.Key("a"),
			value.NewRecordID(1, 2)),
		newLeafEntry(
			value.Key("b"),
			value.NewRecordID(3, 4)),
		newLeafEntry(
			value.Key("c"),
			value.NewRecordID(5, 6)),
	}

	// Insert in mixed order
	order := []int{0, 2, 1}
	for _, i := range order {
		e := data[i]
		if a, b, ok := ln.Insert(e.key, e.recordID); !ok {
			fmt.Println("a...")
			for key, rid := range a {
				fmt.Println(key, rid)
			}
			fmt.Println("b...")
			for key, rid := range b {
				fmt.Println(key, rid)
			}
		}
	}

	// // Verify items are in order
	// result := make([]LeafEntry, 0, len(data))
	// for c := range ln.Children() {
	// 	result = append(result, c)
	// }

	// if !slices.Equal(result, data) {
	// 	t.Errorf("incorrect result: %v, expected: %v", result, data)
	// }

	// // Verify items are all searchable
	// for _, e := range data {
	// 	if rid, ok := ln.Search(e.key); !ok {
	// 		t.Error("should find record id")
	// 	} else if rid != e.recordID {
	// 		t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
	// 	}
	// }
}

func createBuffer(t *testing.T, s *storage.Storage) *storage.Buffer {
	b, err := s.Request(t.Context(), storage.BlockID{})
	if err != nil {
		t.Fatal(err)
	}

	return b
}
