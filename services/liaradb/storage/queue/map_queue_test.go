package queue

import "testing"

func TestQueue_Count(t *testing.T) {
	mq := MapQueue[int, string]{}

	if c := mq.Count(); c != 0 {
		t.Errorf("should be empty")
	}

	mq.Push(1, "a")

	if c := mq.Count(); c != 1 {
		t.Errorf("should be 1")
	}

	mq.Push(2, "b")

	if c := mq.Count(); c != 2 {
		t.Errorf("should be 2")
	}

	mq.Pop()

	if c := mq.Count(); c != 1 {
		t.Errorf("should be 1")
	}

	mq.Pop()

	if c := mq.Count(); c != 0 {
		t.Errorf("should be empty")
	}
}

func TestQueue_PopDefault(t *testing.T) {
	mq := MapQueue[int, string]{}

	if v, ok := mq.Pop(); ok {
		t.Error("should not return a value")
	} else if v != "" {
		t.Errorf("should not return a value")
	}
}

func TestQueue_PopValue(t *testing.T) {
	mq := MapQueue[int, string]{}

	mq.Push(1, "a")
	mq.Push(2, "b")

	if v, ok := mq.Pop(); !ok {
		t.Error("should return a value")
	} else if v != "a" {
		t.Errorf("returned incorrect value: %v, expected: %v", v, "a")
	}

	if v, ok := mq.Pop(); !ok {
		t.Error("should return a value")
	} else if v != "b" {
		t.Errorf("returned incorrect value: %v, expected: %v", v, "b")
	}
}

func TestQueue_RemoveDefault(t *testing.T) {
	mq := MapQueue[int, string]{}

	if v, ok := mq.Remove(1); ok {
		t.Error("should not return a value")
	} else if v != "" {
		t.Errorf("should not return a value")
	}
}

func TestQueue_RemoveValue(t *testing.T) {
	mq := MapQueue[int, string]{}

	mq.Push(1, "a")
	mq.Push(2, "b")
	mq.Push(3, "c")

	if v, ok := mq.Remove(2); !ok {
		t.Error("should return a value")
	} else if v != "b" {
		t.Errorf("returned incorrect value: %v, expected: %v", v, "b")
	}

	if v, ok := mq.Remove(3); !ok {
		t.Error("should return a value")
	} else if v != "c" {
		t.Errorf("returned incorrect value: %v, expected: %v", v, "c")
	}

	if v, ok := mq.Remove(1); !ok {
		t.Error("should return a value")
	} else if v != "a" {
		t.Errorf("returned incorrect value: %v, expected: %v", v, "a")
	}
}
