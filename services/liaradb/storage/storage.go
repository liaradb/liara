package storage

import (
	"context"
)

type Storage struct {
	ctx context.Context
	in  chan *request
	bm  *BufferManager
}

type request struct {
	blockID BlockID
	out     chan *response
}

type response struct {
	buffer *Buffer
	err    error
}

func (s *Storage) Run(ctx context.Context, bm *BufferManager) {
	s.in = make(chan *request)
	s.ctx = ctx
	s.bm = bm
	go s.run(ctx)
}

func (s *Storage) run(ctx context.Context) {
	for {
		select {
		case r, ok := <-s.in:
			if ok {
				r.out <- s.request(r.blockID)
			} else {
				close(r.out)
			}
		case <-ctx.Done():
			close(s.in)
			return
		}
	}
}

func (s *Storage) Request(ctx context.Context, bid BlockID) (*Buffer, error) {
	if s.in == nil {
		return nil, ErrNotInitialized
	}

	r := &request{
		blockID: bid,
		out:     make(chan *response),
	}
	select {
	case s.in <- r:
	case <-s.ctx.Done():
	case <-ctx.Done():
	}

	select {
	case o, ok := <-r.out:
		if ok {
			return o.buffer, o.err
		} else {
			return nil, ErrRequestClosed
		}
	case <-s.ctx.Done():
	case <-ctx.Done():
	}

	return nil, context.Canceled
}

func (s *Storage) request(bid BlockID) *response {
	b := s.bm.Buffer(bid)
	err := b.Load()
	return &response{
		buffer: b,
		err:    err,
	}
}
