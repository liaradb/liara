package storage

import (
	"context"

	"github.com/cardboardrobots/liaradb/storage/queue"
)

type Storage struct {
	pinned   map[BlockID]*Buffer
	unpinned queue.MapQueue[BlockID, *Buffer]
	requests chan *request
	returns  chan *Buffer
	max      int
	bm       *BufferManager
}

func NewStorage(fs FS, max int, bs int64) *Storage {
	return &Storage{
		requests: make(chan *request),
		returns:  make(chan *Buffer, max),
		pinned:   make(map[BlockID]*Buffer, max),
		bm:       NewBufferManager(fs, bs),
		max:      max,
	}
}

func (s *Storage) CountPinned() int {
	return len(s.pinned)
}

func (s *Storage) Count() int {
	return len(s.pinned) + s.unpinned.Count()
}

func (s *Storage) Run(ctx context.Context) {
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
	// TODO: Create second goroutine
	// One for loaded Buffers, one for non-loaded Buffers
	// This will allow loaded traffic to continue
	b, err := s.getBuffer(r.ctx, r.blockID)
	r.respond(r.ctx, b, err)
}

func (s *Storage) getBuffer(ctx context.Context, bid BlockID) (*Buffer, error) {
	if b, ok := s.getLoaded(bid); ok {
		return b, nil
	}

	return s.getUnloaded(ctx, bid)
}

func (s *Storage) getLoaded(bid BlockID) (*Buffer, bool) {
	if b, ok := s.getPinned(bid); ok {
		b.pin()
		return b, true
	}

	if b, ok := s.unpinned.Remove(bid); ok {
		b.pin()
		s.moveToPinned(b)
		return b, true
	}

	return nil, false
}

func (s *Storage) getUnloaded(ctx context.Context, bid BlockID) (*Buffer, error) {
	b, err := s.popAllocateOrWait(ctx, bid)
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
	s.unpinned.Remove(b.blockID)
	s.pinned[b.blockID] = b
}

func (s *Storage) moveToUnpinned(b *Buffer) {
	delete(s.pinned, b.blockID)
	s.unpinned.Push(b.blockID, b)
}

func (s *Storage) getPinned(bid BlockID) (*Buffer, bool) {
	b, ok := s.pinned[bid]
	return b, ok
}

func (s *Storage) popAllocateOrWait(ctx context.Context, bid BlockID) (*Buffer, error) {
	if b, ok := s.popUnpinned(); ok {
		return b, nil
	}

	if b, ok := s.allocate(bid); ok {
		return b, nil
	}

	return s.waitForRelease(ctx)
}

func (s *Storage) popUnpinned() (*Buffer, bool) {
	b, ok := s.unpinned.Pop()
	if !ok {
		return nil, false
	}

	b.pin()
	s.moveToPinned(b)
	return b, true
}

func (s *Storage) allocate(bid BlockID) (*Buffer, bool) {
	if s.Count() >= s.max {
		return nil, false
	}

	b := NewBuffer(s)
	s.pinned[bid] = b
	b.pin()
	return b, true
}

func (s *Storage) waitForRelease(ctx context.Context) (*Buffer, error) {
	select {
	case b := <-s.returns:
		b.pin()
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
func (s *Storage) release(b *Buffer) {
	s.returns <- b
}
