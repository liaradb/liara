package domain

import "testing"

func TestVersion_Increment(t *testing.T) {
	v := Version(0)
	v.Increment()

	if v != 1 {
		t.Error("wrong value")
	}
}
