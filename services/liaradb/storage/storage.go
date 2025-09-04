package storage

import (
	"context"
)

type Storage struct {
	ctx context.Context
	in  chan *request
	bm  *BufferManager
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
				b, err := s.loadBuffer(r.blockID)
				r.respond(b, err)
			} else {
				r.close()
			}
		case <-ctx.Done():
			s.close()
			return
		}
	}
}

func (s *Storage) close() {
	close(s.in)
}

func (s *Storage) Request(ctx context.Context, bid BlockID) (*Buffer, error) {
	if s.in == nil {
		return nil, ErrNotInitialized
	}

	r := newRequest(bid)
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

func (s *Storage) loadBuffer(bid BlockID) (*Buffer, error) {
	b := s.bm.Buffer(bid)
	err := b.Load()
	if err != nil {
		return nil, err
	}

	return b, nil
}
