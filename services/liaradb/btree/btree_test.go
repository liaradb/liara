package btree

import "testing"

func TestBTree_Default(t *testing.T) {
	t.Parallel()

	bt := &BTree[int, string]{}

	testFanout(t, bt, 3)
	testHeight(t, bt, 0)

	if v, ok := bt.GetValue(0); ok {
		t.Error("should have no value by default")
	} else if v != "" {
		t.Error("should have no value by default")
	}
}

func TestBTree_Insert(t *testing.T) {
	t.Parallel()

	bt := &BTree[int, string]{}

	items := newItems(2)

	for _, i := range items {
		bt.Insert(i.key, i.value)
	}

	testFanout(t, bt, 3)
	testHeight(t, bt, 1)
	testItems(t, bt, items)
}

func TestBTree_SplitLeafNode(t *testing.T) {
	t.Parallel()

	bt := &BTree[int, string]{}

	items := newItems(4)

	for _, i := range items {
		bt.Insert(i.key, i.value)
	}

	testFanout(t, bt, 3)
	testHeight(t, bt, 2)
	testItems(t, bt, items)
}

func TestBTree_SplitKeyNode(t *testing.T) {
	t.Parallel()

	bt := &BTree[int, string]{}

	items := newItems(9)

	for _, i := range items {
		bt.Insert(i.key, i.value)
	}

	testFanout(t, bt, 3)
	testHeight(t, bt, 3)
	testItems(t, bt, items)
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

func testFanout(t *testing.T, bt *BTree[int, string], fanout int) {
	t.Helper()

	if f := bt.FanOut(); f != fanout {
		t.Errorf("should have a fanout of %v, recieved: %v", fanout, f)
	}
}

func testHeight(t *testing.T, bt *BTree[int, string], height int) {
	t.Helper()

	if h := bt.Height(); h != height {
		t.Errorf("should have a height of %v, recieved: %v", height, h)
	}
}

func testItems(t *testing.T, bt *BTree[int, string], items []item) {
	t.Helper()

	for _, i := range items {
		if v, ok := bt.GetValue(i.key); !ok {
			t.Error("should have a value")
		} else if v != i.value {
			t.Errorf("incorrect value: %v, expected: %v", v, i.value)
		}
	}
}
