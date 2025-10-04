package sharedpool

import (
	"context"
	"testing"
	"testing/synctest"
	"time"
)

func TestSharedPool(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testSharedPool)
}

func testSharedPool(t *testing.T) {
	sp := NewSharedPool[string, *testItem](2)
	t.Cleanup(sp.Close)

	for i := range 2 {
		sp.Add(&testItem{id: i})
	}

	ctx := context.Background()

	if a, ok := sp.Request(ctx, "a"); !ok || a == nil {
		t.Error("should get value")
	}

	var b *testItem
	var ok bool
	if b, ok = sp.Request(ctx, "b"); !ok || b == nil {
		t.Error("should get value")
	}

	ctx2, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if c, ok := sp.Request(ctx2, "c"); ok || c != nil {
		t.Error("should not get value")
	}

	sp.Release(b)

	if c, ok := sp.Request(ctx, "c"); !ok || c == nil {
		t.Error("should get value")
	}
}

type testItem struct {
	id int
}

func (ti *testItem) Id() int {
	return ti.id
}

func (ti *testItem) Block() (string, bool) {
	return "", true
}

func (ti *testItem) Pin() {

}

func (ti *testItem) ReplaceBlock(string) error {
	return nil
}
