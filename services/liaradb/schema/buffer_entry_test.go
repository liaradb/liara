package storage

import (
	"path"
	"testing"

	"github.com/cardboardrobots/liaradb/file"
	"github.com/cardboardrobots/liaradb/storage"
)

func TestBufferEntry(t *testing.T) {
	b, close := testCreateBuffer(t)
	defer close()

	if err := b.Load(); err != nil {
		t.Fatal(err)
	}

	number := newUInt64Entry(0)

	var want uint64 = 12345

	if err := number.Set(b, want); err != nil {
		t.Fatal(err)
	}

	if v, err := number.Get(b); err != nil {
		t.Error(err)
	} else if v != want {
		t.Errorf("value does not match: expected: %v, recieved: %v", want, v)
	}
}

func testCreateBuffer(t *testing.T) (*storage.Buffer, func() error) {
	dir := t.TempDir()
	fs := &file.FileSystem{}

	bm := storage.NewBufferManager(fs)
	bid := storage.BlockID{FileName: path.Join(dir, "testfile"), Position: 0}
	return bm.Buffer(bid), fs.Close
}
