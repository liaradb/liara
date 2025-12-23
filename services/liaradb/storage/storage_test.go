package storage

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/filetesting"
	"github.com/liaradb/liaradb/storage/link"
)

func TestStorage(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage)
}

func testStorage(t *testing.T) {
	s := createStorage(t, 2, 16)
	ctx := t.Context()

	fn := link.NewFileName("testfile")

	if b, err := s.Request(ctx, fn.BlockID(1)); err != nil {
		t.Error(err)
	} else if p := b.blockID.Position(); p != 1 {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, p)
	}

	if b, err := s.Request(ctx, fn.BlockID(2)); err != nil {
		t.Error(err)
	} else if p := b.blockID.Position(); p != 2 {
		t.Errorf("incorrect result: expected %v, recieved %v", 2, p)
	}
}

func TestStorage_RequestBeforeRun(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_RequestBeforeRun)
}

func testStorage_RequestBeforeRun(t *testing.T) {
	s := Storage{}

	if b, err := s.Request(t.Context(), link.BlockID{}); b != nil || err == nil {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, b.blockID.Position())
	}
}

func TestStorage_CancelRun(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_CancelRun)
}

func testStorage_CancelRun(t *testing.T) {
	fsys := filetesting.NewMockFileSystem(t, nil)
	s := New(fsys, 2, 1024, t.TempDir())

	ctx, cancel := context.WithCancel(t.Context())
	if err := s.Run(ctx); err != nil {
		t.Fatal(err)
	}

	ctx2, cancel2 := context.WithTimeout(t.Context(), 1*time.Second)
	defer cancel2()

	fn := link.NewFileName("testfile")
	if b, err := s.Request(ctx2, fn.BlockID(1)); err != nil {
		t.Error(err)
	} else if p := b.blockID.Position(); p != 1 {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, p)
	}

	cancel()

	if r, err := s.Request(ctx2, link.BlockID{}); r != nil || err == nil {
		t.Errorf("incorrect result: expected %v, recieved %v", 0, r.blockID.Position())
	}
}

func TestStorage_Pinned(t *testing.T) {
	t.Parallel()

	s := createStorage(t, 2, 32)
	ctx := t.Context()

	fn := link.NewFileName("testfile")
	bid := fn.BlockID(0)

	b, err := s.Request(ctx, bid)
	if err != nil {
		t.Fatal(err)
	}

	if b.Dirty() {
		t.Error("should not be dirty")
	}

	if _, err := b.Write([]byte{1}); err != nil {
		t.Fatal(err)
	}

	if !b.Dirty() {
		t.Error("should be dirty")
	}

	b, err = s.Request(ctx, bid)
	if err != nil {
		t.Fatal(err)
	}

	if !b.Dirty() {
		t.Error("should be dirty")
	}
}

func TestStorage_Flush(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_Flush)
}

func testStorage_Flush(t *testing.T) {
	s := createStorage(t, 2, 32)
	ctx := t.Context()

	fn := link.NewFileName("testfile")
	bid0 := fn.BlockID(0)
	bid1 := fn.BlockID(1)
	bid2 := fn.BlockID(2)

	// Request Buffer 0
	b0, err := s.Request(ctx, bid0)
	if err != nil {
		t.Fatal(err)
	}

	if b0.Dirty() {
		t.Error("should not be dirty")
	}

	// Request Buffer 1
	b1, err := s.Request(ctx, bid1)
	if err != nil {
		t.Fatal(err)
	}

	if b1.Dirty() {
		t.Error("should not be dirty")
	}

	ctx2, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// Request Buffer 2 - none available
	_, err = s.Request(ctx2, bid2)
	if err != context.Canceled {
		t.Error("should be cancelled")
	}

	if _, err := b1.Write([]byte{1, 2}); err != nil {
		t.Fatal(err)
	}

	// Release Buffer 1
	// TODO: Prove this is non-blocking
	b1.Release()

	ctx2, cancel = context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// Request Buffer 2 again - available
	// TODO: How do we test that it flushed?
	_, err = s.Request(ctx2, bid2)
	if err != nil {
		t.Fatal(err)
	}

	synctest.Wait()

	if c := s.Count(); c != 2 {
		t.Errorf("incorrect number of Buffers.  Expected: %v, Recieved: %v", 2, c)
	}

	if c := s.CountPinned(); c != 2 {
		t.Errorf("incorrect number of Buffers.  Expected: %v, Recieved: %v", 2, c)
	}
}

func TestStorage_Wait(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_Wait)
}

func testStorage_Wait(t *testing.T) {
	s := createStorage(t, 1, 16)
	ctx := t.Context()

	fn := link.NewFileName("testfile")

	go func() {
		b, err := s.Request(ctx, fn.BlockID(0))
		if err != nil {
			t.Error(err)
			return
		}

		b.status = BufferStatusDirty

		b.Release()
	}()

	go func() {
		b, err := s.Request(ctx, fn.BlockID(1))
		if err != nil {
			t.Error(err)
			return
		}

		b.status = BufferStatusDirty

		b.Release()
	}()

	go func() {
		b, err := s.Request(ctx, fn.BlockID(2))
		if err != nil {
			t.Error(err)
			return
		}

		b.status = BufferStatusDirty

		b.Release()
	}()

	synctest.Wait()

	if c := s.Count(); c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}

	if c := s.CountPinned(); c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}
}

// TODO: Test with noPin true
func createStorage(t *testing.T, max int, bs int64) *Storage {
	t.Helper()

	fsys := filetesting.NewMockFileSystem(t, nil)
	return createStorageWithFileSystem(t, max, bs, fsys, false)
}

func createStorageWithFileSystem(t *testing.T, max int, bs int64, fsys file.FileSystem, noPin bool) *Storage {
	t.Helper()

	s := New(fsys, max, bs, t.TempDir())
	if err := s.Run(t.Context()); err != nil {
		t.Fatal(err)
	}

	if noPin {
		t.Cleanup(func() {
			synctest.Wait()

			if p := s.CountPinned(); p != 0 {
				t.Errorf("incorrect pin count: %v, expected: %v", p, 0)
			}
		})
	}

	return s
}
