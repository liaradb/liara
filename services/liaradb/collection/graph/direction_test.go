package graph

import "testing"

func TestDirection(t *testing.T) {
	t.Parallel()

	d := BiDirection

	if v := d.Size(); v != DirectionSize {
		t.Errorf("incorrect size: %v, expected: %v", v, DirectionSize)
	}
}
