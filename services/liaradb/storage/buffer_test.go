package storage

import (
	"path"
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/filetesting"
	"github.com/liaradb/liaradb/storage/page"
)

func TestBuffer(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testBuffer)
}

func testBuffer(t *testing.T) {
	b, bid := testCreateBuffer(t)

	if b.Dirty() {
		t.Error("should not be dirty")
	}

	if err := b.Load(bid); err != nil {
		t.Fatal(err)
	}

	want := [][]byte{{1, 2, 3, 4, 5}}

	if err := b.Add(want[0]); err != nil {
		t.Fatal(err)
	}

	if !b.Dirty() {
		t.Error("should be dirty")
	}

	if err := b.Flush(); err != nil {
		t.Fatal(err)
	}

	if b.Dirty() {
		t.Error("should not be dirty")
	}

	if err := b.Flush(); err == nil {
		t.Fatal("should not flush clean buffers")
	}

	if err := b.Load(bid); err != nil {
		t.Fatal(err)
	}

	if b.Dirty() {
		t.Error("should not be dirty")
	}

	result := make([]page.Item, 0)

	for i, err := range b.Items() {
		if err != nil {
			t.Error(err)
		}

		result = append(result, i)
	}

	if !slices.EqualFunc(result, want, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func testCreateBuffer(t *testing.T) (*Buffer, BlockID) {
	fsys := filetesting.NewDiskFileSystem(t)
	return NewBuffer(NewStorage(fsys, 2, 1024)),
		BlockID{FileName: path.Join(t.TempDir(), "testfile"), Position: 0}
}
