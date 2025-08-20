package btree

import (
	"slices"
	"testing"
)

func TestBTree_Default(t *testing.T) {
	t.Parallel()

	bt := &BTree[int, string]{}

	testFanout(t, "default", bt, 3)
	testHeight(t, "default", bt, 0)

	if v, ok := bt.GetValue(0); ok {
		t.Error("should have no value by default")
	} else if v != "" {
		t.Error("should have no value by default")
	}
}

func TestBTree_Insert(t *testing.T) {
	t.Parallel()

	for _, row := range []struct {
		message string
		items   []item
		fanout  int
		height  int
	}{
		{message: "should insert",
			items: newItems(2), fanout: 3, height: 1},
		{message: "should split leaf nodes",
			items: newItems(4), fanout: 3, height: 2},
		{message: "should split key nodes",
			items: newItems(9), fanout: 3, height: 3},
		{message: "should insert in any order",
			items: newItemsReverse(9), fanout: 3, height: 3},
	} {
		t.Run(row.message, func(t *testing.T) {
			t.Parallel()

			bt := &BTree[int, string]{}

			for _, i := range row.items {
				bt.Insert(i.key, i.value)
			}

			testFanout(t, row.message, bt, row.fanout)
			testHeight(t, row.message, bt, row.height)
			testItems(t, row.message, bt, row.items)
		})
	}
}

type item struct {
	key   int
	value string
}

func newItems(count int) []item {
	items := make([]item, 0, count)
	for i := range count {
		items = append(items, item{i + 1, string(rune('a' + i))})
	}
	return items
}

func newItemsReverse(count int) []item {
	i := newItems(count)
	slices.Reverse(i)
	return i
}

func testFanout(t *testing.T, message string, bt *BTree[int, string], fanout int) {
	t.Helper()

	if f := bt.FanOut(); f != fanout {
		t.Errorf("%v: should have a fanout of %v, recieved: %v", message, fanout, f)
	}
}

func testHeight(t *testing.T, message string, bt *BTree[int, string], height int) {
	t.Helper()

	if h := bt.Height(); h != height {
		t.Errorf("%v: should have a height of %v, recieved: %v", message, height, h)
	}
}

func testItems(t *testing.T, message string, bt *BTree[int, string], items []item) {
	t.Helper()

	for _, i := range items {
		if v, ok := bt.GetValue(i.key); !ok {
			t.Errorf("%v: %v should have a value", message, i.key)
		} else if v != i.value {
			t.Errorf("%v: incorrect value: %v, expected: %v", message, v, i.value)
		}
	}
}
