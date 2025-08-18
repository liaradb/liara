package btree

import "testing"

func TestBTree_Default(t *testing.T) {
	bt := BTree[int, string]{}

	if f := bt.FanOut(); f != 3 {
		t.Error("should have a fanout of 3")
	}

	// if h := bt.Height(); h != 0 {
	// 	t.Error("should have a height of 0")
	// }

	if v, ok := bt.getValue(0); ok {
		t.Error("should have no value by default")
	} else if v != "" {
		t.Error("should have no value by default")
	}
}

func TestBTree_Insert(t *testing.T) {
	bt := BTree[int, string]{}

	bt.insert(1, "a")
	bt.insert(2, "b")

	if f := bt.FanOut(); f != 3 {
		t.Error("should have a fanout of 3")
	}

	// if h := bt.Height(); h != 1 {
	// 	t.Error("should have a height of 1")
	// }

	if v, ok := bt.getValue(1); !ok {
		t.Error("should have a value")
	} else if v != "a" {
		t.Error("should have a value")
	}

	if v, ok := bt.getValue(2); !ok {
		t.Error("should have a value")
	} else if v != "b" {
		t.Error("should have a value")
	}
}

func TestBTree_Height(t *testing.T) {
	t.Skip()
	bt := BTree[int, string]{}

	bt.insert(1, "a")
	bt.insert(2, "b")
	bt.insert(3, "c")
	bt.insert(4, "d")

	if f := bt.FanOut(); f != 3 {
		t.Error("should have a fanout of 3")
	}

	// if h := bt.Height(); h != 2 {
	// 	t.Error("should have a height of 2")
	// }

	if v, ok := bt.getValue(1); !ok {
		t.Error("should have a value")
	} else if v != "a" {
		t.Error("should have a value")
	}

	if v, ok := bt.getValue(2); !ok {
		t.Error("should have a value")
	} else if v != "b" {
		t.Error("should have a value")
	}
}
