package storage

import (
	"context"
)

type Storage struct {
	ctx   context.Context
	in    chan *request
	count int
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
	s.ctx = ctx
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

func (s *Storage) Request(ctx context.Context, value int) int {
	r := &request{
		value: value,
		out:   make(chan int),
	}
	select {
	case s.in <- r:
	case <-s.ctx.Done():
	case <-ctx.Done():
	}

	select {
	case o := <-r.out:
		return o
	case <-s.ctx.Done():
	case <-ctx.Done():
	}

	return 0
}

func (s *Storage) request(value int) int {
	s.count++
	return s.count + value
}
