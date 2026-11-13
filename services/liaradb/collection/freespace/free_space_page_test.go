package freespace

import (
	"fmt"
	"slices"
	"testing"
)

func TestMidIDPosition(t *testing.T) {
	for message, c := range map[string]struct {
		skip     bool
		id       MidID
		position int64
	}{
		"should handle 0": {
			id:       0,
			position: 1,
		},
		"should handle 1": {
			id:       1,
			position: 4002,
		},
		"should handle n": {
			id:       2,
			position: 8003,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}
			if p := MidIDPosition(c.id); p != c.position {
				t.Errorf("incorrect position: %v, expected: %v", p, c.position)
			}
		})
	}
}

func TestPageIDPosition(t *testing.T) {
	for message, c := range map[string]struct {
		skip     bool
		id       PageID
		position int64
	}{
		"should handle 0": {
			id:       0,
			position: 2,
		},
		"should handle 1": {
			id:       1,
			position: 3,
		},
		"should handle 4003": {
			id:       4003,
			position: 4006,
		},
		"should handle n": {
			id:       8004,
			position: 8008,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}
			if p := PageIDPosition(c.id); p != c.position {
				t.Errorf("incorrect position: %v, expected: %v", p, c.position)
			}
		})
	}
}

func TestNewPage(t *testing.T) {
	t.Parallel()
	t.Skip()

	p := NewPage(PageID(0), make([]byte, 4096))
	fmt.Println(1 << 1)
	fmt.Println(1 << 2)
	fmt.Println(p)

}

func TestNextIndex(t *testing.T) {
	for message, c := range map[string]struct {
		skip     bool
		position int16
		left     int16
		right    int16
	}{
		"should handle root": {
			position: 1,
			left:     2,
			right:    3,
		},
		"should handle 2nd left": {
			position: 2,
			left:     4,
			right:    5,
		},
		"should handle 2nd right": {
			position: 3,
			left:     6,
			right:    7,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			left := NextLeft(c.position)
			if left != c.left {
				t.Errorf("incorrect left: %v, expected: %v", left, c.left)
			}
			if right := left + 1; right != c.right {
				t.Errorf("incorrect right: %v, expected: %v", right, c.right)
			}
		})
	}
}

func TestGetPage(t *testing.T) {
	data := []byte{0,
		4,
		4, 2,
		3, 4, 0, 2,
	}
	i, ok := GetPage(data, 4, 4)
	if !ok {
		t.Error("should find value")
	}
	if i != 5 {
		t.Errorf("incorrect index: %v, expected: %v", i, 5)
	}
}

func TestSetPage(t *testing.T) {
	data := []byte{0,
		3,
		3, 2,
		3, 0, 0, 2,
	}
	SetPage(data, 5, 4, 4)
	want := []byte{0,
		4,
		4, 2,
		3, 4, 0, 2,
	}
	if !slices.Equal(data, want) {
		t.Errorf("incorrect dataa: %v, expected: %v", data, want)
	}
}
