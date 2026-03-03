package storage

import (
	"context"
	"slices"
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

	bid0 := fn.BlockID(1)
	b, err := s.Request(ctx, bid0)
	if err != nil {
		t.Error(err)
	}

	if i := b.BlockID(); i != bid0 {
		t.Errorf("incorrect block id: %v, expected: %v", i, bid0)
	}

	if p := b.blockID.Position(); p != 1 {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, p)
	}

	if p := b.Pins(); p != 1 {
		t.Errorf("incorrect pins: %v, expected: %v", p, 1)
	}

	if s := b.Size(); s != 16 {
		t.Errorf("incorrect size: %v, expected: %v", s, 16)
	}

	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	if n, err := b.Write(data); err != nil {
		t.Fatal(err)
	} else if n != 16 {
		t.Fatalf("incorrect length: %v, expected: %v", n, 16)
	}

	if bytes := b.Raw(); !slices.Equal(bytes, data) {
		t.Errorf("incorrect raw value: %v, expected: %v", bytes, data)
	}

	bid1 := fn.BlockID(2)
	b0, err := s.Request(ctx, bid1)
	if err != nil {
		t.Error(err)
	}

	if i := b0.BlockID(); i != bid1 {
		t.Errorf("incorrect block id: %v, expected: %v", i, bid1)
	}

	if p := b0.blockID.Position(); p != 2 {
		t.Errorf("incorrect result: expected %v, recieved %v", 2, p)
	}

	if p := b0.Pins(); p != 1 {
		t.Errorf("incorrect pins: %v, expected: %v", p, 1)
	}

	b.Release()
	b0.Release()

	synctest.Wait()

	if p := b.Pins(); p != 0 {
		t.Errorf("incorrect pins: %v, expected: %v", p, 0)
	}

	if p := b0.Pins(); p != 0 {
		t.Errorf("incorrect pins: %v, expected: %v", p, 0)
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

	if b, err := s.RequestCurrent(t.Context(), link.NewFileName("")); b != nil || err == nil {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, b.blockID.Position())
	}

	if b, err := s.RequestNext(t.Context(), link.NewFileName("")); b != nil || err == nil {
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
	b, err := s.Request(ctx2, fn.BlockID(1))
	if err != nil {
		t.Error(err)
	}

	if p := b.blockID.Position(); p != 1 {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, p)
	}

	cancel()

	if r, err := s.Request(ctx2, link.BlockID{}); r != nil || err == nil {
		t.Errorf("incorrect result: expected %v, recieved %v", 0, r.blockID.Position())
	}

	b.Release()

	synctest.Wait()
}

func TestStorage_Pinned(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_Pinned)
}

func testStorage_Pinned(t *testing.T) {
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

	b0, err := s.Request(ctx, bid)
	if err != nil {
		t.Fatal(err)
	}

	if !b0.Dirty() {
		t.Error("should be dirty")
	}

	b.Release()
	b0.Release()

	synctest.Wait()
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
	b3, err := s.Request(ctx2, bid2)
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

	b0.Release()
	b3.Release()

	synctest.Wait()
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
	return createStorageWithFileSystem(t, max, bs, fsys)
}

func createStorageWithFileSystem(t *testing.T, max int, bs int64, fsys file.FileSystem) *Storage {
	t.Helper()

	s := New(fsys, max, bs, t.TempDir())
	if err := s.Run(t.Context()); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		synctest.Wait()

		if p := s.CountPinned(); p != 0 {
			t.Errorf("incorrect pin count: %v, expected: %v", p, 0)
		}
	})

	return s
}

func TestStorage_FlushAll(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_FlushAll)
}

func testStorage_FlushAll(t *testing.T) {
	s := createStorage(t, 2, 16)
	ctx := t.Context()

	fn := link.NewFileName("testfile")
	b0, err := s.Request(ctx, fn.BlockID(0))
	if err != nil {
		t.Fatal(err)
	}

	b0.Write([]byte{1})
	if !b0.Dirty() {
		t.Error("should be dirty")
	}

	b1, err := s.Request(ctx, fn.BlockID(1))
	if err != nil {
		t.Fatal(err)
	}

	b1.Write([]byte{1})
	if !b1.Dirty() {
		t.Error("should be dirty")
	}

	b1.Release()

	synctest.Wait()

	if err := s.FlushAll(); err != nil {
		t.Fatal(err)
	}

	if b0.Dirty() {
		t.Error("should not be dirty")
	}

	if b1.Dirty() {
		t.Error("should not be dirty")
	}

	b0.Release()

	synctest.Wait()
}

func TestStorage_Request__AlreadyLoaded(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_Request__AlreadyLoaded)
}

func testStorage_Request__AlreadyLoaded(t *testing.T) {
	s := createStorage(t, 1, 16)
	ctx := t.Context()

	fn := link.NewFileName("testfile")
	b0, err := s.Request(ctx, fn.BlockID(0))
	if err != nil {
		t.Fatal(err)
	}

	b0.Release()

	synctest.Wait()

	b1, err := s.Request(ctx, fn.BlockID(0))
	if err != nil {
		t.Fatal(err)
	}

	b1.Release()

	synctest.Wait()
}

func TestStorage_RequestCurrent_RequestNext(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_RequestCurrent_RequestNext)
}

func testStorage_RequestCurrent_RequestNext(t *testing.T) {
	s := createStorage(t, 1, 16)
	ctx := t.Context()

	fn := link.NewFileName("testfile")

	b0, err := s.RequestCurrent(ctx, fn)
	if err != nil {
		t.Fatal(err)
	}

	if p := b0.BlockID().Position().Value(); p != 0 {
		t.Errorf("incorrect position: %v, expected: %v", p, 0)
	}

	b0.Release()

	b1, err := s.RequestNext(ctx, fn)
	if err != nil {
		t.Fatal(err)
	}

	if p := b1.BlockID().Position().Value(); p != 1 {
		t.Errorf("incorrect position: %v, expected: %v", p, 1)
	}

	b1.Release()

	b2, err := s.RequestNext(ctx, fn)
	if err != nil {
		t.Fatal(err)
	}

	if p := b2.BlockID().Position().Value(); p != 2 {
		t.Errorf("incorrect position: %v, expected: %v", p, 2)
	}

	b2.Release()

	b3, err := s.RequestCurrent(ctx, fn)
	if err != nil {
		t.Fatal(err)
	}

	if p := b3.BlockID().Position().Value(); p != 2 {
		t.Errorf("incorrect position: %v, expected: %v", p, 2)
	}

	b3.Release()

	synctest.Wait()
}
