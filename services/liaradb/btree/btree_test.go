package btree

import "testing"

func TestBTree_Default(t *testing.T) {
	bt := BTree[int, string]{}
	fanout := 3

	if f := bt.FanOut(); f != fanout {
		t.Errorf("should have a fanout of %v, recieved: %v", fanout, f)
	}

	if h := bt.Height(); h != 0 {
		t.Errorf("should have a height of %v, recieved: %v", 0, h)
	}

	if v, ok := bt.getValue(0); ok {
		t.Error("should have no value by default")
	} else if v != "" {
		t.Error("should have no value by default")
	}
}

func TestBTree_Insert(t *testing.T) {
	bt := BTree[int, string]{}
	fanout := 3

	if f := bt.FanOut(); f != fanout {
		t.Errorf("should have a fanout of %v, recieved: %v", fanout, f)
	}

	items := []struct {
		key   int
		value string
	}{
		{1, "a"},
		{2, "b"},
	}

	for _, i := range items {
		bt.insert(i.key, i.value)
	}

	if h := bt.Height(); h != 1 {
		t.Errorf("should have a height of %v, recieved: %v", 1, h)
	}

	if v, ok := bt.getValue(1); !ok {
		t.Error("should have a value")
	} else if v != "a" {
		t.Error("should have a value")
	}

	for _, i := range items {
		if v, ok := bt.getValue(i.key); !ok {
			t.Error("should have a value")
		} else if v != i.value {
			t.Errorf("incorrect value: %v, expected: %v", v, i.value)
		}
	}
}

func TestBTree_SplitLeafNode(t *testing.T) {
	bt := BTree[int, string]{}
	fanout := 3

	if f := bt.FanOut(); f != fanout {
		t.Errorf("should have a fanout of %v, recieved: %v", fanout, f)
	}

	items := []struct {
		key   int
		value string
	}{
		{1, "a"},
		{2, "b"},
		{3, "c"},
		{4, "d"},
	}

	for _, i := range items {
		bt.insert(i.key, i.value)
	}

	if h := bt.Height(); h != 2 {
		t.Errorf("should have a height of %v, recieved: %v", 2, h)
	}

	for _, i := range items {
		if v, ok := bt.getValue(i.key); !ok {
			t.Error("should have a value")
		} else if v != i.value {
			t.Errorf("incorrect value: %v, expected: %v", v, i.value)
		}
	}
}

func TestBTree_SplitKeyNode(t *testing.T) {
	bt := BTree[int, string]{}
	fanout := 3

	if f := bt.FanOut(); f != fanout {
		t.Errorf("should have a fanout of %v, recieved: %v", fanout, f)
	}

	items := []struct {
		key   int
		value string
	}{
		{1, "a"},
		{2, "b"},
		{3, "c"},
		{4, "d"},
		{5, "e"},
		{6, "f"},
		{7, "g"},
		{8, "h"},
		{9, "i"},
	}

	for _, i := range items {
		bt.insert(i.key, i.value)
	}

	if h := bt.Height(); h != 3 {
		t.Errorf("should have a height of %v, recieved: %v", 3, h)
	}

	for _, i := range items {
		if v, ok := bt.getValue(i.key); !ok {
			t.Error("should have a value")
		} else if v != i.value {
			t.Errorf("incorrect value: %v, expected: %v", v, i.value)
		}
	}
}
