package locktable

import (
	"context"
	"testing"
	"time"

	"github.com/liaradb/liaradb/util/testutil"
)

func TestLockTable(t *testing.T) {
	t.Parallel()

	testutil.Run(t, "should XLock", func(t *testing.T) {
		lt := NewLockTable[int](1)
		lt.Run(t.Context())
		ctx := t.Context()

		if !lt.xLock(ctx, 0) {
			t.Fatal("should get lock")
		}
	})

	testutil.Run(t, "should not XLock twice", func(t *testing.T) {
		lt := NewLockTable[int](1)
		lt.Run(t.Context())
		ctx := t.Context()

		lt.xLock(ctx, 0)

		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		if lt.xLock(ctx, 0) {
			t.Fatal("should not get lock")
		}
	})

	testutil.Run(t, "should XLock after release XLock", func(t *testing.T) {
		lt := NewLockTable[int](1)
		lt.Run(t.Context())
		ctx := t.Context()

		lt.xLock(ctx, 0)
		go func() {
			time.Sleep(1 * time.Second)

			lt.release(0)
		}()

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if !lt.xLock(ctx, 0) {
			t.Fatal("should get lock")
		}
	})

	testutil.Run(t, "should XLock after failed XLock", func(t *testing.T) {
		lt := NewLockTable[int](1)
		lt.Run(t.Context())
		ctx := t.Context()

		lt.xLock(ctx, 0)

		ctx1, cancel1 := context.WithTimeout(ctx, 1*time.Second)
		defer cancel1()
		if lt.xLock(ctx1, 0) {
			t.Fatal("should not get lock")
		}

		lt.release(0)
		if !lt.xLock(ctx, 0) {
			t.Fatal("should get lock")
		}
	})

	testutil.Run(t, "should SLock", func(t *testing.T) {
		lt := NewLockTable[int](1)
		lt.Run(t.Context())
		ctx := t.Context()

		if !lt.sLock(ctx, 0) {
			t.Fatal("should get lock")
		}
	})

	testutil.Run(t, "should SLock twice", func(t *testing.T) {
		lt := NewLockTable[int](1)
		lt.Run(t.Context())
		ctx := t.Context()

		if !lt.sLock(ctx, 0) {
			t.Fatal("should get lock")
		}

		if !lt.sLock(ctx, 0) {
			t.Fatal("should get lock")
		}
	})

	testutil.Run(t, "should not XLock while multiple SLock", func(t *testing.T) {
		lt := NewLockTable[int](1)
		lt.Run(t.Context())
		ctx := t.Context()

		if !lt.sLock(ctx, 0) {
			t.Fatal("should get lock")
		}

		if !lt.sLock(ctx, 0) {
			t.Fatal("should get lock")
		}

		lt.release(0)

		if !lt.sLock(ctx, 0) {
			t.Fatal("should get lock")
		}

		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		if lt.xLock(ctx, 0) {
			t.Fatal("should not get XLock")
		}
	})

	testutil.Run(t, "should not SLock after XLock", func(t *testing.T) {
		lt := NewLockTable[int](1)
		lt.Run(t.Context())
		ctx := t.Context()

		lt.xLock(ctx, 0)

		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		if lt.sLock(ctx, 0) {
			t.Fatal("should not get lock")
		}
	})

	testutil.Run(t, "should SLock after release XLock", func(t *testing.T) {
		lt := NewLockTable[int](1)
		lt.Run(t.Context())
		ctx := t.Context()

		lt.xLock(ctx, 0)
		go func() {
			time.Sleep(1 * time.Second)

			lt.release(0)
		}()

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if !lt.sLock(ctx, 0) {
			t.Fatal("should get lock")
		}
	})

	testutil.Run(t, "should SLock after failed XLock", func(t *testing.T) {
		lt := NewLockTable[int](1)
		lt.Run(t.Context())
		ctx := t.Context()

		lt.xLock(ctx, 0)

		ctx1, cancel1 := context.WithTimeout(ctx, 1*time.Second)
		defer cancel1()
		if lt.xLock(ctx1, 0) {
			t.Fatal("should not get lock")
		}

		lt.release(0)
		if !lt.sLock(ctx, 0) {
			t.Fatal("should get lock")
		}
	})
}
