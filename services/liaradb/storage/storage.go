package storage

import (
	"context"
)

type Storage struct {
	pinned   map[BlockID]*Buffer
	unpinned map[BlockID]*Buffer
	requests chan *request
	returns  chan *Buffer
	max      int
	bm       *BufferManager
}

func (s *Storage) Run(ctx context.Context, bm *BufferManager, max int) {
	if s.requests == nil {
		s.requests = make(chan *request)
	}
	if s.pinned == nil {
		s.pinned = make(map[BlockID]*Buffer)
	}
	s.bm = bm
	s.max = max
	go s.run(ctx)
}

func (s *Storage) run(ctx context.Context) {
	for {
		select {
		case r := <-s.requests:
			s.respond(ctx, r)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Storage) respond(ctx context.Context, r *request) {
	b, err := s.loadBuffer(r.ctx, r.blockID)
	r.respond(ctx, b, err)
}

func (s *Storage) loadBuffer(ctx context.Context, bid BlockID) (*Buffer, error) {
	if b, ok := s.getPinned(bid); ok {
		return b, nil
	}

	if b, ok := s.getUnpinned(bid); ok {
		return b, nil
	}

	b, err := s.allocateOrWait(ctx, bid)
	if err != nil {
		return nil, err
	}

	return b, b.Load()
}

func (s *Storage) getPinned(bid BlockID) (*Buffer, bool) {
	b, ok := s.pinned[bid]
	return b, ok
}

func (s *Storage) getUnpinned(bid BlockID) (*Buffer, bool) {
	b, ok := s.unpinned[bid]
	return b, ok
}

func (s *Storage) allocateOrWait(ctx context.Context, bid BlockID) (*Buffer, error) {
	if b, ok := s.allocate(bid); ok {
		return b, nil
	}

	return s.waitForRelease(ctx)
}

func (s *Storage) allocate(bid BlockID) (*Buffer, bool) {
	if len(s.pinned) == s.max {
		return nil, false
	}

	b := s.bm.Buffer(bid)
	s.pinned[bid] = b
	return b, true
}

func (s *Storage) waitForRelease(ctx context.Context) (*Buffer, error) {
	select {
	case b := <-s.returns:
		return b, nil
	case <-ctx.Done():
		return nil, context.Canceled
	}
}

// External thread
func (s *Storage) Request(ctx context.Context, bid BlockID) (*Buffer, error) {
	if s.requests == nil {
		return nil, ErrNotInitialized
	}

	r := newRequest(ctx, bid)
	select {
	case s.requests <- r:
	case <-ctx.Done():
		return nil, context.Canceled
	}

	return r.wait(ctx)
}
