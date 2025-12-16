package btree

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/file/filetesting"
	"github.com/liaradb/liaradb/storage"
)

// TODO: Test this
// {message: "should insert",
// 	items: newItemsAscending(2), fanout: 3, height: 1, skip: false},
// {message: "should split leaf nodes",
// 	items: newItemsAscending(4), fanout: 3, height: 2, skip: false},
// {message: "should split key nodes",
// 	items: newItemsAscending(9), fanout: 3, height: 3, skip: false},
// {message: "should insert in any order",
// 	items: newItemsReversed(9), fanout: 3, height: 3, skip: true},
// {message: "should handle repeated items",
// 	items: newItems(1, 2, 2, 3), fanout: 3, height: 1, skip: true},

func TestCursor_GetRoot_Default(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor)
}

func testCursor(t *testing.T) {
	s := createStorage(t, 2, 256)
	ctx := t.Context()

	n := "testfile"
	r, err := NewCursor(s).GetPage(ctx, storage.NewBlockID(n, 0))
	if err != nil {
		t.Error(err)
	}

	if l := r.Level(); l != 0 {
		t.Errorf("incorrect level: %v, expected: %v", l, 0)
	}

	r.Release()

	synctest.Wait()

	if p := s.CountPinned(); p != 0 {
		t.Errorf("incorrect pin count: %v, expected: %v", p, 0)
	}
}

func TestCursor_Insert__Root(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_Insert__Root)
}

func testCursor_Insert__Root(t *testing.T) {
	s := createStorage(t, 2, 256)
	ctx := t.Context()
	n := "testfile"

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

	for _, e := range data {
		if err := NewCursor(s).Insert(ctx, n, e.key, e.recordID); err != nil {
			t.Error(err)
		}
		// TODO: Need to flush to disk
	}

	for _, e := range data {
		if rid, err := NewCursor(s).Search(ctx, n, e.key); err != nil {
			t.Fatal(err)
		} else if rid != e.recordID {
			t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
		}
	}

	synctest.Wait()

	if p := s.CountPinned(); p != 0 {
		t.Errorf("incorrect pin count: %v, expected: %v", p, 0)
	}
}

func TestCursor_Insert__RootSplit(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_Insert__RootSplit)
}

func testCursor_Insert__RootSplit(t *testing.T) {
	// TODO: Why does this use so many buffers?
	//                            [7]
	//               .............   .........
	//         [3        5]                 [9]
	//     ....   ......   ....          ...   ..
	// [1   2]   [3   4]   [5   6]   [7   8]   [9]

	s := createStorage(t, 8, 62)
	ctx := t.Context()
	n := "testfile"

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
		newLeafEntry(
			Key("d"),
			NewRecordID(7, 8)),
		newLeafEntry(
			Key("e"),
			NewRecordID(9, 10)),
		newLeafEntry(
			Key("f"),
			NewRecordID(11, 12)),
		newLeafEntry(
			Key("g"),
			NewRecordID(13, 14)),
		newLeafEntry(
			Key("h"),
			NewRecordID(15, 16)),
		newLeafEntry(
			Key("i"),
			NewRecordID(17, 18)),
	}

	for _, e := range data {
		if err := NewCursor(s).Insert(ctx, n, e.key, e.recordID); err != nil {
			t.Fatal(e.key, err)
		}
		// TODO: Need to flush to disk
	}

	for _, e := range data {
		if rid, err := NewCursor(s).Search(ctx, n, e.key); err != nil {
			t.Error(err, e.key)
		} else if rid != e.recordID {
			t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
		}
	}

	synctest.Wait()

	if p := s.CountPinned(); p != 0 {
		t.Errorf("incorrect pin count: %v, expected: %v", p, 0)
	}
}

func TestCursor_Insert__Random(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_Insert__Random)
}

func testCursor_Insert__Random(t *testing.T) {
	s := createStorage(t, 8, 62)
	ctx := t.Context()
	n := "testfile"

	data := []LeafEntry{
		newLeafEntry(
			Key("0"),
			NewRecordID(1, 2)),
		newLeafEntry(
			Key("1"),
			NewRecordID(3, 4)),
		newLeafEntry(
			Key("2"),
			NewRecordID(5, 6)),
		newLeafEntry(
			Key("3"),
			NewRecordID(7, 8)),
		newLeafEntry(
			Key("4"),
			NewRecordID(9, 10)),
		newLeafEntry(
			Key("5"),
			NewRecordID(11, 12)),
		newLeafEntry(
			Key("6"),
			NewRecordID(13, 14)),
		newLeafEntry(
			Key("7"),
			NewRecordID(15, 16)),
		newLeafEntry(
			Key("8"),
			NewRecordID(17, 18)),
	}

	// Insert in mixed order
	order := []int{
		5,
		8,
		2,
		4,
		3,
		6,
		0,
		7,
		1,
	}
	for index, i := range order {
		e := data[i]
		if err := NewCursor(s).Insert(ctx, n, e.key, e.recordID); err != nil {
			t.Fatal(index, i, e.key, err)
		}
		// TODO: Need to flush to disk
	}

	for _, i := range order {
		e := data[i]
		if rid, err := NewCursor(s).Search(ctx, n, e.key); err != nil {
			t.Error(err, e.key)
		} else if rid != e.recordID {
			t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
		}
	}

	synctest.Wait()

	if p := s.CountPinned(); p != 0 {
		t.Errorf("incorrect pin count: %v, expected: %v", p, 0)
	}
}

func createStorage(t *testing.T, max int, bs int64) *storage.Storage {
	fsys := filetesting.NewMockFileSystem(t, nil)
	s := storage.New(fsys, max, bs, t.TempDir())
	if err := s.Run(t.Context()); err != nil {
		t.Fatal(err)
	}

	return s
}
