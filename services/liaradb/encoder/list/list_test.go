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

func TestList_Append(t *testing.T) {
	t.Run("should append", func(t *testing.T) {
		l := New(make([]byte, 16))

		if length := l.Length(); length != 16 {
			t.Errorf("incorrect length: %v, expected: %v", length, 16)
		}

		if i, err := l.Append(1); err != nil {
			t.Error(err)
		} else if i != 0 {
			t.Errorf("incorrect index: %v, expected: %v", i, 0)
		}

		if s := l.Size(); s != 1 {
			t.Errorf("incorrect size: %v, expected: %v", s, 1)
		}

		if i, err := l.Append(2); err != nil {
			t.Error(err)
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
