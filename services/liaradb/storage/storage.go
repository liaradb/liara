package storage

import (
	"context"
)

type Storage struct {
	ctx   context.Context
	in    chan *request
	count int
}

type request struct {
	value int
	out   chan int
}

func (s *Storage) Run(ctx context.Context) {
	s.in = make(chan *request)
	s.ctx = ctx
	go s.run(ctx)
}

func (s *Storage) run(ctx context.Context) {
	for {
		select {
		case r, ok := <-s.in:
			if ok {
				r.out <- s.request(r.value)
			} else {
				close(r.out)
			}
		case <-ctx.Done():
			close(s.in)
			return
		}
	}
}

func (s *Storage) Request(ctx context.Context, value int) (int, bool) {
	if s.in == nil {
		return 0, false
	}

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
	case o, ok := <-r.out:
		if ok {
			return o, true
		}
	case <-s.ctx.Done():
	case <-ctx.Done():
	}

	return 0, false
}

func (s *Storage) request(value int) int {
	s.count++
	return s.count + value
}
