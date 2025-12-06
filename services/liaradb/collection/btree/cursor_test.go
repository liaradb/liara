package btree

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/file/filetesting"
	"github.com/liaradb/liaradb/storage"
)

func TestCursor_GetRoot_Default(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor)
}

func testCursor(t *testing.T) {
	s := createStorage(t, 2, 256)
	ctx := t.Context()

	n := "testfile"
	c := NewCursor[Key, any](s)
	r, err := c.GetRoot(ctx, n)
	if err != nil {
		t.Error(err)
	}

	if l := r.Level(); l != 0 {
		t.Errorf("incorrect level: %v, expected: %v", l, 0)
	}
}

func TestCursor_Insert(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_Insert)
}

func testCursor_Insert(t *testing.T) {
	s := createStorage(t, 2, 256)
	ctx := t.Context()

	n := "testfile"
	c := NewCursor[Key, any](s)
	r, err := c.GetRoot(ctx, n)
	if err != nil {
		t.Error(err)
	}

	if l := r.Level(); l != 0 {
		t.Errorf("incorrect level: %v, expected: %v", l, 0)
	}

	ln := NewLeafNode(r)
	_, ok := ln.Insert("a", NewRecordID(1, 2))
	if !ok {
		t.Error("should insert")
	}

	// TODO: Need to flush to disk
}

func createStorage(t *testing.T, max int, bs int64) *storage.Storage {
	fsys := filetesting.NewMockFileSystem(t, nil)
	s := storage.New(fsys, max, bs, t.TempDir())
	if err := s.Run(t.Context()); err != nil {
		t.Fatal(err)
	}

	return s
}
