package storage

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestStorage(t *testing.T) {
	s := NewStorage()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	s.Run(ctx)
	fmt.Printf("%v\n", s.Request(0))
	fmt.Printf("%v\n", s.Request(0))
	fmt.Printf("%v\n", s.Request(0))
	fmt.Printf("%v\n", s.Request(0))
	fmt.Printf("%v\n", s.Request(0))
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
