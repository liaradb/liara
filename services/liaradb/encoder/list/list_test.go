package list

import (
	"testing"
)

func TestList_Default(t *testing.T) {
	t.Run("default data", func(t *testing.T) {
		l := New([]byte{})

		if length := l.Length(); length != 0 {
			t.Errorf("incorrect length: %v, expected: %v", length, 0)
		}

		if s := l.Size(); s != 0 {
			t.Errorf("incorrect size: %v, expected: %v", s, 0)
		}
	})
}

func TestList_Push(t *testing.T) {
	t.Run("should push", func(t *testing.T) {
		l := New(make([]byte, 16))

		if length := l.Length(); length != 16 {
			t.Errorf("incorrect length: %v, expected: %v", length, 16)
		}

		if i, ok := l.Push(1); !ok {
			t.Error("should push")
		} else if i != 0 {
			t.Errorf("incorrect index: %v, expected: %v", i, 0)
		}

		if s := l.Size(); s != 1 {
			t.Errorf("incorrect size: %v, expected: %v", s, 1)
		}

		if i, ok := l.Push(2); !ok {
			t.Error("should push")
		} else if i != 1 {
			t.Errorf("incorrect index: %v, expected: %v", i, 0)
		}

		if s := l.Size(); s != 2 {
			t.Errorf("incorrect size: %v, expected: %v", s, 2)
		}

		if v, ok := l.Item(0); !ok {
			t.Errorf("should have value")
		} else if v != 1 {
			t.Errorf("incorrect value: %v, expected: %v", v, 1)
		}

		if v, ok := l.Item(1); !ok {
			t.Errorf("should have value")
		} else if v != 2 {
			t.Errorf("incorrect value: %v, expected: %v", v, 2)
		}
	})
}

func TestList_Pop(t *testing.T) {
	l := New(make([]byte, 16))

	if _, ok := l.Push(1); !ok {
		t.Error("should push")
	}

	if _, ok := l.Push(2); !ok {
		t.Error("should push")
	}

	if v, ok := l.Pop(); !ok {
		t.Error("should pop")
	} else if v != 2 {
		t.Errorf("incorrect value: %v, expected: %v", v, 2)
	}

	if s := l.Size(); s != 1 {
		t.Errorf("incorrect size: %v, expected: %v", s, 1)
	}

	if v, ok := l.Pop(); !ok {
		t.Error("should pop")
	} else if v != 1 {
		t.Errorf("incorrect value: %v, expected: %v", v, 1)
	}

	if s := l.Size(); s != 0 {
		t.Errorf("incorrect size: %v, expected: %v", s, 0)
	}

	if _, ok := l.Pop(); ok {
		t.Error("should not pop beyond empty")
	}
}
