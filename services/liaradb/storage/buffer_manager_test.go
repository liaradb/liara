package storage

import (
	"path"
	"testing"

	"github.com/liaradb/liaradb/filetesting"
)

func TestBufferManager(t *testing.T) {
	t.Parallel()

	b, bid := testCreateBuffer(t)

	if b.Dirty() {
		t.Error("should not be dirty")
	}

	if err := b.Load(bid); err != nil {
		t.Fatal(err)
	}

	var want uint64 = 12345

	if err := b.WriteUint64(want, 0); err != nil {
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

	if v, err := b.ReadUint64(0); err != nil {
		t.Error(err)
	} else if v != want {
		t.Errorf("value does not match: expected: %v, recieved: %v", want, v)
	}
}

func testCreateBuffer(t *testing.T) (*Buffer, BlockID) {
	fsys := filetesting.NewDiskFileSystem(t)
	return NewBuffer(NewStorage(fsys, 2, 1024)),
		BlockID{FileName: path.Join(t.TempDir(), "testfile"), Position: 0}
}
