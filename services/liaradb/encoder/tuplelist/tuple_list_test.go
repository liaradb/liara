package tuplelist

import (
	"slices"
	"testing"
)

func TestTupleList_Default(t *testing.T) {
	t.Parallel()

	l := New([]byte{})

	if length := l.Length(); length != 0 {
		t.Errorf("incorrect length: %v, expected: %v", length, 0)
	}

	if s := l.Size(); s != 2 {
		t.Errorf("incorrect size: %v, expected: %v", s, 2)
	}

	if c := l.Count(); c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}
}

func TestTupleList_Push(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 16))

	if length := l.Length(); length != 16 {
		t.Errorf("incorrect length: %v, expected: %v", length, 16)
	}

	if i, ok := l.Push(1, 10); !ok {
		t.Error("should push")
	} else if i != 0 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	if s := l.Size(); s != 6 {
		t.Errorf("incorrect size: %v, expected: %v", s, 6)
	}

	if c := l.Count(); c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}

	if i, ok := l.Push(2, 20); !ok {
		t.Error("should push")
	} else if i != 1 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	if s := l.Size(); s != 10 {
		t.Errorf("incorrect size: %v, expected: %v", s, 10)
	}

	if c := l.Count(); c != 2 {
		t.Errorf("incorrect count: %v, expected: %v", c, 2)
	}

	if a, b, ok := l.Item(0); !ok {
		t.Errorf("should have value")
	} else if a != 1 {
		t.Errorf("incorrect value: %v, expected: %v", a, 1)
	} else if b != 10 {
		t.Errorf("incorrect value: %v, expected: %v", b, 10)
	}

	if a, b, ok := l.Item(1); !ok {
		t.Errorf("should have value")
	} else if a != 2 {
		t.Errorf("incorrect value: %v, expected: %v", a, 2)
	} else if b != 20 {
		t.Errorf("incorrect value: %v, expected: %v", b, 20)
	}
}

func TestTupleList_Pop(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 16))

	if _, ok := l.Push(1, 10); !ok {
		t.Error("should push")
	}

	if _, ok := l.Push(2, 20); !ok {
		t.Error("should push")
	}

	if a, b, ok := l.Pop(); !ok {
		t.Error("should pop")
	} else if a != 2 {
		t.Errorf("incorrect value: %v, expected: %v", a, 2)
	} else if b != 20 {
		t.Errorf("incorrect value: %v, expected: %v", b, 20)
	}

	if s := l.Size(); s != 6 {
		t.Errorf("incorrect size: %v, expected: %v", s, 6)
	}

	if c := l.Count(); c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}

	if a, b, ok := l.Pop(); !ok {
		t.Error("should pop")
	} else if a != 1 {
		t.Errorf("incorrect value: %v, expected: %v", a, 1)
	} else if b != 10 {
		t.Errorf("incorrect value: %v, expected: %v", b, 10)
	}

	if s := l.Size(); s != 2 {
		t.Errorf("incorrect size: %v, expected: %v", s, 2)
	}

	if c := l.Count(); c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}

	if _, _, ok := l.Pop(); ok {
		t.Error("should not pop beyond empty")
	}
}

func TestTupleList_Items(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 32))

	as := []int16{10, 20, 30, 40, 50}
	bs := []int16{60, 70, 80, 90, 100}

	for i, item := range as {
		if _, ok := l.Push(item, bs[i]); !ok {
			t.Error("should push")
		}
	}

	resultA := make([]int16, 0, len(as))
	resultB := make([]int16, 0, len(bs))
	for a, b := range l.Items() {
		resultA = append(resultA, a)
		resultB = append(resultB, b)
	}

	if !slices.Equal(resultA, as) {
		t.Errorf("incorrect result: %v, expected: %v", resultA, as)
	}

	if !slices.Equal(resultB, bs) {
		t.Errorf("incorrect result: %v, expected: %v", resultB, bs)
	}
}

// TODO: Should not affect items outside of range
func TestTupleList_Insert(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 32))

	as := []int16{10, 20, 40, 50}
	bs := []int16{60, 70, 90, 100}

	wantAs := []int16{10, 20, 30, 40, 50}
	wantBs := []int16{60, 70, 80, 90, 100}

	for i, item := range as {
		if _, ok := l.Push(item, bs[i]); !ok {
			t.Error("should push")
		}
	}

	if _, ok := l.Insert(30, 80, 2); !ok {
		t.Error("should insert")
	}

	if c := l.Count(); c != 5 {
		t.Errorf("incorrect count: %v, expected: %v", c, 5)
	}

	resultA := make([]int16, 0, len(as))
	resultB := make([]int16, 0, len(bs))
	for a, b := range l.Items() {
		resultA = append(resultA, a)
		resultB = append(resultB, b)
	}

	if !slices.Equal(resultA, wantAs) {
		t.Errorf("incorrect result: %v, expected: %v", resultA, wantAs)
	}

	if !slices.Equal(resultB, wantBs) {
		t.Errorf("incorrect result: %v, expected: %v", resultB, wantBs)
	}
}
