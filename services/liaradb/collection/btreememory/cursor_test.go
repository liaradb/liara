package btreememory

import (
	"cmp"
	"context"
	"slices"
	"testing"
)

func TestCursor_Default(t *testing.T) {
	t.Parallel()

	bt := NewCursor(&mockStorage[int]{})

	testFanout(t, "default", bt, 3)
	testHeight(t, "default", bt, 0)

	if rid, err := bt.GetValue(t.Context(), 0); err == nil {
		t.Error("should have no value by default")
	} else if rid != (RecordID{}) {
		t.Error("should have no value by default")
	}
}

func TestCursor_Insert(t *testing.T) {
	t.Parallel()

	for _, row := range []struct {
		message string
		items   []item
		fanout  int
		height  int
	}{
		{message: "should insert",
			items: newItemsAscending(2), fanout: 3, height: 1},
		{message: "should split leaf nodes",
			items: newItemsAscending(4), fanout: 3, height: 2},
		{message: "should split key nodes",
			items: newItemsAscending(9), fanout: 3, height: 3},
		{message: "should insert in any order",
			items: newItemsReversed(9), fanout: 3, height: 3},
		// {message: "should handle repeated items",
		// 	items: newItems(1, 2, 2, 3), fanout: 3, height: 1},
	} {
		t.Run(row.message, func(t *testing.T) {
			t.Parallel()

			bt := NewCursor(&mockStorage[int]{})

			for _, i := range row.items {
				bt.Insert(t.Context(), i.key, i.rid)
			}

			testFanout(t, row.message, bt, row.fanout)
			testHeight(t, row.message, bt, row.height)
			testCount(t, row.message, bt, len(row.items))
			testItems(t, row.message, bt, row.items)
		})
	}
}

func TestCursor_Delete(t *testing.T) {
	t.Parallel()

	bt := NewCursor(&mockStorage[int]{})

	if err := bt.Insert(t.Context(), 1, NewRecordID(0, 1)); err != nil {
		t.Error(err)
	}

	if err := bt.DeleteAll(t.Context(), 1); err != nil {
		t.Error(err)
	}

	message := "should delete"

	testFanout(t, message, bt, 3)
	testHeight(t, message, bt, 1)
	testCount(t, message, bt, 0)
	testItems(t, message, bt, []item{})
}

type item struct {
	key int
	rid RecordID
}

func newItem(i int) item {
	return item{i, NewRecordID(0, int8(i))}
}

func newItems(i ...int) []item {
	items := make([]item, 0, len(i))
	for _, i := range i {
		items = append(items, newItem(i))
	}
	return items
}

func newItemsAscending(count int) []item {
	items := make([]item, 0, count)
	for i := range count {
		items = append(items, newItem(i+1))
	}
	return items
}

func newItemsReversed(count int) []item {
	i := newItemsAscending(count)
	slices.Reverse(i)
	return i
}

func testFanout(t *testing.T, message string, bt *Cursor[int], fanout int) {
	t.Helper()

	if f := bt.FanOut(); f != fanout {
		t.Errorf("%v: should have a fanout of %v, recieved: %v", message, fanout, f)
	}
}

func testHeight(t *testing.T, message string, bt *Cursor[int], height int) {
	t.Helper()

	if h, err := bt.Height(t.Context()); err != nil {
		t.Error(err)
	} else if h != height {
		t.Errorf("%v: should have a height of %v, recieved: %v", message, height, h)
	}
}

func testCount(t *testing.T, message string, bt *Cursor[int], count int) {
	t.Helper()

	if c, err := bt.Count(t.Context()); err != nil {
		t.Error(err)
	} else if c != count {
		t.Errorf("%v: should have a count of %v, recieved: %v", message, count, c)
	}
}

func testItems(t *testing.T, message string, bt *Cursor[int], items []item) {
	t.Helper()

	for _, i := range items {
		if rid, err := bt.GetValue(t.Context(), i.key); err != nil {
			t.Error(err)
		} else if rid != i.rid {
			t.Errorf("%v: incorrect value: %v, expected: %v", message, rid, i.rid)
		}
	}
}

type mockStorage[K cmp.Ordered] struct {
	root node[K]
}

var _ Storage[int] = (*mockStorage[int])(nil)

func (m *mockStorage[K]) GetPage(context.Context) (node[K], error) {
	return nil, nil
}

func (m *mockStorage[K]) GetRoot(context.Context) (node[K], error) {
	if m.root == nil {
		return nil, ErrEmptyTree
	}

	return m.root, nil
}

func (m *mockStorage[K]) SetRoot(ctx context.Context, root node[K]) error {
	m.root = root
	return nil
}
