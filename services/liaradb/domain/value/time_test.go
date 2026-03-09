package value

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t.Parallel()

	now := time.Now()
	tm := NewTime(now)

	value := now.UTC()
	if v := tm.Value(); v != value {
		t.Errorf("incorrect value: %v, expected: %v", v, value)
	}
}
