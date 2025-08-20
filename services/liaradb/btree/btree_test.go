package btree

import "testing"

func TestBTree_Default(t *testing.T) {
	bt := BTree[int, string]{}
	fanout := 3

	if f := bt.FanOut(); f != fanout {
		t.Errorf("should have a fanout of %v, recieved: %v", fanout, f)
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
	fanout := 3

	bt.insert(1, "a")
	bt.insert(2, "b")

	if f := bt.FanOut(); f != fanout {
		t.Errorf("should have a fanout of %v, recieved: %v", fanout, f)
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

func TestBTree_SplitLeafNode(t *testing.T) {
	bt := BTree[int, string]{}
	fanout := 3

	bt.insert(1, "a")
	bt.insert(2, "b")
	bt.insert(3, "c")
	bt.insert(4, "d")

	if f := bt.FanOut(); f != fanout {
		t.Errorf("should have a fanout of %v, recieved: %v", fanout, f)
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

	if v, ok := bt.getValue(3); !ok {
		t.Error("should have a value")
	} else if v != "c" {
		t.Error("should have a value")
	}

	if v, ok := bt.getValue(4); !ok {
		t.Error("should have a value")
	} else if v != "d" {
		t.Error("should have a value")
	}
}

func TestBTree_SplitKeyNode(t *testing.T) {
	bt := BTree[int, string]{}
	fanout := 3

	bt.insert(1, "a")
	bt.insert(2, "b")
	bt.insert(3, "c")
	bt.insert(4, "d")
	bt.insert(5, "e")
	bt.insert(6, "f")
	bt.insert(7, "g")
	bt.insert(8, "h")
	bt.insert(9, "i")

	if f := bt.FanOut(); f != fanout {
		t.Errorf("should have a fanout of %v, recieved: %v", fanout, f)
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

	if v, ok := bt.getValue(3); !ok {
		t.Error("should have a value")
	} else if v != "c" {
		t.Error("should have a value")
	}

	if v, ok := bt.getValue(4); !ok {
		t.Error("should have a value")
	} else if v != "d" {
		t.Error("should have a value")
	}

	if v, ok := bt.getValue(5); !ok {
		t.Error("should have a value")
	} else if v != "e" {
		t.Error("should have a value")
	}

	if v, ok := bt.getValue(6); !ok {
		t.Error("should have a value")
	} else if v != "f" {
		t.Error("should have a value")
	}

	if v, ok := bt.getValue(7); !ok {
		t.Error("should have a value")
	} else if v != "g" {
		t.Error("should have a value")
	}

	if v, ok := bt.getValue(8); !ok {
		t.Error("should have a value")
	} else if v != "h" {
		t.Error("should have a value")
	}

	if v, ok := bt.getValue(9); !ok {
		t.Error("should have a value")
	} else if v != "i" {
		t.Error("should have a value")
	}
}
