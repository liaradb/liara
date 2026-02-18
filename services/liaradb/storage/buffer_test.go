package storage

import (
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/storage/link"
)

func TestBuffer_Clone(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testBuffer_Clone)
}

func testBuffer_Clone(t *testing.T) {
	s := createStorage(t, 2, 16)
	ctx := t.Context()

	b0, err := s.Request(ctx, link.NewBlockID(link.NewFileName(""), 0))
	if err != nil {
		t.Fatal(err)
	}

	defer b0.Release()

	b1, err := s.Request(ctx, link.NewBlockID(link.NewFileName(""), 0))
	if err != nil {
		t.Fatal(err)
	}

	defer b1.Release()

	if _, err := b0.Write([]byte{1, 2, 3}); err != nil {
		t.Fatal(err)
	}

	b1.Clone(b0)

	if !slices.Equal(b0.Raw(), b1.Raw()) {
		t.Error("Should copy")
	}

	synctest.Wait()
}
