package testutil

import (
	"testing"
	"testing/synctest"
)

func Run(t *testing.T, m string, f func(*testing.T)) {
	t.Helper()
	t.Run(m, func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, f)
	})
}

func RunWait(t *testing.T, m string, f func(*testing.T)) {
	t.Helper()
	t.Run(m, func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			f(t)
			synctest.Wait()
		})
	})
}
