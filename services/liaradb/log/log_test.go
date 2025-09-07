package log

import "testing"

func TestLog_Default(t *testing.T) {
	t.Parallel()

	l := &Log{}

	if h := l.HighWater(); h != 0 {
		t.Errorf("incorrect value: %v, expected: %v", h, 0)
	}

	if l := l.LowWater(); l != 0 {
		t.Errorf("incorrect value: %v, expected: %v", l, 0)
	}
}

func TestLog_Append(t *testing.T) {
	t.Parallel()

	l := &Log{}

	if lsn, err := l.Append([]byte{0, 1, 2, 3, 4, 5}); err != nil {
		t.Error(err)
	} else if lsn != 1 {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	if h := l.HighWater(); h != 1 {
		t.Errorf("incorrect value: %v, expected: %v", h, 1)
	}

	if l := l.LowWater(); l != 0 {
		t.Errorf("incorrect value: %v, expected: %v", l, 0)
	}
}

func TestLog_Flush(t *testing.T) {
	t.Parallel()

	t.Run("should flush", func(t *testing.T) {
		t.Parallel()

		l := &Log{}

		lsn1, err := l.Append([]byte{0, 1, 2, 3, 4, 5})
		if err != nil {
			t.Error(err)
		}

		_, err = l.Append([]byte{0, 1, 2, 3, 4, 5})
		if err != nil {
			t.Error(err)
		}

		if err := l.Flush(lsn1); err != nil {
			t.Error(err)
		}

		if h := l.HighWater(); h != 2 {
			t.Errorf("incorrect value: %v, expected: %v", h, 2)
		}

		if l := l.LowWater(); l != 1 {
			t.Errorf("incorrect value: %v, expected: %v", l, 1)
		}
	})

	t.Run("should not flush beyond HighWater", func(t *testing.T) {
		t.Parallel()

		l := &Log{}

		_, err := l.Append([]byte{0, 1, 2, 3, 4, 5})
		if err != nil {
			t.Error(err)
		}

		_, err = l.Append([]byte{0, 1, 2, 3, 4, 5})
		if err != nil {
			t.Error(err)
		}

		if err := l.Flush(10); err != nil {
			t.Error(err)
		}

		if h := l.HighWater(); h != 2 {
			t.Errorf("incorrect value: %v, expected: %v", h, 2)
		}

		if l := l.LowWater(); l != 2 {
			t.Errorf("incorrect value: %v, expected: %v", l, 1)
		}
	})
}
