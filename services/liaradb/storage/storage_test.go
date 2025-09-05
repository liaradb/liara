package storage

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/cardboardrobots/liaradb/file"
)

func TestStorage(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage)
}

func testStorage(t *testing.T) {
	s := Storage{}

	ctx := t.Context()
	s.Run(ctx, NewBufferManager(&file.FileSystem{}))

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
	s := Storage{}

	ctx, cancel := context.WithCancel(t.Context())
	s.Run(ctx, NewBufferManager(&file.FileSystem{}))

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
	s := Storage{}

	ctx := t.Context()
	s.Run(ctx, NewBufferManager(&file.FileSystem{}))

	n := path.Join(t.TempDir(), "testfile")
	bid := BlockID{FileName: n, Position: 0}

	b, err := s.Request(ctx, bid)
	if err != nil {
		t.Fatal(err)
	}

	if b.dirty {
		t.Error("should not be dirty")
	}

	err = b.WriteUint64(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	if !b.dirty {
		t.Error("should be dirty")
	}

	b, err = s.Request(ctx, bid)
	if err != nil {
		t.Fatal(err)
	}

	if !b.dirty {
		t.Error("should be dirty")
	}
}

func TestMultipleWriter(t *testing.T) {
	t.Skip()
	f, _ := os.OpenFile("multiple.text", os.O_RDWR|os.O_CREATE, 0644)
	wg := sync.WaitGroup{}
	for i := range 100 {
		wg.Go(func() {
			fmt.Fprintf(f, "abcdefg: %v\n", i)
		})
	}
	wg.Wait()
}
