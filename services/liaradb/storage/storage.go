package storage

import (
	"context"
	"errors"
	"io"
	"path"

	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/storage/queue"
)

type Storage struct {
	bufferSize int64 // TODO: Do we need this?
	fs         file.FileSystem
	dir        string
	pinned     map[BlockID]*Buffer
	unpinned   queue.MapQueue[BlockID, *Buffer]
	bufferReqs async.Handler[bufferQuery, *Buffer]
	highWReqs  async.Handler[string, BlockID]
	returns    chan *Buffer
	max        int
	highWater  map[string]Offset
}

func NewStorage(fs file.FileSystem, max int, bs int64, dir string) *Storage {
	return &Storage{
		bufferSize: bs,
		fs:         fs,
		dir:        dir,
		bufferReqs: make(chan *bufferRequest),
		highWReqs:  make(async.Handler[string, BlockID]),
		returns:    make(chan *Buffer, max),
		pinned:     make(map[BlockID]*Buffer, max),
		max:        max,
		highWater:  make(map[string]Offset),
	}
}

func (s *Storage) BufferSize() int64 { return s.bufferSize }

func (s *Storage) CountPinned() int {
	return len(s.pinned)
}

func (s *Storage) Count() int {
	return len(s.pinned) + s.unpinned.Count()
}

func (s *Storage) Run(ctx context.Context) error {
	// TODO: Test this
	if err := s.fs.MkDirAll(s.dir); err != nil {
		return err
	}

	go s.run(ctx)

	return nil
}

func (s *Storage) run(ctx context.Context) {
	for {
		select {
		case r := <-s.bufferReqs:
			s.requestBuffer(r)
		case r := <-s.highWReqs:
			s.getHighWater(r)
		case b := <-s.returns:
			s.returnBuffer(b)
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

func (s *Storage) requestBuffer(r *bufferRequest) {
	// TODO: Create second goroutine
	// One for loaded Buffers, one for non-loaded Buffers
	// This will allow loaded traffic to continue
	if bid, err := s.getBufferID(r.Value()); err != nil {
		r.Reply(nil, err)
	} else {
		b, err := s.getBuffer(r.Context(), bid)
		r.Reply(b, err)
	}
}

func (s *Storage) getBufferID(v bufferQuery) (BlockID, error) {
	switch v.queryType {
	case bufferQueryTypeByID:
		return v.bid, nil
	case bufferQueryTypeCurrent:
		return s.highBlockID(v.fileName), nil
	case bufferQueryTypeNext:
		s.incrementHighWater(v.fileName)
		return s.highBlockID(v.fileName), nil
	default:
		return BlockID{}, errors.New("invalid request")
	}
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

	return s.getUnpinned(bid)
}

func (s *Storage) getUnpinned(bid BlockID) (*Buffer, bool) {
	b, ok := s.unpinned.Remove(bid)
	if ok {
		b.pin()
		s.pinned[b.blockID] = b
	}
	return b, ok
}

func (s *Storage) getUnloaded(ctx context.Context, bid BlockID) (*Buffer, error) {
	b, err := s.popAllocateOrWait(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: Don't load here.  Do this in separate goroutine.
	// TODO: Don't load if just allocated
	if err := b.load(bid); err != nil {
		return nil, err
	}

	b.pin()
	s.pinned[b.blockID] = b

	return b, nil
}

func (s *Storage) getPinned(bid BlockID) (*Buffer, bool) {
	b, ok := s.pinned[bid]
	return b, ok
}

func (s *Storage) popAllocateOrWait(ctx context.Context) (*Buffer, error) {
	if b, ok := s.popUnpinned(); ok {
		return b, nil
	}

	if b, ok := s.allocate(); ok {
		return b, nil
	}

	return s.waitForRelease(ctx)
}

func (s *Storage) popUnpinned() (*Buffer, bool) {
	return s.unpinned.Pop()
}

func (s *Storage) allocate() (*Buffer, bool) {
	if s.Count() >= s.max {
		return nil, false
	}

	return newBuffer(s), true
}

func (s *Storage) waitForRelease(ctx context.Context) (*Buffer, error) {
	for {
		b, err := s.getReturn(ctx)
		if err != nil {
			return nil, err
		}

		if s.unpinAfterRelease(b) {
			return b, nil
		}
	}
}

func (s *Storage) getReturn(ctx context.Context) (*Buffer, error) {
	select {
	case b := <-s.returns:
		return b, nil
	case <-ctx.Done():
		return nil, context.Canceled
	}
}

func (s *Storage) unpinAfterRelease(b *Buffer) bool {
	if b.unpin() {
		delete(s.pinned, b.blockID)
		return true
	}
	return false
}

func (s *Storage) getHighWater(r *async.Request[string, BlockID]) {
	r.Reply(s.highBlockID(r.Value()), nil)
}

// Doesn't change BlockID
func (s *Storage) returnBuffer(b *Buffer) {
	if b.unpin() {
		s.moveToUnpinned(b)
	}
}

func (s *Storage) moveToUnpinned(b *Buffer) {
	delete(s.pinned, b.blockID)
	s.unpinned.Push(b.blockID, b)
}

func (s *Storage) Highwater(ctx context.Context, fileName string) (BlockID, error) {
	if s.highWReqs == nil {
		return BlockID{}, ErrNotInitialized
	}

	return s.highWReqs.Send(ctx, fileName)
}

// TODO: Is this still needed?
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

	return s.bufferReqs.Send(ctx, newBufferByIDQuery(bid))
}

// External thread
// TODO: Test this
func (s *Storage) RequestCurrent(ctx context.Context, fileName string) (*Buffer, error) {
	if s.bufferReqs == nil {
		return nil, ErrNotInitialized
	}

	return s.bufferReqs.Send(ctx, newCurrentBufferQuery(fileName))
}

// External thread
// TODO: Test this
func (s *Storage) RequestNext(ctx context.Context, fileName string) (*Buffer, error) {
	if s.bufferReqs == nil {
		return nil, ErrNotInitialized
	}

	return s.bufferReqs.Send(ctx, newNextBufferQuery(fileName))
}

// External thread
func (s *Storage) release(b *Buffer) {
	s.returns <- b
}

func (s *Storage) load(b *Buffer) error {
	f, err := s.openFile(b)
	if err != nil {
		return err
	}

	// TODO: Do we need to check io.EOF?
	if err := b.read(f); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (s *Storage) flush(b *Buffer) error {
	f, err := s.openFile(b)
	if err != nil {
		return err
	}

	return b.write(f)
}

func (s *Storage) openFile(b *Buffer) (file.File, error) {
	// TODO: Test this
	return s.fs.OpenFile(path.Join(s.dir, b.blockID.FileName))
}

// TODO: Test this
func (s *Storage) FlushAll() error {
	for _, b := range s.pinned {
		if err := b.FlushIfDirty(); err != nil {
			return err
		}
	}

	for b := range s.unpinned.Iterate() {
		if err := b.FlushIfDirty(); err != nil {
			return err
		}
	}

	return nil
}
