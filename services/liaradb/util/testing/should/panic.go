package should

import "testing"

func Panic(t *testing.T, m string, p func()) {
	t.Helper()

	defer assertPanic(t, m)

	p()
}

func assertPanic(t *testing.T, m string) {
	t.Helper()

	if r := recover(); r == nil {
		t.Error(m)
	}
}
