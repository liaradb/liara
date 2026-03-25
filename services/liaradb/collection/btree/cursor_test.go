package btree

import (
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

// TODO: Test latching
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
	s := storagetesting.CreateStorage(t, 2, 256)
	ctx := t.Context()
	fn := link.NewFileName("testfile")

	if l, err := NewCursor(s).Level(ctx, fn); err != nil {
		t.Error(err)
	} else if l != 0 {
		t.Errorf("incorrect level: %v, expected: %v", l, 0)
	}
}

func TestCursor_Insert__Root(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_Insert__Root)
}

type leafEntry struct {
	key      key.Key
	recordID link.RecordLocator
}

func newLeafEntry(key key.Key, recordID link.RecordLocator) leafEntry {
	return leafEntry{
		key:      key,
		recordID: recordID,
	}
}

func testCursor_Insert__Root(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	ctx := t.Context()
	fn := link.NewFileName("testfile")

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

	for _, e := range data {
		if err := NewCursor(s).Insert(ctx, fn, e.key, e.recordID); err != nil {
			t.Error(err)
		}
		// TODO: Need to flush to disk
	}

	for _, e := range data {
		if rid, err := NewCursor(s).Search(ctx, fn, e.key); err != nil {
			t.Fatal(err)
		} else if rid != e.recordID {
			t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
		}
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

	s := storagetesting.CreateStorage(t, 8, 72)
	ctx := t.Context()
	fn := link.NewFileName("testfile")

	data := createData()

	for _, e := range data {
		if err := NewCursor(s).Insert(ctx, fn, e.key, e.recordID); err != nil {
			t.Fatal(e.key, err)
		}
		// TODO: Need to flush to disk
	}

	for _, e := range data {
		if rid, err := NewCursor(s).Search(ctx, fn, e.key); err != nil {
			t.Error(err, e.key)
		} else if rid != e.recordID {
			t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
		}
	}
}

func TestCursor_Insert__Reverse(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_Insert__Reverse)
}

func testCursor_Insert__Reverse(t *testing.T) {
	s := storagetesting.CreateStorage(t, 8, 72)
	ctx := t.Context()
	fn := link.NewFileName("testfile")

	data := createData()

	for _, e := range reverseData(data) {
		if err := NewCursor(s).Insert(ctx, fn, e.key, e.recordID); err != nil {
			t.Fatal(e.key, err)
		}
		// TODO: Need to flush to disk
	}

	for _, e := range data {
		if rid, err := NewCursor(s).Search(ctx, fn, e.key); err != nil {
			t.Error(err, e.key)
		} else if rid != e.recordID {
			t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
		}
	}
}

func TestCursor_Insert__Random(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_Insert__Random)
}

func testCursor_Insert__Random(t *testing.T) {
	s := storagetesting.CreateStorage(t, 8, 72)
	ctx := t.Context()
	fn := link.NewFileName("testfile")

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
		if err := NewCursor(s).Insert(ctx, fn, e.key, e.recordID); err != nil {
			t.Fatal(i, e.key, err)
		}
		// TODO: Need to flush to disk
	}

	for _, i := range order {
		e := data[i]
		if rid, err := NewCursor(s).Search(ctx, fn, e.key); err != nil {
			t.Error(err, e.key)
		} else if rid != e.recordID {
			t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
		}
	}
}

func TestCursor_Insert__Existing(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_Insert__Existing)
}

func testCursor_Insert__Existing(t *testing.T) {
	s := storagetesting.CreateStorage(t, 8, 62)
	ctx := t.Context()
	fn := link.NewFileName("testfile")

	le := newLeafEntry(
		key.NewKey([]byte("0")),
		link.NewRecordLocator(1, 2))

	if err := NewCursor(s).Insert(ctx, fn, le.key, le.recordID); err != nil {
		t.Fatal(le.key, err)
	}

	if err := NewCursor(s).Insert(ctx, fn, le.key, le.recordID); err == nil {
		t.Error("should not insert the same key")
	}

	c := 0
	for _, err := range NewCursor(s).All(ctx, fn, 0, 0) {
		if err != nil {
			t.Error(err)
		}

		c++
	}

	if c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}
}

func TestCursor_SearchRange(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_SearchRange)
}

func testCursor_SearchRange(t *testing.T) {
	s := storagetesting.CreateStorage(t, 8, 72)
	ctx := t.Context()
	fn := link.NewFileName("testfile")

	data := createData()

	for _, e := range data {
		if err := NewCursor(s).Insert(ctx, fn, e.key, e.recordID); err != nil {
			t.Fatal(e.key, err)
		}
		// TODO: Need to flush to disk
	}

	wantAll := make([]link.RecordLocator, 0, len(data))
	for _, e := range data {
		wantAll = append(wantAll, e.recordID)
	}

	for i, e := range data {
		c := NewCursor(s)
		result := make([]link.RecordLocator, 0, len(data))
		for rid, err := range c.SearchRange(ctx, fn, e.key, 0, 0) {
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
		result := make([]link.RecordLocator, 0, len(data))
		for rid, err := range c.SearchRange(ctx, fn, key.NewKey([]byte("1")), 1, 3) {
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
}

func TestCursor_All(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_All)
}

func testCursor_All(t *testing.T) {
	s := storagetesting.CreateStorage(t, 8, 72)
	ctx := t.Context()
	fn := link.NewFileName("testfile")

	data := createData()

	for _, e := range data {
		if err := NewCursor(s).Insert(ctx, fn, e.key, e.recordID); err != nil {
			t.Fatal(e.key, err)
		}
		// TODO: Need to flush to disk
	}

	wantAll := make([]link.RecordLocator, 0, len(data))
	for _, e := range data {
		wantAll = append(wantAll, e.recordID)
	}

	for i := range data {
		c := NewCursor(s)
		result := make([]link.RecordLocator, 0, len(data))
		for rid, err := range c.All(ctx, fn, i, 0) {
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
		result := make([]link.RecordLocator, 0, len(data))
		for rid, err := range c.All(ctx, fn, 1, 3) {
			if err != nil {
				t.Fatal(err)
			}

			result = append(result, rid)
		}

		want := wantAll[1:4]
		if !slices.Equal(result, want) {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}
	}
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
			key.NewKey([]byte("0")),
			link.NewRecordLocator(1, 2)),
		newLeafEntry(
			key.NewKey([]byte("1")),
			link.NewRecordLocator(3, 4)),
		newLeafEntry(
			key.NewKey([]byte("2")),
			link.NewRecordLocator(5, 6)),
		newLeafEntry(
			key.NewKey([]byte("3")),
			link.NewRecordLocator(7, 8)),
		newLeafEntry(
			key.NewKey([]byte("4")),
			link.NewRecordLocator(9, 10)),
		newLeafEntry(
			key.NewKey([]byte("5")),
			link.NewRecordLocator(11, 12)),
		newLeafEntry(
			key.NewKey([]byte("6")),
			link.NewRecordLocator(13, 14)),
		newLeafEntry(
			key.NewKey([]byte("7")),
			link.NewRecordLocator(15, 16)),
		newLeafEntry(
			key.NewKey([]byte("8")),
			link.NewRecordLocator(17, 18)),
	}
}
