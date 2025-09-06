package queue

import "testing"

func TestQueue_PopDefault(t *testing.T) {
	q := Queue[int, string]{}

	if v, ok := q.Pop(); ok {
		t.Error("should not return a value")
	} else if v != "" {
		t.Errorf("should not return a value")
	}
}

func TestQueue_PopValue(t *testing.T) {
	q := Queue[int, string]{}

	q.Push(1, "a")
	q.Push(2, "b")

	if v, ok := q.Pop(); !ok {
		t.Error("should return a value")
	} else if v != "a" {
		t.Errorf("returned incorrect value: %v, expected: %v", v, "a")
	}

	if v, ok := q.Pop(); !ok {
		t.Error("should return a value")
	} else if v != "b" {
		t.Errorf("returned incorrect value: %v, expected: %v", v, "b")
	}
}

func TestQueue_RemoveDefault(t *testing.T) {
	q := Queue[int, string]{}

	if v, ok := q.Remove(1); ok {
		t.Error("should not return a value")
	} else if v != "" {
		t.Errorf("should not return a value")
	}
}

func TestQueue_RemoveValue(t *testing.T) {
	q := Queue[int, string]{}

	q.Push(1, "a")
	q.Push(2, "b")
	q.Push(3, "c")

	if v, ok := q.Remove(2); !ok {
		t.Error("should return a value")
	} else if v != "b" {
		t.Errorf("returned incorrect value: %v, expected: %v", v, "b")
	}

	if v, ok := q.Remove(3); !ok {
		t.Error("should return a value")
	} else if v != "c" {
		t.Errorf("returned incorrect value: %v, expected: %v", v, "c")
	}

	if v, ok := q.Remove(1); !ok {
		t.Error("should return a value")
	} else if v != "a" {
		t.Errorf("returned incorrect value: %v, expected: %v", v, "a")
	}
}
