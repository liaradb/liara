package iterator

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	t.Parallel()

	want := errors.New("error")
	count := 0
	for _, err := range Error[string](want) {
		if err != want {
			t.Errorf("incorrect error: %v, expected: %v", err, want)
		}
		count++
	}

	if count != 1 {
		t.Errorf("incorrect count: %v, expected: %v", count, 1)
	}
}
