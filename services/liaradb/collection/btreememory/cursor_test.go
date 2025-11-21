package btreememory

import (
	"cmp"
	"context"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/storage"
)

func TestCursor_Default(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	bt := NewCursor(&mockStorage[int]{})
	if err := bt.CreateBTree(ctx); err != nil {
		t.Error(err)
	}

	testFanout(t, "default", bt, 3)
	testHeight(t, "default", bt, 1)

	if rid, err := bt.GetValue(ctx, 0); err == nil {
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
		skip    bool
	}{
		{message: "should insert",
			items: newItemsAscending(2), fanout: 3, height: 1, skip: false},
		{message: "should split leaf nodes",
			items: newItemsAscending(4), fanout: 3, height: 2, skip: false},
		{message: "should split key nodes",
			items: newItemsAscending(9), fanout: 3, height: 3, skip: true},
		{message: "should insert in any order",
			items: newItemsReversed(9), fanout: 3, height: 3, skip: true},
		{message: "should handle repeated items",
			items: newItems(1, 2, 2, 3), fanout: 3, height: 1, skip: true},
	} {
		t.Run(row.message, func(t *testing.T) {
			t.Parallel()
			if row.skip {
				t.Skip()
			}

			bt := NewCursor(&mockStorage[int]{})

			for _, i := range row.items {
				bt.Insert(t.Context(), i.key, i.rid)
			}

			testFanout(t, row.message, bt, row.fanout)
			testHeight(t, row.message, bt, row.height)
			// testCount(t, row.message, bt, len(row.items))
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
	// testCount(t, message, bt, 0)
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

// func testCount(t *testing.T, message string, bt *Cursor[int], count int) {
// 	t.Helper()

// 	if c, err := bt.Count(t.Context()); err != nil {
// 		t.Error(err)
// 	} else if c != count {
// 		t.Errorf("%v: should have a count of %v, recieved: %v", message, count, c)
// 	}
// }

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
	root  node[K]
	nodes map[storage.BlockID]node[K]
}

var _ Storage[int] = (*mockStorage[int])(nil)

func (m *mockStorage[K]) GetNode(ctx context.Context, bid storage.BlockID) (node[K], error) {
	if m.nodes == nil {
		return nil, ErrNotFound
	}

	n, ok := m.nodes[bid]
	if ok {
		return n, nil
	}

	// m.nodes = newKeyNode(m, )

	return nil, ErrNotFound
}

func (m *mockStorage[K]) GetKeyNode(ctx context.Context, bid storage.BlockID) (*keyNode[K], error) {
	if m.nodes == nil {
		return nil, ErrNotFound
	}

	n, ok := m.nodes[bid]
	if ok {
		return n.(*keyNode[K]), nil
	}

	// m.nodes = newKeyNode(m, )

	return nil, ErrNotFound
}

func (m *mockStorage[K]) GetLeafNode(ctx context.Context, bid storage.BlockID) (*leafNode[K], error) {
	if m.nodes == nil {
		return nil, ErrNotFound
	}

	n, ok := m.nodes[bid]
	if ok {
		return n.(*leafNode[K]), nil
	}

	// m.nodes = newKeyNode(m, )

	return nil, ErrNotFound
}

func (m *mockStorage[K]) InsertNode(ctx context.Context, bid storage.BlockID, n node[K]) error {
	if m.nodes == nil {
		m.nodes = make(map[storage.BlockID]node[K])
	}

	m.nodes[bid] = n
	return nil
}
