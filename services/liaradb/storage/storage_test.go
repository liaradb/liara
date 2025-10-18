package storage

import (
	"context"
	"path"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/filetesting"
	"github.com/liaradb/liaradb/raw"
)

func TestStorage(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage)
}

func testStorage(t *testing.T) {
	s := createStorage(t, 2, 16)
	ctx := t.Context()

	n := path.Join(t.TempDir(), "testfile")

	if b, err := s.Request(ctx, BlockID{FileName: n, Position: 1}); err != nil {
		t.Error(err)
	} else if b.blockID.Position != 1 {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, b.blockID.Position)
	}

	if b, err := s.Request(ctx, BlockID{FileName: n, Position: 2}); err != nil {
		t.Error(err)
	} else if b.blockID.Position != 2 {
		t.Errorf("incorrect result: expected %v, recieved %v", 2, b.blockID.Position)
	}
}

func TestStorage_RequestBeforeRun(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_RequestBeforeRun)
}

func testStorage_RequestBeforeRun(t *testing.T) {
	s := Storage{}

	if b, err := s.Request(t.Context(), BlockID{}); b != nil || err == nil {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, b.blockID.Position)
	}
}

func TestStorage_CancelRun(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_CancelRun)
}

func testStorage_CancelRun(t *testing.T) {
	fsys := filetesting.NewDiskFileSystem(t)
	s := NewStorage(fsys, 2, 1024)

	ctx, cancel := context.WithCancel(t.Context())
	s.Run(ctx)

	ctx2, cancel2 := context.WithTimeout(t.Context(), 1*time.Second)
	defer cancel2()

	n := path.Join(t.TempDir(), "testfile")
	if b, err := s.Request(ctx2, BlockID{FileName: n, Position: 1}); err != nil {
		t.Error(err)
	} else if b.blockID.Position != 1 {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, b.blockID.Position)
	}

	cancel()

	if r, err := s.Request(ctx2, BlockID{}); r != nil || err == nil {
		t.Errorf("incorrect result: expected %v, recieved %v", 0, r.blockID.Position)
	}
}

func TestStorage_Pinned(t *testing.T) {
	t.Parallel()

	s := createStorage(t, 2, 16)
	ctx := t.Context()

	n := path.Join(t.TempDir(), "testfile")
	bid := BlockID{FileName: n, Position: 0}

	b, err := s.Request(ctx, bid)
	if err != nil {
		t.Fatal(err)
	}

	if b.Dirty() {
		t.Error("should not be dirty")
	}

	if err = b.WriteUint64(1, 1); err != nil {
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
	s := createStorage(t, 2, 16)
	ctx := t.Context()

	n := path.Join(t.TempDir(), "testfile")
	bid0 := BlockID{FileName: n, Position: 0}
	bid1 := BlockID{FileName: n, Position: 1}
	bid2 := BlockID{FileName: n, Position: 2}

	b0, err := s.Request(ctx, bid0)
	if err != nil {
		t.Fatal(err)
	}

	if b0.Dirty() {
		t.Error("should not be dirty")
	}

	b1, err := s.Request(ctx, bid1)
	if err != nil {
		t.Fatal(err)
	}

	if b1.Dirty() {
		t.Error("should not be dirty")
	}

	ctx2, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err = s.Request(ctx2, bid2)
	if err != context.Canceled {
		t.Error("should be cancelled")
	}

	if err := b1.WriteUint64(12345, 0); err != nil {
		t.Fatal(err)
	}

	// TODO: Prove this is non-blocking
	b1.Release()

	ctx2, cancel = context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// TODO: How do we test that it flushed?
	_, err = s.Request(ctx2, bid2)
	if err != nil {
		t.Fatal(err)
	}

	if c := s.Count(); c != 2 {
		t.Errorf("incorrect number of Buffers.  Expected: %v, Recieved: %v", 2, c)
	}

	if c := s.CountPinned(); c != 2 {
		t.Errorf("incorrect number of Buffers.  Expected: %v, Recieved: %v", 2, c)
	}
}

func TestStorage_Append(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_Append)
}

func testStorage_Append(t *testing.T) {
	ctx := t.Context()
	s := createStorage(t, 1, 32)
	n := path.Join(t.TempDir(), "testfile")

	records := [][]byte{
		{1, 2},
		{3, 4},
		// {5, 6},
	}

	(func() {
		var bid BlockID
		var err error
		for _, r := range records {
			if bid, err = s.Append(ctx, n, raw.NewBufferFromSlice(r)); err != nil {
				t.Error(err)
			}
		}

		b, err := s.Request(ctx, bid)
		if err != nil {
			t.Fatal(err)
		}

		defer b.Release()

		err = b.Flush()
		if err != nil {
			t.Fatal(err)
		}
	})()

	c := 0
	for b, err := range s.Iterate(ctx, "testfile") {
		if err != nil {
			t.Error(err)
		}

		for i, err := range b.page.Items() {
			if err != nil {
				t.Error(err)
			}

			r := records[c]
			if !slices.Equal(i, r) {
				t.Errorf("incorrect record: %v, expected: %v", i, r)
			}

			c++
		}

		if c != len(records) {
			t.Errorf("incorrect count: %v, expected: %v", c, len(records))
		}
	}
}

func createStorage(t *testing.T, max int, bs int64) *Storage {
	fsys := filetesting.NewMockFileSystem(t, nil)
	s := NewStorage(fsys, max, bs)
	s.Run(t.Context())
	return s
}
