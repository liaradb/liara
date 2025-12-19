package btree

import (
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/btree/leafnode"
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

	if l, err := NewCursor(s).Level(ctx, n); err != nil {
		t.Error(err)
	} else if l != 0 {
		t.Errorf("incorrect level: %v, expected: %v", l, 0)
	}

	synctest.Wait()

	if p := s.CountPinned(); p != 0 {
		t.Errorf("incorrect pin count: %v, expected: %v", p, 0)
	}
}

func TestCursor_Insert__Root(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_Insert__Root)
}

type leafEntry struct {
	key      key.Key
	recordID leafnode.RecordID
}

func newLeafEntry(key key.Key, recordID leafnode.RecordID) leafEntry {
	return leafEntry{
		key:      key,
		recordID: recordID,
	}
}

func testCursor_Insert__Root(t *testing.T) {
	s := createStorage(t, 2, 256)
	ctx := t.Context()
	n := "testfile"

	data := []leafEntry{
		newLeafEntry(
			key.Key("a"),
			leafnode.NewRecordID(1, 2)),
		newLeafEntry(
			key.Key("b"),
			leafnode.NewRecordID(3, 4)),
		newLeafEntry(
			key.Key("c"),
			leafnode.NewRecordID(5, 6)),
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

	data := createData()

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

func TestCursor_Insert__Reverse(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_Insert__Reverse)
}

func testCursor_Insert__Reverse(t *testing.T) {
	s := createStorage(t, 8, 62)
	ctx := t.Context()
	n := "testfile"

	data := createData()

	for _, e := range reverseData(data) {
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

	data := createData()

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

	for i, e := range reorderData(order, data) {
		if err := NewCursor(s).Insert(ctx, n, e.key, e.recordID); err != nil {
			t.Fatal(i, e.key, err)
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

func TestCursor_SearchRange(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_SearchRange)
}

func testCursor_SearchRange(t *testing.T) {
	s := createStorage(t, 8, 62)
	ctx := t.Context()
	n := "testfile"

	data := createData()

	for _, e := range data {
		if err := NewCursor(s).Insert(ctx, n, e.key, e.recordID); err != nil {
			t.Fatal(e.key, err)
		}
		// TODO: Need to flush to disk
	}

	wantAll := make([]leafnode.RecordID, 0, len(data))
	for _, e := range data {
		wantAll = append(wantAll, e.recordID)
	}

	for i, e := range data {
		c := NewCursor(s)
		result := make([]leafnode.RecordID, 0, len(data))
		for rid, err := range c.SearchRange(ctx, n, e.key, 0, 0) {
			if err != nil {
				t.Fatal(err)
			}

			result = append(result, rid)
		}

		want := wantAll[i:]
		if !slices.Equal(result, want) {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}
	}

	// Skip and Limit
	{
		c := NewCursor(s)
		result := make([]leafnode.RecordID, 0, len(data))
		for rid, err := range c.SearchRange(ctx, n, "1", 1, 3) {
			if err != nil {
				t.Fatal(err)
			}

			result = append(result, rid)
		}

		want := wantAll[2:5]
		if !slices.Equal(result, want) {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
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

func reverseData(data []leafEntry) []leafEntry {
	data = slices.Clone(data)
	slices.Reverse(data)
	return data
}

func reorderData(order []int, data []leafEntry) []leafEntry {
	data = slices.Clone(data)

	hash := make(map[key.Key]int)
	for i, le := range data {
		if i >= len(order) {
			break
		}

		hash[le.key] = order[i]
	}
	l := len(data)

	slices.SortFunc(data, func(a, b leafEntry) int {
		c, ok := hash[a.key]
		if !ok {
			c = l
		}

		d, ok := hash[b.key]
		if !ok {
			d = l
		}

		return c - d
	})

	return data
}

func createData() []leafEntry {
	return []leafEntry{
		newLeafEntry(
			key.Key("0"),
			leafnode.NewRecordID(1, 2)),
		newLeafEntry(
			key.Key("1"),
			leafnode.NewRecordID(3, 4)),
		newLeafEntry(
			key.Key("2"),
			leafnode.NewRecordID(5, 6)),
		newLeafEntry(
			key.Key("3"),
			leafnode.NewRecordID(7, 8)),
		newLeafEntry(
			key.Key("4"),
			leafnode.NewRecordID(9, 10)),
		newLeafEntry(
			key.Key("5"),
			leafnode.NewRecordID(11, 12)),
		newLeafEntry(
			key.Key("6"),
			leafnode.NewRecordID(13, 14)),
		newLeafEntry(
			key.Key("7"),
			leafnode.NewRecordID(15, 16)),
		newLeafEntry(
			key.Key("8"),
			leafnode.NewRecordID(17, 18)),
	}
}
