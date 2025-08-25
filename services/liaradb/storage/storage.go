package storage

import (
	"context"
)

type Storage struct {
	count int
	file  file
	in    chan *request
}

func NewStorage() *Storage {
	return &Storage{
		in: make(chan *request),
	}
}

type request struct {
	value int
	out   chan int
}

func (s *Storage) Run(ctx context.Context) {
	go s.run(ctx)
}

func (s *Storage) run(ctx context.Context) {
	for {
		select {
		case r := <-s.in:
			r.out <- s.request(r.value)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Storage) Request(value int) int {
	r := &request{
		value: value,
		out:   make(chan int),
	}
	s.in <- r
	return <-r.out
}

func (s *Storage) request(value int) int {
	s.count++
	return s.count + value
}
