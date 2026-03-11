package storage

import (
	"context"
	"errors"
	"path"

	"github.com/liaradb/liaradb/async"
	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/storage/queue"
)

type Storage struct {
	bufferSize int64
	fs         file.FileSystem
	dir        string
	pinned     map[link.BlockID]*Buffer
	unpinned   queue.MapQueue[link.BlockID, *Buffer]
	bufferReqs async.Handler[bufferQuery, *Buffer]
	highWReqs  async.Handler[link.FileName, link.BlockID]
	returns    chan *Buffer
	max        int
	highWater  map[link.FileName]link.FilePosition
}

func New(fs file.FileSystem, max int, bs int64, dir string) *Storage {
	return &Storage{
		bufferSize: bs,
		fs:         fs,
		dir:        dir,
		bufferReqs: make(chan *bufferRequest),
		highWReqs:  make(async.Handler[link.FileName, link.BlockID]),
		returns:    make(chan *Buffer, max),
		pinned:     make(map[link.BlockID]*Buffer, max),
		max:        max,
		highWater:  make(map[link.FileName]link.FilePosition),
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

func (s *Storage) incrementHighWater(fn link.FileName) {
	s.highWater[fn]++
}

func (s *Storage) highBlockID(fn link.FileName) link.BlockID {
	return fn.BlockID(s.highWater[fn])
}

func (s *Storage) requestBuffer(r *bufferRequest) {
	if bid, err := s.getBufferID(r.Value()); err != nil {
		r.Reply(nil, err)
	} else {
		b, err := s.getBuffer(r.Context(), bid, r.Value().isNext())
		r.Reply(b, err)
	}
}

func (s *Storage) getBufferID(v bufferQuery) (link.BlockID, error) {
	switch v.queryType {
	case bufferQueryTypeByID:
		return v.bid, nil
	case bufferQueryTypeCurrent:
		return s.highBlockID(v.fileName), nil
	case bufferQueryTypeNext:
		s.incrementHighWater(v.fileName)
		return s.highBlockID(v.fileName), nil
	default:
		return link.BlockID{}, ErrInvalidRequest
	}
}

func (s *Storage) getBuffer(ctx context.Context, bid link.BlockID, next bool) (*Buffer, error) {
	if b, ok := s.getLoaded(bid); ok {
		return b, nil
	}

	return s.getUnloaded(ctx, bid, next)
}

func (s *Storage) getLoaded(bid link.BlockID) (*Buffer, bool) {
	if b, ok := s.getPinned(bid); ok {
		b.pin()
		return b, true
	}

	return s.getUnpinned(bid)
}

func (s *Storage) getUnpinned(bid link.BlockID) (*Buffer, bool) {
	b, ok := s.unpinned.Remove(bid)
	if ok {
		b.pin()
		s.pinned[b.blockID] = b
	}
	return b, ok
}

func (s *Storage) getUnloaded(ctx context.Context, bid link.BlockID, next bool) (*Buffer, error) {
	b, err := s.popAllocateOrWait(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: Don't load here.  Do this in separate goroutine.
	// Create second goroutine
	// One for loaded Buffers, one for non-loaded Buffers
	// This will allow loaded traffic to continue
	if err := b.load(bid, next); err != nil {
		return nil, err
	}

	b.pin()
	s.pinned[b.blockID] = b

	return b, nil
}

func (s *Storage) getPinned(bid link.BlockID) (*Buffer, bool) {
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

func (s *Storage) getHighWater(r *async.Request[link.FileName, link.BlockID]) {
	fn := r.Value()
	if _, err := s.openHighwater(fn); err != nil {
		r.Reply(link.BlockID{}, err)
		return
	}

	r.Reply(s.highBlockID(fn), nil)
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

func (s *Storage) Highwater(ctx context.Context, fn link.FileName) (link.BlockID, error) {
	if s.highWReqs == nil {
		return link.BlockID{}, ErrNotInitialized
	}

	return s.highWReqs.Send(ctx, fn)
}

// External thread
func (s *Storage) Request(ctx context.Context, bid link.BlockID) (*Buffer, error) {
	if s.bufferReqs == nil {
		return nil, ErrNotInitialized
	}

	return s.bufferReqs.Send(ctx, newBufferByIDQuery(bid))
}

// External thread
func (s *Storage) RequestCurrent(ctx context.Context, fn link.FileName) (*Buffer, error) {
	if s.bufferReqs == nil {
		return nil, ErrNotInitialized
	}

	return s.bufferReqs.Send(ctx, newCurrentBufferQuery(fn))
}

// External thread
func (s *Storage) RequestNext(ctx context.Context, fn link.FileName) (*Buffer, error) {
	if s.bufferReqs == nil {
		return nil, ErrNotInitialized
	}

	return s.bufferReqs.Send(ctx, newNextBufferQuery(fn))
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

	return b.read(f)
}

func (s *Storage) flush(b *Buffer) error {
	f, err := s.openFile(b)
	if err != nil {
		return err
	}

	return b.write(f)
}

func (s *Storage) openFile(b *Buffer) (file.File, error) {
	return s.openHighwater(b.blockID.FileName())
}

func (s *Storage) openHighwater(fn link.FileName) (file.File, error) {
	// TODO: Test this
	f, err := s.fs.OpenFile(path.Join(s.dir, fn.String()))
	if err != nil {
		return nil, err
	}

	if err := s.initHighwater(fn, f); err != nil {
		return nil, errors.Join(err, f.Close())
	}

	return f, nil
}

func (s *Storage) initHighwater(fn link.FileName, f file.File) error {
	if _, ok := s.highWater[fn]; ok {
		return nil
	}

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	size := stat.Size()
	s.highWater[fn] = link.FilePosition(size / s.bufferSize)

	return nil
}

func (s *Storage) FlushAll() error {
	for _, b := range s.pinned {
		if err := b.flushIfDirty(); err != nil {
			return err
		}
	}

	for b := range s.unpinned.Iterate() {
		if err := b.flushIfDirty(); err != nil {
			return err
		}
	}

	return nil
}
