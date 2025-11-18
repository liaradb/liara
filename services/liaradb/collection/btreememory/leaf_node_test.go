package btreememory

import "testing"

func TestLeafNode(t *testing.T) {
	t.Parallel()

	s := &mockStorage[int]{}

	ln := newLeafNode(s, 1, NewRecordID(0, 1))
	ln.insert(2, 2, NewRecordID(0, 2))
	ln2, ok := ln.insert(2, 3, NewRecordID(0, 3))
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

func TestLeafNode_Reverse(t *testing.T) {
	t.Parallel()

	s := &mockStorage[int]{}

	ln := newLeafNode(s, 3, NewRecordID(0, 3))
	ln.insert(2, 2, NewRecordID(0, 2))
	ln2, ok := ln.insert(2, 1, NewRecordID(0, 1))
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
