package storage

import (
	"path"
	"testing"

	"github.com/cardboardrobots/liaradb/file"
)

func TestBufferManager(t *testing.T) {
	b, close := testCreateBuffer(t)
	defer close()

	if b.dirty {
		t.Error("should not be dirty")
	}

	if err := b.Load(); err != nil {
		t.Fatal(err)
	}

	var want uint64 = 12345

	if err := b.WriteUint64(want, 0); err != nil {
		t.Fatal(err)
	}

	if !b.dirty {
		t.Error("should be dirty")
	}

	if err := b.Flush(); err != nil {
		t.Fatal(err)
	}

	if b.dirty {
		t.Error("should not be dirty")
	}

	if err := b.Load(); err != nil {
		t.Fatal(err)
	}

	if b.dirty {
		t.Error("should not be dirty")
	}

	if v, err := b.ReadUint64(0); err != nil {
		t.Error(err)
	} else if v != want {
		t.Errorf("value does not match: expected: %v, recieved: %v", want, v)
	}
}

func testCreateBuffer(t *testing.T) (*Buffer, func() error) {
	dir := t.TempDir()
	fs := &file.FileSystem{}

	bm := NewBufferManager(fs)
	bid := BlockID{FileName: path.Join(dir, "testfile"), Position: 0}
	return bm.Buffer(bid), fs.Close
}
