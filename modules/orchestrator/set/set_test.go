package set

import (
	"slices"
	"testing"
)

func TestSet_Add(t *testing.T) {
	t.Parallel()

	s := Set[int]{}
	s.Add(1, 2, 3)
	result := s.Slice()
	if len(result) != 3 {
		t.Error("incorrect length")
	}
	if !slices.ContainsFunc(result, func(i int) bool {
		return i == 1 || i == 2 || i == 3
	}) {
		t.Error("incorrect values")
	}
}

func TestSet_Remove(t *testing.T) {
	t.Parallel()

	s := New(1)
	s.Remove(1)
	slice := s.Slice()
	if len(slice) != 0 {
		t.Error("incorrect length")
	}
}

func TestSet_Clear(t *testing.T) {
	t.Parallel()

	s := New(1)
	s.Clear()
	slice := s.Slice()
	if len(slice) != 0 {
		t.Error("incorrect length")
	}
}

func TestSet_Includes(t *testing.T) {
	t.Parallel()

	s := New(1)
	if s.Includes(0) {
		t.Error("should only include 1")
	}
	if !s.Includes(1) {
		t.Error("should include 1")
	}
}

func TestSet_Merge(t *testing.T) {
	t.Parallel()

	a := New(0)
	b := New(1)
	a.Merge(b)
	if !a.Includes(0) {
		t.Error("should include 0")
	}
	if !a.Includes(1) {
		t.Error("should include 1")
	}
}

func TestSet_Union(t *testing.T) {
	t.Parallel()

	a := New(0)
	b := New(1)
	c := a.Union(b)
	if !c.Includes(0) {
		t.Error("should include 0")
	}
	if !c.Includes(1) {
		t.Error("should include 1")
	}
}

func TestSet_Intersection(t *testing.T) {
	t.Parallel()

	a := New(0, 1)
	b := New(1, 2)
	c := a.Intersection(b)
	if c.Includes(0) {
		t.Error("should not include 0")
	}
	if !c.Includes(1) {
		t.Error("should include 1")
	}
	if c.Includes(2) {
		t.Error("should not include 2")
	}
}
