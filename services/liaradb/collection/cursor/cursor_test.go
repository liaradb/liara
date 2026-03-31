package cursor

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestCursor(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor)
}

func testCursor(t *testing.T) {
	s := storagetesting.CreateStorage(t, 3, 256)
	ctx := t.Context()
	fn := link.NewFileName("testfile")

	b0, err := s.Request(ctx, fn.BlockID(0))
	if err != nil {
		t.Fatal(err)
	}

	b1, err := s.Request(ctx, fn.BlockID(1))
	if err != nil {
		t.Fatal(err)
	}

	b2, err := s.Request(ctx, fn.BlockID(2))
	if err != nil {
		t.Fatal(err)
	}

	c := New(b0, b1, b2)

	c.Release()

	synctest.Wait()
}
