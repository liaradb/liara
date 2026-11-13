package freespace

import (
	"fmt"
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

			if left := NextIndex(c.position, false); left != c.left {
				t.Errorf("incorrect left: %v, expected: %v", left, c.left)
			}
			if right := NextIndex(c.position, true); right != c.right {
				t.Errorf("incorrect right: %v, expected: %v", right, c.right)
			}
		})
	}
}
