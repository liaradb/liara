package storage

import (
	"path"
	"testing"
)

func TestBufferManager(t *testing.T) {
	dir := t.TempDir()
	fs := &fileSystem{}
	defer fs.Close()

	bm := newBufferManager(fs)
	bid := BlockID{FileName: path.Join(dir, "testfile"), Position: 0}
	b := bm.Buffer(bid)

	if err := b.Load(); err != nil {
		t.Fatal(err)
	}

	var want uint64 = 12345

	if err := b.WriteUint64(want, 0); err != nil {
		t.Fatal(err)
	}

	if err := b.Flush(); err != nil {
		t.Fatal(err)
	}

	if err := b.Load(); err != nil {
		t.Fatal(err)
	}

	v, err := b.ReadUint64(0)
	if err != nil {
		t.Fatal(err)
	}
	if v != want {
		t.Errorf("value does not match: expected: %v, recieved: %v", want, v)
	}
}

func TestBufferEntry(t *testing.T) {
	dir := t.TempDir()
	fs := &fileSystem{}
	defer fs.Close()

	bm := newBufferManager(fs)
	bid := BlockID{FileName: path.Join(dir, "testfile"), Position: 0}
	b := bm.Buffer(bid)

	if err := b.Load(); err != nil {
		t.Fatal(err)
	}

	number := newUInt64Entry(0)

	var want uint64 = 12345
	if err := number.Set(b, want); err != nil {
		t.Fatal(err)
	}

	v, err := number.Get(b)
	if err != nil {
		t.Fatal(err)
	}
	if v != want {
		t.Errorf("value does not match: expected: %v, recieved: %v", want, v)
	}
}
