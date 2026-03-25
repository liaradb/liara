package leafnode

import (
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/btree/node"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestLeafNode_Fill(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLeafNode_Fill)
}

func testLeafNode_Fill(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	p := node.New(b)
	ln := New(p)

	data := []leafEntry{
		newLeafEntry(
			key.NewKey([]byte("abcde")),
			link.NewRecordLocator(1, 2)),
		newLeafEntry(
			key.NewKey([]byte("fghij")),
			link.NewRecordLocator(3, 4)),
	}

	ln.Fill(1, 2, func(yield func(key.Key, link.RecordLocator) bool) {
		for _, le := range data {
			if !yield(le.Key(), le.RecordID()) {
				return
			}
		}
	})

	result := make([]leafEntry, 0, len(data))
	for key, rl := range ln.Children() {
		result = append(result, newLeafEntry(key, rl))
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}

	if l := ln.LeftID(); l != 1 {
		t.Errorf("incorrect left id: %v, expected: %v", l, 1)
	}

	if r := ln.RightID(); r != 2 {
		t.Errorf("incorrect left id: %v, expected: %v", r, 2)
	}
}

func TestLeafNode_Child(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLeafNode_Child)
}

func testLeafNode_Child(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	p := node.New(b)
	ln := New(p)

	data := []leafEntry{
		newLeafEntry(
			key.NewKey([]byte("abcde")),
			link.NewRecordLocator(1, 2)),
		newLeafEntry(
			key.NewKey([]byte("fghij")),
			link.NewRecordLocator(3, 4)),
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
	synctest.Test(t, testLeafNode_Children)
}

func testLeafNode_Children(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	p := node.New(b)
	ln := New(p)

	data := []leafEntry{
		newLeafEntry(
			key.NewKey([]byte("abcde")),
			link.NewRecordLocator(1, 2)),
		newLeafEntry(
			key.NewKey([]byte("fghij")),
			link.NewRecordLocator(3, 4)),
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
	synctest.Test(t, testLeafNode_Insert)
}

func testLeafNode_Insert(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	bp := node.New(b)
	ln := New(bp)

	data := []leafEntry{
		newLeafEntry(
			key.NewKey([]byte("a")),
			link.NewRecordLocator(1, 2)),
		newLeafEntry(
			key.NewKey([]byte("b")),
			link.NewRecordLocator(3, 4)),
		newLeafEntry(
			key.NewKey([]byte("c")),
			link.NewRecordLocator(5, 6)),
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
		want := make([]link.RecordLocator, 0, len(data))
		for _, e := range data {
			want = append(want, e.RecordID())
		}

		result := make([]link.RecordLocator, 0, len(want))
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
		want := make([]link.RecordLocator, 0, len(data[i:]))
		for _, e := range data[i:] {
			want = append(want, e.RecordID())
		}

		result := make([]link.RecordLocator, 0, len(want))
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
	synctest.Test(t, testLeafNode_Insert__Split)
}

func testLeafNode_Insert__Split(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)

	bp := node.New(b)
	ln := New(bp)

	data := []leafEntry{
		newLeafEntry(
			key.NewKey([]byte("a")),
			link.NewRecordLocator(1, 2)),
		newLeafEntry(
			key.NewKey([]byte("b")),
			link.NewRecordLocator(3, 4)),
		newLeafEntry(
			key.NewKey([]byte("c")),
			link.NewRecordLocator(5, 6)),
	}

	// Insert in mixed order
	order := []int{0, 2, 1}
	for _, i := range order {
		e := data[i]
		if a, b, ok := ln.Insert(e.key, e.recordID); !ok {
			t.Errorf("should insert:\n%v\n%v", a, b)
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

	b.Release()

	synctest.Wait()
}

func TestLeafNode_SetLeftID(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLeafNode_SetLeftID)
}

func testLeafNode_SetLeftID(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	bp := node.New(b)
	ln := New(bp)

	if fp := ln.LeftID(); fp != 0 {
		t.Errorf("incorrect left id: %v, expected: %v", fp, 0)
	}

	if b.Dirty() {
		t.Error("should not be dirty")
	}

	ln.SetLeftID(link.FilePosition(1))
	if fp := ln.LeftID(); fp != 1 {
		t.Errorf("incorrect left id: %v, expected: %v", fp, 1)
	}

	if !b.Dirty() {
		t.Error("should be dirty")
	}
}

func TestLeafNode_SetRightID(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLeafNode_SetRightID)
}

func testLeafNode_SetRightID(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	bp := node.New(b)
	ln := New(bp)

	if fp := ln.RightID(); fp != 0 {
		t.Errorf("incorrect right id: %v, expected: %v", fp, 0)
	}

	if b.Dirty() {
		t.Error("should not be dirty")
	}

	ln.SetRightID(link.FilePosition(1))
	if fp := ln.RightID(); fp != 1 {
		t.Errorf("incorrect right id: %v, expected: %v", fp, 1)
	}

	if !b.Dirty() {
		t.Error("should be dirty")
	}
}

func createBuffer(t *testing.T, s *storage.Storage) *storage.Buffer {
	b, err := s.Request(t.Context(), link.BlockID{})
	if err != nil {
		t.Fatal(err)
	}

	return b
}
