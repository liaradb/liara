package int32list

import "testing"

func TestInt32List_Default(t *testing.T) {
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

func TestList_GetSet(t *testing.T) {
	t.Run("should set and get", func(t *testing.T) {
		l := New(make([]byte, 16))

		if length := l.Length(); length != 16 {
			t.Errorf("incorrect length: %v, expected: %v", length, 16)
		}

		if s := l.Size(); s != 4 {
			t.Errorf("incorrect size: %v, expected: %v", s, 1)
		}

		for i := range int32(4) {
			if ok := l.Set(i, (i+1)*11); !ok {
				t.Error("should set value")
			}
		}

		if ok := l.Set(4, 55); ok {
			t.Error("should not set value beyond size")
		}

		for i := range int32(4) {
			want := (i + 1) * 11
			if v, ok := l.Get(i); !ok {
				t.Error("should set value")
			} else if v != want {
				t.Errorf("incorrect value: %v, expected: %v", v, want)
			}
		}

		if _, ok := l.Get(4); ok {
			t.Error("should not set value beyond size")
		}
	})
}
