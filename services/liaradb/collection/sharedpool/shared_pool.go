package sharedpool

import (
	"context"
	"iter"
)

type SharedPool[K comparable, V item[K]] struct {
	unpinned map[int]V
	claimed  map[K]V
	requests chan *request[K, V]
	in       chan V
	ctx      context.Context
	cancel   context.CancelFunc
}

type item[K comparable] interface {
	comparable
	Id() int
	Block() (K, bool)
	ReplaceBlock(K) error
	Pin()
}

func NewSharedPool[K comparable, V item[K]](size int) *SharedPool[K, V] {
	return &SharedPool[K, V]{
		unpinned: make(map[int]V, size),
		claimed:  make(map[K]V, size),
		requests: make(chan *request[K, V]),
		in:       make(chan V),
	}
}

func (sp *SharedPool[K, V]) Close() {
	if sp.cancel == nil {
		return
	}

	sp.cancel()
}

func (sp *SharedPool[K, V]) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	sp.ctx = ctx
	sp.cancel = cancel

	go sp.run()
}

func (sp *SharedPool[K, V]) run() {
	for {
		select {
		case v := <-sp.in:
			sp.addUnpinned(v)
		case r := <-sp.requests:
			sp.send(r)
		case <-sp.ctx.Done():
			return
		}
	}
}

func (sp *SharedPool[K, V]) addUnpinned(v V) {
	if k, ok := v.Block(); ok {
		sp.claimed[k] = v
	}
	sp.unpinned[v.Id()] = v
}

func (sp *SharedPool[K, V]) send(r *request[K, V]) {
	if v, ok := sp.get(r.ctx, r.key); ok {
		r.send(v)
	} else {
		r.cancel(v)
	}
}

func (sp *SharedPool[K, V]) get(ctx context.Context, k K) (V, bool) {
	// Find claimed
	if v, ok := sp.getClaimed(k); ok {
		return v, true
	}

	// Find any
	if v, ok := sp.getUnclaimed(k); ok {
		return v, true
	}

	// Wait for Release
	return sp.waitForRelease(ctx, k)
}

func (sp *SharedPool[K, V]) getClaimed(k K) (V, bool) {
	v, ok := sp.claimed[k]
	if !ok {
		return v, false
	}

	// We want to remove if it is claimed but unpinned
	sp.removeUnpinned(v)
	v.Pin()
	return v, true
}

func (sp *SharedPool[K, V]) removeUnpinned(v V) {
	delete(sp.unpinned, v.Id())
}

func (sp *SharedPool[K, V]) getUnclaimed(k K) (V, bool) {
	v, ok := sp.getUnpinned()
	if !ok {
		return v, false
	}

	sp.removeUnpinned(v)
	return sp.replaceClaimed(k, v)
}

func (sp *SharedPool[K, V]) getUnpinned() (V, bool) {
	for _, v := range sp.unpinned {
		return v, true
	}

	return sp.zero()
}

func (sp *SharedPool[K, V]) waitForRelease(ctx context.Context, k K) (V, bool) {
	select {
	case v := <-sp.in:
		return sp.replaceClaimed(k, v)
	case <-ctx.Done():
		return sp.zero()
	case <-sp.ctx.Done():
		return sp.zero()
	}
}

func (sp *SharedPool[K, V]) replaceClaimed(k K, v V) (V, bool) {
	if oldK, ok := v.Block(); ok {
		delete(sp.claimed, oldK)
	}

	// Change key on v
	_ = v.ReplaceBlock(k)

	sp.claimed[k] = v

	return v, true
}

func (sp *SharedPool[K, V]) zero() (V, bool) {
	var v V
	return v, false
}

func (sp *SharedPool[K, V]) Add(v V) {
	sp.unpinned[v.Id()] = v
}

func (sp *SharedPool[K, V]) Request(ctx context.Context, k K) (V, bool) {
	r := newRequest[K, V](ctx, k)

	select {
	case <-ctx.Done():
		return r.close()
	case sp.requests <- r:
		return r.receive()
	}
}

func (sp *SharedPool[K, V]) Release(v V) {
	sp.in <- v
}

func (sp *SharedPool[K, V]) Iterate() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range sp.unpinned {
			if !yield(v) {
				return
			}
		}
	}
}
