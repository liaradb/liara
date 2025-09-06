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

func (s *Storage) CountPinned() int {
	return len(s.pinned)
}

func (s *Storage) Count() int {
	return len(s.pinned) + len(s.unpinned)
}

func (s *Storage) Run(ctx context.Context, bm *BufferManager, max int) {
	if s.requests == nil {
		s.requests = make(chan *request)
	}
	if s.returns == nil {
		s.returns = make(chan *Buffer)
	}
	if s.pinned == nil {
		s.pinned = make(map[BlockID]*Buffer)
	}
	if s.unpinned == nil {
		s.unpinned = make(map[BlockID]*Buffer)
	}
	s.bm = bm
	s.max = max
	go s.run(ctx)
}

func (s *Storage) run(ctx context.Context) {
	for {
		select {
		case r := <-s.requests:
			s.respond(r)
		case b := <-s.returns:
			s.unpin(b)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Storage) respond(r *request) {
	b, err := s.loadBuffer(r.ctx, r.blockID)
	r.respond(r.ctx, b, err)
}

func (s *Storage) loadBuffer(ctx context.Context, bid BlockID) (*Buffer, error) {
	if b, ok := s.getPinned(bid); ok {
		b.pin()
		return b, nil
	}

	if b, ok := s.getUnpinned(bid); ok {
		b.pin()
		s.moveToPinned(b)
		return b, nil
	}

	b, err := s.allocateOrWait(ctx, bid)
	if err != nil {
		return nil, err
	}

	return b, b.Load(bid)
}

func (s *Storage) unpin(b *Buffer) {
	if b.unpin() {
		s.moveToUnpinned(b)
	}
}

func (s *Storage) moveToPinned(b *Buffer) {
	delete(s.unpinned, b.blockID)
	s.pinned[b.blockID] = b
}

func (s *Storage) moveToUnpinned(b *Buffer) {
	delete(s.pinned, b.blockID)
	s.unpinned[b.blockID] = b
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
	if b, ok := s.popUnpinned(); ok {
		b.pin()
		s.moveToPinned(b)
		return b, nil
	}

	if b, ok := s.allocate(bid); ok {
		b.pin()
		return b, nil
	}

	if b, err := s.waitForRelease(ctx); err != nil {
		return nil, err
	} else {
		b.pin()
		return b, nil
	}
}

func (s *Storage) popUnpinned() (*Buffer, bool) {
	if len(s.unpinned) == 0 {
		return nil, false
	}

	for _, b := range s.unpinned {
		return b, true
	}

	return nil, false
}

func (s *Storage) allocate(bid BlockID) (*Buffer, bool) {
	if s.Count() >= s.max {
		return nil, false
	}

	b := s.bm.Buffer()
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

// External thread
func (s *Storage) Release(ctx context.Context, b *Buffer) {
	select {
	case s.returns <- b:
	case <-ctx.Done():
		return
	}
}
