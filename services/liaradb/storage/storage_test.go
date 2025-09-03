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
	s := Storage{}

	ctx := t.Context()
	s.Run(ctx)

	if r, _ := s.Request(ctx, 1); r.id != 1 {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, r)
	}

	if r, _ := s.Request(ctx, 2); r.id != 2 {
		t.Errorf("incorrect result: expected %v, recieved %v", 2, r)
	}
}

func TestStorage_RequestBeforeRun(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_RequestBeforeRun)
}

func testStorage_RequestBeforeRun(t *testing.T) {
	s := Storage{}

	if r, ok := s.Request(t.Context(), 0); r != nil || ok {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, r)
	}
}

func TestStorage_CancelRun(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testStorage_CancelRun)
}

func testStorage_CancelRun(t *testing.T) {
	s := Storage{}

	ctx, cancel := context.WithCancel(t.Context())
	s.Run(ctx)

	ctx2 := t.Context()

	if r, _ := s.Request(ctx2, 1); r.id != 1 {
		t.Errorf("incorrect result: expected %v, recieved %v", 1, r.id)
	}

	cancel()

	if r, ok := s.Request(ctx2, 0); r != nil || ok {
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
