package btree

import "testing"

func TestLeafNode(t *testing.T) {
	ln := newLeafNode(1, "a")
	ln.insert(2, 2, "b")
	ln2, ok := ln.insert(2, 3, "c")
	if !ok {
		t.Error("should split")
	}

	if l := ln.count(); l != 1 {
		t.Errorf("incorrect split.  Expected: %v, Recieved: %v", 1, l)
	}

	if l := ln2.count(); l != 2 {
		t.Errorf("incorrect split.  Expected: %v, Recieved: %v", 2, l)
	}
}
