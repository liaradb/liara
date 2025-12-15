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

	data := []LeafEntry{
		newLeafEntry(
			Key("a"),
			NewRecordID(1, 2)),
		newLeafEntry(
			Key("b"),
			NewRecordID(3, 4)),
		newLeafEntry(
			Key("c"),
			NewRecordID(5, 6)),
	}

	for _, e := range data {
		if err := c.Insert(ctx, n, e.key, e.recordID); err != nil {
			t.Error(err)
		}
		// TODO: Need to flush to disk
	}

	for _, e := range data {
		if rid, err := getRecordID(ctx, s, n, e.key); err != nil {
			t.Fatal(err)
		} else if rid != e.recordID {
			t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
		}
	}
}

func TestCursor_Insert__RootSplit(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCursor_Insert__RootSplit)
}

func testCursor_Insert__RootSplit(t *testing.T) {
	s := createStorage(t, 2, 62)
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

	data := []LeafEntry{
		newLeafEntry(
			Key("a"),
			NewRecordID(1, 2)),
		newLeafEntry(
			Key("b"),
			NewRecordID(3, 4)),
		newLeafEntry(
			Key("c"),
			NewRecordID(5, 6)),
	}

	for _, e := range data {
		if err := c.Insert(ctx, n, e.key, e.recordID); err != nil {
			t.Error(err)
		}
		// TODO: Need to flush to disk
	}

	for _, e := range data {
		if rid, err := getRecordID(ctx, s, n, e.key); err != nil {
			t.Error(err, e.key)
		} else if rid != e.recordID {
			t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
		}
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
