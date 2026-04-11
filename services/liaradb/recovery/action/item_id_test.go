package action

import "testing"

func TestItemID(t *testing.T) {
	t.Parallel()

	want := "id"
	i := ItemID(want)

	if s := i.String(); s != want {
		t.Errorf("incorrect string: %v, expected: %v", i, want)
	}
}
