package storage

import (
	"context"
)

type Storage struct {
	in chan *request
	bm *BufferManager
}

func (s *Storage) Run(ctx context.Context, bm *BufferManager) {
	if s.in == nil {
		s.in = make(chan *request)
	}
	s.bm = bm
	go s.run(ctx)
}

func (s *Storage) run(ctx context.Context) {
	for {
		select {
		case r := <-s.in:
			s.respond(r)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Storage) respond(r *request) {
	b, err := s.loadBuffer(r.blockID)
	r.respond(b, err)
}

func (s *Storage) loadBuffer(bid BlockID) (*Buffer, error) {
	b := s.bm.Buffer(bid)
	err := b.Load()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// External thread
func (s *Storage) Request(ctx context.Context, bid BlockID) (*Buffer, error) {
	if s.in == nil {
		return nil, ErrNotInitialized
	}

	r := newRequest(bid)
	select {
	case s.in <- r:
	case <-ctx.Done():
		return nil, context.Canceled
	}

	return r.wait(ctx)
}
