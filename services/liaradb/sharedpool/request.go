package sharedpool

import "context"

type item[K comparable] interface {
	comparable
	Id() int
	Block() (K, bool)
	ReplaceBlock(K) error
	Pin()
}

type request[K comparable, V item[K]] struct {
	key      K
	ctx      context.Context
	reply    chan V
	canceled bool
}

func newRequest[K comparable, V item[K]](ctx context.Context, k K) *request[K, V] {
	return &request[K, V]{
		key:   k,
		ctx:   ctx,
		reply: make(chan V),
	}
}

func (r *request[K, V]) send(v V) {
	r.reply <- v
}

func (r *request[K, V]) cancel(v V) {
	r.canceled = true
	r.reply <- v
}

func (r *request[K, V]) receive() (V, bool) {
	return <-r.reply, !r.canceled
}

func (r *request[K, V]) close() (V, bool) {
	r.canceled = true
	var v V
	return v, false
}
