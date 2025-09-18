package storage

import (
	"path"
	"testing"

	"github.com/liaradb/liaradb/file/disk"
)

func TestBufferManager(t *testing.T) {
	t.Parallel()

	b, bid, close := testCreateBuffer(t)
	defer close()

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

func testCreateBuffer(t *testing.T) (*Buffer, BlockID, func() error) {
	dir := t.TempDir()
	fs := &disk.FileSystem{}

	return NewBuffer(NewStorage(fs, 2, 1024)),
		BlockID{FileName: path.Join(dir, "testfile"), Position: 0},
		fs.Close
}
