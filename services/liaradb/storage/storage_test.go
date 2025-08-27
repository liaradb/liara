package storage

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"testing/synctest"
)

func TestStorage(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage)
}

func testStorage(t *testing.T) {
	s := NewStorage()

	ctx := t.Context()
	s.Run(ctx)

	if r := s.Request(ctx, 0); r != 1 {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, r)
	}

	if r := s.Request(ctx, 0); r != 2 {
		t.Errorf("incorrect result: expected %v, recieved %v", 2, r)
	}
}

func TestStorage_CancelRun(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_CancelRun)
}

func testStorage_CancelRun(t *testing.T) {
	s := NewStorage()

	ctx, cancel := context.WithCancel(t.Context())
	s.Run(ctx)

	ctx2 := t.Context()

	if r := s.Request(ctx2, 0); r != 1 {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, r)
	}

	cancel()

	if r := s.Request(ctx2, 0); r != 0 {
		t.Errorf("incorrect result: expected %v, recieved %v", 0, r)
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
