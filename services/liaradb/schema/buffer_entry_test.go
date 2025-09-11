package storage

import (
	"path"
	"testing"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/storage"
)

func TestBufferEntry(t *testing.T) {
	t.Parallel()

	b, bid, close := testCreateBuffer(t)
	defer close()

	if err := b.Load(bid); err != nil {
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

func testCreateBuffer(t *testing.T) (*storage.Buffer, storage.BlockID, func() error) {
	dir := t.TempDir()
	fs := &file.FileSystem{}

	return storage.NewBuffer(storage.NewStorage(fs, 2, 1024)),
		storage.BlockID{FileName: path.Join(dir, "testfile"), Position: 0},
		fs.Close
}
