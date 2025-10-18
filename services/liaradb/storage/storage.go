package storage

import (
	"context"
	"io"

	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/raw"
	"github.com/liaradb/liaradb/storage/queue"
)

type Storage struct {
	pinned     map[BlockID]*Buffer
	unpinned   queue.MapQueue[BlockID, *Buffer]
	bufferReqs async.Handler[BlockID, *Buffer]
	appendReqs async.Handler[appendValue, BlockID]
	returns    chan *Buffer
	max        int
	bm         *BufferManager
	highWater  map[string]raw.Offset
}

type bufferRequest = async.Request[BlockID, *Buffer]

type appendRequest = async.Request[appendValue, BlockID]

type appendValue struct {
	fileName string
	reader   io.Reader
}

func NewStorage(fs file.FileSystem, max int, bs int64) *Storage {
	return &Storage{
		bufferReqs: make(chan *bufferRequest),
		appendReqs: make(chan *appendRequest),
		returns:    make(chan *Buffer, max),
		pinned:     make(map[BlockID]*Buffer, max),
		bm:         NewBufferManager(fs, bs),
		max:        max,
		highWater:  make(map[string]raw.Offset),
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
		case r := <-s.appendReqs:
			s.append(r)
		case r := <-s.bufferReqs:
			s.respond(r)
		case b := <-s.returns:
			s.unpin(b)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Storage) append(r *appendRequest) {
	v := r.Value()
	// TODO: Increment
	h := s.highWater[v.fileName]
	bid := BlockID{
		FileName: v.fileName,
		Position: h}

	b, err := s.getBuffer(r.Context(), bid)
	if err != nil {
		r.Reply(bid, err)
		return
	}

	defer b.Release()

	data := make([]byte, b.buffer.Length())
	n, err := v.reader.Read(data)
	if err != nil && err != io.EOF {
		r.Reply(bid, err)
		return
	}

	r.Reply(bid, b.Add(data[:n]))
}

func (s *Storage) respond(r *bufferRequest) {
	// TODO: Create second goroutine
	// One for loaded Buffers, one for non-loaded Buffers
	// This will allow loaded traffic to continue
	b, err := s.getBuffer(r.Context(), r.Value())
	r.Reply(b, err)
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

	// TODO: Don't load here.  Do this in separate goroutine.
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

func (s *Storage) RequestLatest(ctx context.Context, fileName string) (*Buffer, error) {
	return s.Request(ctx, BlockID{
		FileName: fileName,
		Position: -1,
	})
}

// External thread
func (s *Storage) Request(ctx context.Context, bid BlockID) (*Buffer, error) {
	if s.bufferReqs == nil {
		return nil, ErrNotInitialized
	}

	return s.bufferReqs.Send(ctx, bid)
}

// External thread
func (s *Storage) release(b *Buffer) {
	s.returns <- b
}

// TODO: Should this be multiple BlockIDs?
func (s *Storage) Append(ctx context.Context, fileName string, rd io.Reader) (BlockID, error) {
	if s.appendReqs == nil {
		return BlockID{}, ErrNotInitialized
	}

	return s.appendReqs.Send(ctx, appendValue{fileName, rd})
}
