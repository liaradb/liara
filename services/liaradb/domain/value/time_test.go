package value

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	tm := NewTime(now)

	if v := tm.Value(); v != now {
		t.Errorf("incorrect value: %v, expected: %v", v, now)
	}
}
