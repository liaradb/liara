package storage

import (
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/storage/link"
)

func TestBuffer_ReadWrite(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testBuffer_ReadWrite)
}

func testBuffer_ReadWrite(t *testing.T) {
	s := createStorage(t, 2, 16)
	ctx := t.Context()

	b, err := s.Request(ctx, link.NewBlockID(link.NewFileName(""), 0))
	if err != nil {
		t.Fatal(err)
	}

	data := []byte{1, 2, 3}

	if l, err := b.Write(data); err != nil {
		t.Fatal(err)
	} else if l != len(data) {
		t.Fatalf("incorrect length: %v, expected: %v", l, 3)
	}

	if l, err := b.Seek(0, 0); err != nil {
		t.Fatal(err)
	} else if l != 0 {
		t.Fatalf("incorrect length: %v, expected: %v", l, 0)
	}

	result := make([]byte, 3)
	if l, err := b.Read(result); err != nil {
		t.Fatal(err)
	} else if l != len(data) {
		t.Fatalf("incorrect length: %v, expected: %v", l, 3)
	}

	if !slices.Equal(result, data) {
		t.Error("Should copy")
	}

	b.Release()

	synctest.Wait()
}

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

	b1, err := s.Request(ctx, link.NewBlockID(link.NewFileName(""), 0))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := b0.Write([]byte{1, 2, 3}); err != nil {
		t.Fatal(err)
	}

	b1.Clone(b0)

	if !slices.Equal(b0.Raw(), b1.Raw()) {
		t.Error("Should copy")
	}

	b0.Release()
	b1.Release()

	synctest.Wait()
}
