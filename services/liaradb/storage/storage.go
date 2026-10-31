package storage

import (
	"context"
	"errors"
	"path"
	"sync"

	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/async"
)

// TODO: Create dirty list in FreePool. This will make flushes simpler.
// TODO: Verify pool is filled up to max
type Storage struct {
	bufferSize int64
	fs         filecache.FileSystem
	dir        string
	pinned     map[link.FileBlockID]*Buffer
	unpinned   FreePool
	bufferReqs async.Handler[bufferQuery, *Buffer]
	highWReqs  async.Handler[link.FileName, link.FileBlockID]
	flushReqs  async.CommandHandler[struct{}]
	returns    chan *Buffer
	max        int
	highWater  map[link.FileName]link.BlockID
	hwMux      sync.RWMutex
}

func New(fs filecache.FileSystem, unpinned FreePool, max int, bs int64, dir string) *Storage {
	return &Storage{
		bufferSize: bs,
		fs:         fs,
		dir:        dir,
		bufferReqs: make(chan *bufferRequest),
		highWReqs:  make(async.Handler[link.FileName, link.FileBlockID]),
		flushReqs:  make(async.CommandHandler[struct{}]),
		returns:    make(chan *Buffer, max),
		pinned:     make(map[link.FileBlockID]*Buffer, max),
		unpinned:   unpinned,
		max:        max,
		highWater:  make(map[link.FileName]link.BlockID),
	}
}

func (s *Storage) BufferSize() int64 { return s.bufferSize }
func (s *Storage) Dir() string       { return s.dir }

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
		case r := <-s.flushReqs:
			s.flush(r)
		case b := <-s.returns:
			s.returnBuffer(b)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Storage) requestBuffer(r *bufferRequest) {
	if fbid, err := s.getBufferID(r.Value()); err != nil {
		r.Reply(nil, err)
	} else {
		r.Reply(s.getBuffer(r.Context(), fbid, r.Value().isNext()))
	}
}

func (s *Storage) getBufferID(v bufferQuery) (link.FileBlockID, error) {
	switch v.queryType {
	case bufferQueryTypeByID:
		return v.fbid, nil
	case bufferQueryTypeCurrent:
		return s.highFileBlockID(v.fileName), nil
	case bufferQueryTypeNext:
		s.incrementHighWater(v.fileName)
		return s.highFileBlockID(v.fileName), nil
	default:
		return link.FileBlockID{}, ErrInvalidRequest
	}
}

func (s *Storage) getBuffer(ctx context.Context, fbid link.FileBlockID, next bool) (*Buffer, error) {
	if b, ok := s.getLoaded(fbid); ok {
		return b, nil
	}

	return s.getUnloaded(ctx, fbid, next)
}

func (s *Storage) getLoaded(fbid link.FileBlockID) (*Buffer, bool) {
	if b, ok := s.getPinned(fbid); ok {
		b.pin()
		return b, true
	}

	return s.getUnpinned(fbid)
}

func (s *Storage) getUnpinned(fbid link.FileBlockID) (*Buffer, bool) {
	b, ok := s.unpinned.Remove(fbid)
	if ok {
		b.pin()
		s.pinned[b.FileBlockID()] = b
	}
	return b, ok
}

func (s *Storage) getUnloaded(ctx context.Context, fbid link.FileBlockID, next bool) (*Buffer, error) {
	b, err := s.popAllocateOrWait(ctx)
	if err != nil {
		return nil, err
	}

	b.load(fbid, next)
	b.pin()
	s.pinned[b.FileBlockID()] = b

	return b, nil
}

func (s *Storage) getPinned(fbid link.FileBlockID) (*Buffer, bool) {
	b, ok := s.pinned[fbid]
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
		delete(s.pinned, b.FileBlockID())
		return true
	}
	return false
}

func (s *Storage) getHighWater(r *async.Request[link.FileName, link.FileBlockID]) {
	fn := r.Value()
	if _, err := s.openHighwater(fn); err != nil {
		r.Reply(link.FileBlockID{}, err)
		return
	}

	r.Reply(s.highFileBlockID(fn), nil)
}

func (s *Storage) FlushUnpinned(ctx context.Context) error {
	return s.flushReqs.Send(ctx, struct{}{})
}

// TODO: Don't block the Storage thread while flushing
func (s *Storage) flush(r *async.Command[struct{}]) {
	r.Reply(s.flushUnpinned())
}

func (s *Storage) flushUnpinned() error {
	s.drainReturns()

	for b := range s.unpinned.Iterate() {
		if err := b.flushIfDirty(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) drainReturns() {
	for {
		select {
		case b := <-s.returns:
			s.returnBuffer(b)
		default:
			return
		}
	}
}

// Doesn't change FileBlockID
func (s *Storage) returnBuffer(b *Buffer) {
	if b.unpin() {
		s.moveToUnpinned(b)
	}
}

func (s *Storage) moveToUnpinned(b *Buffer) {
	delete(s.pinned, b.FileBlockID())
	s.unpinned.Push(b.FileBlockID(), b)
}

func (s *Storage) Highwater(ctx context.Context, fn link.FileName) (link.FileBlockID, error) {
	if s.highWReqs == nil {
		return link.FileBlockID{}, ErrNotInitialized
	}

	return s.highWReqs.Send(ctx, fn)
}

// External thread
func (s *Storage) Request(ctx context.Context, fbid link.FileBlockID) (*Buffer, error) {
	b, err := s.request(ctx, newBufferByIDQuery(fbid))
	if err != nil {
		return nil, err
	}

	return b, b.loadOnce()
}

// External thread
func (s *Storage) RequestCurrent(ctx context.Context, fn link.FileName) (*Buffer, error) {
	b, err := s.request(ctx, newCurrentBufferQuery(fn))
	if err != nil {
		return nil, err
	}

	return b, b.loadOnce()
}

// External thread
func (s *Storage) RequestNext(ctx context.Context, fn link.FileName) (*Buffer, error) {
	b, err := s.request(ctx, newNextBufferQuery(fn))
	if err != nil {
		return nil, err
	}

	return b, b.loadOnce()
}

func (s *Storage) request(ctx context.Context, q bufferQuery) (*Buffer, error) {
	if s.bufferReqs == nil {
		return nil, ErrNotInitialized
	}

	return s.bufferReqs.SendOrCancel(ctx, q, func(b *Buffer, err error) {
		if err == nil {
			// This will leak if we have a buffer, but the request is cancelled
			b.Release()
		}
	})
}

// External thread
// This is non-blocking as channel is the size of the total buffer count
func (s *Storage) release(b *Buffer) {
	s.returns <- b
}

func (s *Storage) openFile(fbid link.FileBlockID) (filecache.File, error) {
	return s.openHighwater(fbid.FileName())
}

func (s *Storage) openHighwater(fn link.FileName) (filecache.File, error) {
	f, err := s.fs.OpenFile(path.Join(s.dir, fn.String()))
	if err != nil {
		return nil, err
	}

	if err := s.initHighwater(fn, f); err != nil {
		return nil, errors.Join(err, f.Close())
	}

	return f, nil
}

func (s *Storage) initHighwater(fn link.FileName, f filecache.File) error {
	s.hwMux.Lock()
	defer s.hwMux.Unlock()

	if _, ok := s.highWater[fn]; ok {
		return nil
	}

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	size := stat.Size()
	s.highWater[fn] = link.BlockID(size / s.bufferSize)

	return nil
}

func (s *Storage) incrementHighWater(fn link.FileName) {
	s.hwMux.Lock()
	defer s.hwMux.Unlock()

	s.highWater[fn]++
}

func (s *Storage) highFileBlockID(fn link.FileName) link.FileBlockID {
	s.hwMux.RLock()
	defer s.hwMux.RUnlock()

	return fn.FileBlockID(s.highWater[fn])
}
