package locktable

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/util/testing/testutil"
)

func TestConcurrencyMgr_SLock(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testConcurrencyMgr_SLock)
}

func testConcurrencyMgr_SLock(t *testing.T) {
	lt := New[int](1)
	lt.Run(t.Context())
	cm1 := NewConcurrencyMgr(lt)
	cm2 := NewConcurrencyMgr(lt)
	ctx := t.Context()

	if err := cm1.SLock(ctx, 0); err != nil {
		t.Fatal(err)
	}

	if err := cm1.SLock(ctx, 0); err != nil {
		t.Fatal(err)
	}

	if err := cm2.SLock(ctx, 0); err != nil {
		t.Fatal(err)
	}

	cm1.Release()
	cm2.Release()
}

func TestConcurrencyMgr_XLock(t *testing.T) {
	t.Parallel()

	testutil.Run(t, "should lock once", func(t *testing.T) {
		lt := New[int](1)
		lt.Run(t.Context())
		cm := NewConcurrencyMgr(lt)
		ctx := t.Context()

		if err := cm.XLock(ctx, 0); err != nil {
			t.Fatal(err)
		}
	})

	testutil.Run(t, "should not lock twice", func(t *testing.T) {
		lt := New[int](1)
		lt.Run(t.Context())
		cm1 := NewConcurrencyMgr(lt)
		cm2 := NewConcurrencyMgr(lt)
		ctx := t.Context()

		if err := cm1.XLock(ctx, 0); err != nil {
			t.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		if err := cm2.XLock(ctx, 0); err == nil {
			t.Fatal("should not lock")
		}
	})

	testutil.Run(t, "should not XLock after other SLock", func(t *testing.T) {
		lt := New[int](1)
		lt.Run(t.Context())
		cm1 := NewConcurrencyMgr(lt)
		cm2 := NewConcurrencyMgr(lt)
		ctx := t.Context()

		if err := cm1.SLock(ctx, 0); err != nil {
			t.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		if err := cm2.XLock(ctx, 0); err == nil {
			t.Fatal("should not lock")
		}
	})

	testutil.Run(t, "should upgrade lock", func(t *testing.T) {
		lt := New[int](1)
		lt.Run(t.Context())
		cm1 := NewConcurrencyMgr(lt)
		ctx := t.Context()

		if err := cm1.SLock(ctx, 0); err != nil {
			t.Fatal(err)
		}

		if err := cm1.XLock(ctx, 0); err != nil {
			t.Fatal("should upgrade lock")
		}
	})

	testutil.Run(t, "should lock after release", func(t *testing.T) {
		lt := New[int](1)
		lt.Run(t.Context())
		cm1 := NewConcurrencyMgr(lt)
		cm2 := NewConcurrencyMgr(lt)
		ctx := t.Context()

		if err := cm1.SLock(ctx, 0); err != nil {
			t.Fatal(err)
		}

		if err := cm1.XLock(ctx, 0); err != nil {
			t.Fatal("should upgrade lock")
		}

		cm1.Release()

		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		if err := cm2.SLock(ctx, 0); err != nil {
			t.Fatal(err)
		}
	})
}
