package storage

import (
	"path"
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/filetesting"
)

func TestBufferManager(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testBufferManager)
}

func testBufferManager(t *testing.T) {
	b, bid := testCreateBuffer(t)

	if b.Dirty() {
		t.Error("should not be dirty")
	}

	if err := b.Load(bid); err != nil {
		t.Fatal(err)
	}

	want := []byte{1, 2, 3, 4, 5}

	if err := b.Add(want); err != nil {
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

	c := 0
	for i, err := range b.Items() {
		if err != nil {
			t.Error(err)
		}
		if !slices.Equal(want, i) {
			t.Errorf("value does not match: expected: %v, recieved: %v", want, i)
		}
		c++
	}

	if c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}
}

func testCreateBuffer(t *testing.T) (*Buffer, BlockID) {
	fsys := filetesting.NewDiskFileSystem(t)
	return NewBuffer(NewStorage(fsys, 2, 1024)),
		BlockID{FileName: path.Join(t.TempDir(), "testfile"), Position: 0}
}
