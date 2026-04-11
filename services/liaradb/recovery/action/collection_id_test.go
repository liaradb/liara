package action

import "testing"

func TestCollectionId(t *testing.T) {
	t.Parallel()

	wantA := "a"
	a := CollectionID(wantA)
	if s := a.String(); s != wantA {
		t.Errorf("incorrect string: %v, expected: %v", s, wantA)
	}

	wantB := "b"
	b := CollectionID(wantB)
	if s := b.String(); s != wantB {
		t.Errorf("incorrect string: %v, expected: %v", s, wantB)
	}
}
