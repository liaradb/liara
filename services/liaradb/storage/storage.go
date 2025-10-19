package storage

import (
	"context"
	"errors"
	"io"
	"iter"

	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/raw"
	"github.com/liaradb/liaradb/storage/queue"
	"github.com/liaradb/liaradb/storage/record"
)

type Storage struct {
	pinned     map[BlockID]*Buffer
	unpinned   queue.MapQueue[BlockID, *Buffer]
	bufferReqs async.Handler[butterRequestQuery, *Buffer]
	highWReqs  async.Handler[string, BlockID]
	returns    chan *Buffer
	max        int
	bm         *BufferManager
	highWater  map[string]raw.Offset
}

type bufferRequest = async.Request[butterRequestQuery, *Buffer]

type butterRequestQuery struct {
	bid       BlockID
	fileName  string
	queryType bufferRequestQueryType
}

type bufferRequestQueryType int

const (
	bufferRequestQueryTypeByID = iota
	bufferRequestQueryTypeCurrent
	bufferRequestQueryTypeNext
)

func NewStorage(fs file.FileSystem, max int, bs int64) *Storage {
	return &Storage{
		bufferReqs: make(chan *bufferRequest),
		highWReqs:  make(async.Handler[string, BlockID]),
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
		case r := <-s.bufferReqs:
			s.respond(r)
		case r := <-s.highWReqs:
			s.getHighWater(r)
		case b := <-s.returns:
			s.unpin(b)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Storage) incrementHighWater(fileName string) {
	s.highWater[fileName]++
}

func (s *Storage) highBlockID(fileName string) BlockID {
	return BlockID{
		FileName: fileName,
		Position: s.highWater[fileName],
	}
}

func (s *Storage) respond(r *bufferRequest) {
	// TODO: Create second goroutine
	// One for loaded Buffers, one for non-loaded Buffers
	// This will allow loaded traffic to continue
	v := r.Value()

	var bid BlockID
	switch v.queryType {
	case bufferRequestQueryTypeByID:
		bid = v.bid
	case bufferRequestQueryTypeCurrent:
		bid = s.highBlockID(v.fileName)
	case bufferRequestQueryTypeNext:
		s.incrementHighWater(v.fileName)
		bid = s.highBlockID(v.fileName)
	default:
		r.Reply(nil, errors.New("invalid request"))
		return
	}

	b, err := s.getBuffer(r.Context(), bid)
	r.Reply(b, err)
}

func (s *Storage) getHighWater(r *async.Request[string, BlockID]) {
	r.Reply(s.highBlockID(r.Value()), nil)
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

func (s *Storage) Highwater(ctx context.Context, fileName string) (BlockID, error) {
	if s.highWReqs == nil {
		return BlockID{}, ErrNotInitialized
	}

	return s.highWReqs.Send(ctx, fileName)
}

// External thread
func (s *Storage) Request(ctx context.Context, bid BlockID) (*Buffer, error) {
	if s.bufferReqs == nil {
		return nil, ErrNotInitialized
	}

	return s.bufferReqs.Send(ctx, butterRequestQuery{
		bid: bid,
	})
}

// External thread
// TODO: Test this
func (s *Storage) RequestCurrent(ctx context.Context, fileName string) (*Buffer, error) {
	if s.bufferReqs == nil {
		return nil, ErrNotInitialized
	}

	return s.bufferReqs.Send(ctx, butterRequestQuery{
		fileName:  fileName,
		queryType: bufferRequestQueryTypeCurrent,
	})
}

// External thread
// TODO: Test this
func (s *Storage) RequestNext(ctx context.Context, fileName string) (*Buffer, error) {
	if s.bufferReqs == nil {
		return nil, ErrNotInitialized
	}

	return s.bufferReqs.Send(ctx, butterRequestQuery{
		fileName:  fileName,
		queryType: bufferRequestQueryTypeNext,
	})
}

// External thread
func (s *Storage) release(b *Buffer) {
	s.returns <- b
}

// TODO: Should this be multiple BlockIDs?
func (s *Storage) Append(ctx context.Context, fileName string, rd io.Reader) (BlockID, error) {
	// TODO: Find a better way to get this
	data := make([]byte, s.bm.bufferSize)
	n, err := rd.Read(data)
	if err != nil && err != io.EOF {
		return BlockID{}, err
	}

	bid, err := s.appendCurrent(ctx, fileName, data[:n])
	if err == record.ErrInsufficientSpace {
		bid, err = s.appendNext(ctx, fileName, data[:n])
	}

	return bid, err
}

func (s *Storage) appendCurrent(ctx context.Context, fileName string, data []byte) (BlockID, error) {
	b, err := s.RequestCurrent(ctx, fileName)
	if err != nil {
		return BlockID{}, err
	}

	defer b.Release()

	if err := b.Add(data); err != nil {
		return BlockID{}, err
	}

	return b.BlockID(), nil
}

func (s *Storage) appendNext(ctx context.Context, fileName string, data []byte) (BlockID, error) {
	b, err := s.RequestNext(ctx, fileName)
	if err != nil {
		return BlockID{}, err
	}

	defer b.Release()

	if err := b.Add(data); err != nil {
		return BlockID{}, err
	}

	return b.BlockID(), nil
}

func (s *Storage) Iterate(ctx context.Context, fn string) iter.Seq2[*Buffer, error] {
	return func(yield func(*Buffer, error) bool) {
		highBid, err := s.Highwater(ctx, fn)
		if err != nil {
			yield(nil, err)
			return
		}

		bid := BlockID{FileName: fn}
		for bid.Position < highBid.Position {
			b, err := s.Request(ctx, bid)
			if err != nil {
				yield(nil, err)
				return
			}

			defer b.Release()
			if !yield(b, err) {
				return
			}

			bid.Position++
		}
	}
}
