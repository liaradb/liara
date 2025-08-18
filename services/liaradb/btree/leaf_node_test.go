package btree

import "testing"

func TestLeafNode(t *testing.T) {
	ln := newLeafNode(1, "a")
	ln.insert(2, "b")
	ln.insert(3, "c")
	ln1 := ln.split()

	if l := len(ln.children); l != 1 {
		t.Errorf("incorrect split.  Expected: %v, Recieved: %v", 1, l)
	}

	if l := len(ln1.children); l != 2 {
		t.Errorf("incorrect split.  Expected: %v, Recieved: %v", 2, l)
	}
}
