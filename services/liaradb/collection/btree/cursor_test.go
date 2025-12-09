package btree

import (
	"context"
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
	c := NewCursor(s)
	r, err := c.GetRoot(ctx, n)
	if err != nil {
		t.Error(err)
	}

	if l := r.Level(); l != 0 {
		t.Errorf("incorrect level: %v, expected: %v", l, 0)
	}
}

func TestCursor_Insert__Root(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_Insert__Root)
}

func testCursor_Insert__Root(t *testing.T) {
	s := createStorage(t, 2, 256)
	ctx := t.Context()

	n := "testfile"
	c := NewCursor(s)
	r, err := c.GetRoot(ctx, n)
	if err != nil {
		t.Error(err)
	}

	if l := r.Level(); l != 0 {
		t.Errorf("incorrect level: %v, expected: %v", l, 0)
	}

	wantRID := NewRecordID(1, 2)
	if err := c.Insert(ctx, n, "a", wantRID); err != nil {
		t.Error(err)
	}

	// TODO: Need to flush to disk
	if rid, err := getRecordID(ctx, s, n, "a"); err != nil {
		t.Fatal(err)
	} else if rid != wantRID {
		t.Errorf("incorrect record id: %v, expected: %v", rid, wantRID)
	}
}

func getRecordID(
	ctx context.Context,
	s *storage.Storage,
	name string,
	key Key,
) (RecordID, error) {
	c := NewCursor(s)
	return c.Search(ctx, name, key)
}

func createStorage(t *testing.T, max int, bs int64) *storage.Storage {
	fsys := filetesting.NewMockFileSystem(t, nil)
	s := storage.New(fsys, max, bs, t.TempDir())
	if err := s.Run(t.Context()); err != nil {
		t.Fatal(err)
	}

	return s
}
