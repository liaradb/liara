package storage

import (
	"path"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/filetesting"
)

func TestBufferList(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testBufferList)
}

func testBufferList(t *testing.T) {
	fsys := filetesting.NewDiskFileSystem(t)
	s := NewStorage(fsys, 2, 1024)
	bl := NewBufferList(s)

	ctx := t.Context()
	s.Run(ctx)

	n := path.Join(t.TempDir(), "testfile")

	b0, err := bl.Pin(ctx, BlockID{FileName: n, Position: 1})
	if err != nil {
		t.Fatal(err)
	}
	if p := b0.BlockID().Position; p != 1 {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, p)
	}
	if p := b0.Pins(); p != 1 {
		t.Errorf("incorrect pins: %v, expected: %v", p, 1)
	}

	b1, err := bl.Pin(ctx, BlockID{FileName: n, Position: 2})
	if err != nil {
		t.Fatal(err)
	}
	if p := b1.blockID.Position; p != 2 {
		t.Fatalf("incorrect result: expected %v, recieved %v", 2, p)
	}
	if p := b1.Pins(); p != 1 {
		t.Errorf("incorrect pins: %v, expected: %v", p, 1)
	}

	bl.Release()

	synctest.Wait()
	if p := b0.Pins(); p != 0 {
		t.Errorf("incorrect pins: %v, expected: %v", p, 0)
	}
	if p := b1.Pins(); p != 0 {
		t.Errorf("incorrect pins: %v, expected: %v", p, 0)
	}
}
