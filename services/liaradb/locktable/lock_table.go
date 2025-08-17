package locktable

import (
	"context"
)

type LockTable[K comparable] struct {
	requestLists map[K]*lockRequestList[K]
	in           chan K
	requests     chan *lockRequest[K]
	cancel       context.CancelFunc
}

func NewLockTable[K comparable](ctx context.Context, inSize int) *LockTable[K] {
	ctx, cancel := context.WithCancel(ctx)
	qp := &LockTable[K]{
		requestLists: make(map[K]*lockRequestList[K]),
		in:           make(chan K, inSize),
		requests:     make(chan *lockRequest[K]),
		cancel:       cancel,
	}

	go qp.run(ctx)

	return qp
}

func (lt *LockTable[K]) run(ctx context.Context) {
	for {
		select {
		case r := <-lt.requests:
			lt.request(r)
		case k := <-lt.in:
			lt.unlock(k)
		case <-ctx.Done():
			return
		}
	}
}

func (lt *LockTable[K]) Close() {
	lt.cancel()
}

func (lt *LockTable[K]) request(lr *lockRequest[K]) {
	lrl, ok := lt.requestLists[lr.key]
	if !ok {
		lrl = newLockRequestList[K]()
		lt.requestLists[lr.key] = lrl
	}

	lrl.lock(lr)
}

func (lt *LockTable[K]) unlock(k K) {
	lrl, ok := lt.requestLists[k]
	if !ok {
		return
	}

	if lrl.unlock() {
		delete(lt.requestLists, k)
	}
}

func (lt *LockTable[K]) sLock(ctx context.Context, k K) bool {
	lr := newSharedLockRequest(ctx, k)
	lt.requests <- lr
	return lr.wait(ctx)
}

func (lt *LockTable[K]) xLock(ctx context.Context, k K) bool {
	lr := newExclusiveLockRequest(ctx, k)
	lt.requests <- lr
	return lr.wait(ctx)
}

func (lt *LockTable[K]) upgradeLock(ctx context.Context, k K) bool {
	lr := newUpgradeLockRequest(ctx, k)
	lt.requests <- lr
	return lr.wait(ctx)
}

func (lt *LockTable[K]) release(k K) {
	lt.in <- k
}
