package storage

import (
	"context"
	"errors"
	"os"
	"path"
	"slices"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/filecache"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/queue"
	"github.com/liaradb/liaradb/util/testing/filetesting"
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

	if p := b.BlockID().Position(); p != 1 {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, p)
	}

	if p := b.Pins(); p != 1 {
		t.Errorf("incorrect pins: %v, expected: %v", p, 1)
	}

	if r := b.Reads(); r != 1 {
		t.Errorf("incorrect reads: %v, expected: %v", r, 1)
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

	if p := b0.BlockID().Position(); p != 2 {
		t.Errorf("incorrect result: expected %v, recieved %v", 2, p)
	}

	if p := b0.Pins(); p != 1 {
		t.Errorf("incorrect pins: %v, expected: %v", p, 1)
	}

	if r := b.Reads(); r != 1 {
		t.Errorf("incorrect reads: %v, expected: %v", r, 1)
	}

	b.Release()
	b0.Release()

	synctest.Wait()

	if p := b.Pins(); p != 0 {
		t.Errorf("incorrect pins: %v, expected: %v", p, 0)
	}

	if r := b.Reads(); r != 1 {
		t.Errorf("incorrect reads: %v, expected: %v", r, 1)
	}

	if p := b0.Pins(); p != 0 {
		t.Errorf("incorrect pins: %v, expected: %v", p, 0)
	}

	if r := b0.Reads(); r != 1 {
		t.Errorf("incorrect reads: %v, expected: %v", r, 1)
	}
}

func TestStorage_RequestBeforeRun(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_RequestBeforeRun)
}

func testStorage_RequestBeforeRun(t *testing.T) {
	s := Storage{}

	if _, err := s.Highwater(t.Context(), link.NewFileName("")); err == nil {
		t.Error("should return error")
	}

	if b, err := s.Request(t.Context(), link.BlockID{}); b != nil || err == nil {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, b.BlockID().Position())
	}

	if b, err := s.RequestCurrent(t.Context(), link.NewFileName("")); b != nil || err == nil {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, b.BlockID().Position())
	}

	if b, err := s.RequestNext(t.Context(), link.NewFileName("")); b != nil || err == nil {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, b.BlockID().Position())
	}
}

func TestStorage_Run(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_Run)
}

func testStorage_Run(t *testing.T) {
	fsys := filetesting.New(nil)
	dir := path.Join(t.TempDir(), "dir")
	s := New(fsys, &queue.MapQueue[link.BlockID, *Buffer]{}, 2, 1024, dir)

	if _, err := fsys.Stat(dir); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("dir should not exist")
	}

	ctx, cancel := context.WithCancel(t.Context())
	if err := s.Run(ctx); err != nil {
		t.Fatal(err)
	}

	if _, err := fsys.Stat(dir); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("dir should not exist")
	}

	ctx2, cancel2 := context.WithTimeout(t.Context(), 1*time.Second)
	defer cancel2()

	fn := link.NewFileName("testfile")
	b, err := s.Request(ctx2, fn.BlockID(1))
	if err != nil {
		t.Error(err)
	}

	if p := b.BlockID().Position(); p != 1 {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, p)
	}

	b.Release()

	synctest.Wait()

	cancel()

	if r, err := s.Request(ctx2, link.BlockID{}); r != nil || err == nil {
		t.Errorf("incorrect result: expected %v, recieved %v", 0, r.BlockID().Position())
	}
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

	data := []byte{1, 2}
	if _, err := b1.Write(data); err != nil {
		t.Fatal(err)
	}

	if !b1.Dirty() {
		t.Fatal("should be dirty")
	}

	// Release Buffer 1
	b1.Release()

	ctx2, cancel = context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// Request Buffer 2 again - available
	b2, err := s.Request(ctx2, bid2)
	if err != nil {
		t.Fatal(err)
	}

	if b2.Dirty() {
		t.Fatal("should have flushed")
	}

	synctest.Wait()

	if c := s.Count(); c != 2 {
		t.Errorf("incorrect number of Buffers.  Expected: %v, Recieved: %v", 2, c)
	}

	if c := s.CountPinned(); c != 2 {
		t.Errorf("incorrect number of Buffers.  Expected: %v, Recieved: %v", 2, c)
	}

	b0.Release()
	b2.Release()

	// Request Buffer 1 again
	b1, err = s.Request(ctx, bid1)
	if err != nil {
		t.Fatal(err)
	}

	result := make([]byte, 2)
	_, err = b1.Read(result)
	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}

	b1.Release()

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

	wg := sync.WaitGroup{}

	wg.Go(func() {
		b, err := s.Request(ctx, fn.BlockID(0))
		if err != nil {
			t.Error(err)
			return
		}

		b.Release()
	})

	wg.Go(func() {
		b, err := s.Request(ctx, fn.BlockID(1))
		if err != nil {
			t.Error(err)
			return
		}

		b.Release()
	})

	wg.Go(func() {
		b, err := s.Request(ctx, fn.BlockID(2))
		if err != nil {
			t.Error(err)
			return
		}

		b.Release()
	})

	wg.Wait()
	synctest.Wait()

	if c := s.Count(); c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}

	if c := s.CountPinned(); c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}
}

func TestStorage_Wait__NoLeak(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_Wait__NoLeak)
}

func testStorage_Wait__NoLeak(t *testing.T) {
	s := createStorageDelay(t, 1, 16, 1*time.Second)
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()

	fn := link.NewFileName("testfile")

	wg := sync.WaitGroup{}

	wg.Go(func() {
		b, err := s.Request(ctx, fn.BlockID(0))
		if err != nil {

			// Wait until file load is complete
			time.Sleep(5 * time.Second)

			if !errors.Is(err, context.Canceled) {
				t.Error(err)
			}

			return
		}

		b.Release()
	})

	wg.Wait()

	synctest.Wait()
}

func TestStorage_Reads(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_Reads)
}

func testStorage_Reads(t *testing.T) {
	s := createStorage(t, 1, 32)
	ctx := t.Context()

	fn := link.NewFileName("testfile")
	bid0 := fn.BlockID(0)
	bid1 := fn.BlockID(1)

	// Request Buffer 0
	b0, err := s.Request(ctx, bid0)
	if err != nil {
		t.Fatal(err)
	}

	if r := b0.Reads(); r != 1 {
		t.Errorf("incorrect reads: %v, expected: %v", r, 1)
	}

	b0.Release()

	// Request Buffer 0 again
	b0, err = s.Request(ctx, bid0)
	if err != nil {
		t.Fatal(err)
	}

	if r := b0.Reads(); r != 2 {
		t.Errorf("incorrect reads: %v, expected: %v", r, 2)
	}

	b0.Release()

	// Request Buffer 1
	b1, err := s.Request(ctx, bid1)
	if err != nil {
		t.Fatal(err)
	}

	if r := b1.Reads(); r != 1 {
		t.Errorf("incorrect reads: %v, expected: %v", r, 1)
	}

	b1.Release()

	synctest.Wait()
}

func TestStorage_Load__Error(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_Load__Error)
}

func testStorage_Load__Error(t *testing.T) {
	s, fsys := createStorageAndFileSystem(t, 1, 32, 0)
	ctx := t.Context()

	fn := link.NewFileName("testfile")
	bid0 := fn.BlockID(0)

	mfs := fsys.FSYS().(*filetesting.FileSystem)

	mfs.SetFail(true)

	// Request Buffer 0
	b0, err := s.Request(ctx, bid0)
	if err == nil {
		t.Fatal("should return error")
	}

	if r := b0.Reads(); r != 1 {
		t.Errorf("incorrect reads: %v, expected: %v", r, 1)
	}

	b0.Release()

	mfs.SetFail(false)

	// Request Buffer 0 again
	b1, err := s.Request(ctx, bid0)
	if err != nil {
		t.Fatal(err)
	}

	if r := b1.Reads(); r != 2 {
		t.Errorf("incorrect reads: %v, expected: %v", r, 2)
	}

	b1.Release()

	synctest.Wait()
}

func createStorage(t *testing.T, max int, bs int64) *Storage {
	t.Helper()

	fsys := filetesting.New(nil)
	return createStorageWithFileSystem(t, max, bs, fsys)
}

func createStorageDelay(t *testing.T, max int, bs int64, delay time.Duration) *Storage {
	t.Helper()

	fsys := filetesting.NewCacheDelay(nil, delay)
	return createStorageWithFileSystem(t, max, bs, fsys)
}

func createStorageAndFileSystem(t *testing.T, max int, bs int64, delay time.Duration) (*Storage, *filecache.Cache) {
	t.Helper()

	fsys := filetesting.NewCacheDelay(nil, delay)
	return createStorageWithFileSystem(t, max, bs, fsys), fsys
}

func createStorageWithFileSystem(t *testing.T, max int, bs int64, fsys file.FileSystem) *Storage {
	t.Helper()

	dir := t.TempDir()
	s := New(fsys, &queue.MapQueue[link.BlockID, *Buffer]{}, max, bs, dir)

	if d := s.Dir(); d != dir {
		t.Fatalf("incorrect dir: %v, expected: %v", d, dir)
	}

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
	s, fsys := createStorageAndFileSystem(t, 3, 16, 0)
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

	b2, err := s.Request(ctx, fn.BlockID(2))
	if err != nil {
		t.Fatal(err)
	}

	if b2.Dirty() {
		t.Error("should not be dirty")
	}

	b2.Release()

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

	if b2.Dirty() {
		t.Error("should not be dirty")
	}

	b0.Release()

	synctest.Wait()

	f, err := fsys.OpenFile(path.Join(s.Dir(), fn.String()))
	if err != nil {
		t.Fatal(err)
	}

	mf, ok := f.(*filecache.CacheFile).File.(*filetesting.File)
	if !ok {
		t.Fatal("incorrect type")
	}

	if wc := mf.WriteCount(); wc != 2 {
		t.Errorf("incorrect write count: %v, expected: %v", wc, 2)
	}
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

	// Should be clear
	for _, d := range b1.Raw() {
		if d != 0 {
			t.Errorf("incorrect value: %v, expected: %v", d, 0)
		}
	}

	if s := b1.Status(); s != BufferStatusLoaded {
		t.Errorf("incorrect status: %v, expected: %v", s, BufferStatusLoaded)
	}

	b1.Release()

	b2, err := s.RequestNext(ctx, fn)
	if err != nil {
		t.Fatal(err)
	}

	if p := b2.BlockID().Position().Value(); p != 2 {
		t.Errorf("incorrect position: %v, expected: %v", p, 2)
	}

	// Should be clear
	for _, d := range b2.Raw() {
		if d != 0 {
			t.Errorf("incorrect value: %v, expected: %v", d, 0)
		}
	}

	if s := b2.Status(); s != BufferStatusLoaded {
		t.Errorf("incorrect status: %v, expected: %v", s, BufferStatusLoaded)
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

func TestStorage_Highwater(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_Highwater)
}

func testStorage_Highwater(t *testing.T) {
	s := createStorage(t, 1, 16)
	ctx := t.Context()

	fn := link.NewFileName("testfile")

	if h, err := s.Highwater(ctx, fn); err != nil {
		t.Fatal(err)
	} else if h != fn.BlockID(0) {
		t.Errorf("incorrect highwater: %v, expected: %v", h, fn.BlockID(0))
	}

	if b, err := s.RequestNext(ctx, fn); err != nil {
		t.Fatal(err)
	} else {
		b.Release()
	}

	if h, err := s.Highwater(ctx, fn); err != nil {
		t.Fatal(err)
	} else if h != fn.BlockID(1) {
		t.Errorf("incorrect highwater: %v, expected: %v", h, fn.BlockID(1))
	}

	synctest.Wait()
}

func TestStorage_Release__NonBlocking(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_Release__NonBlocking)
}

func testStorage_Release__NonBlocking(t *testing.T) {
	s, fsys := createStorageAndFileSystem(t, 2, 16, 1*time.Second)
	ctx := t.Context()

	fn := link.NewFileName("testfile")

	wg := sync.WaitGroup{}

	step0 := make(chan struct{})
	step1 := make(chan struct{})

	// If blocking, forces waits
	// Request Block 0 -> success
	// Request Block 1 -> success
	// Lock FileSystem
	// Request Highwater -> waiting, blocks releases
	// Release Block 0 -> waiting
	// Release Block 1 -> waiting
	// Unlock FileSystem
	// Complete Highwater
	// Complete Release Block 0
	// Complete Release Block 1

	wg.Go(func() {
		b0, err := s.Request(ctx, fn.BlockID(0))
		if err != nil {
			t.Fatal(err)
		}

		b1, err := s.Request(ctx, fn.BlockID(1))
		if err != nil {
			t.Fatal(err)
		}

		close(step0)
		<-step1

		time.Sleep(2 * time.Second)

		b0.Release()
		b1.Release()
		fsys.FSYS().(*filetesting.FileSystem).UnLock()
	})

	wg.Go(func() {
		<-step0

		fsys.FSYS().(*filetesting.FileSystem).Lock()
		close(step1)

		_, err := s.Highwater(ctx, fn)
		if err != nil {
			t.Fatal(err)
		}
	})

	wg.Wait()
	synctest.Wait()
}
