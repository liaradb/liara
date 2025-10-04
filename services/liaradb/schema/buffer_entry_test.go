package storage

import (
	"path"
	"testing"

	"github.com/liaradb/liaradb/filetesting"
	"github.com/liaradb/liaradb/storage"
)

func TestBufferEntry(t *testing.T) {
	t.Parallel()

	b, bid := testCreateBuffer(t)

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

func testCreateBuffer(t *testing.T) (*storage.Buffer, storage.BlockID) {
	fsys := filetesting.NewDiskFileSystem(t)
	return storage.NewBuffer(storage.NewStorage(fsys, 2, 1024)),
		storage.BlockID{FileName: path.Join(t.TempDir(), "testfile"), Position: 0}
}
