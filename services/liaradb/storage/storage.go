package storage

import (
	"context"
)

type Storage struct {
	ctx context.Context
	in  chan *request
}

type request struct {
	id  int
	out chan *buffer
}

type buffer struct {
	id int
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
				r.out <- s.request(r.id)
			} else {
				close(r.out)
			}
		case <-ctx.Done():
			close(s.in)
			return
		}
	}
}

func (s *Storage) Request(ctx context.Context, id int) (*buffer, bool) {
	if s.in == nil {
		return nil, false
	}

	r := &request{
		id:  id,
		out: make(chan *buffer),
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

	return nil, false
}

func (s *Storage) request(id int) *buffer {
	return &buffer{
		id: id,
	}
}
