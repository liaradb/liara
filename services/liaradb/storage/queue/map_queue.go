package queue

import (
	"container/list"
	"iter"
)

type MapQueue[K comparable, V any] struct {
	list list.List
	hash map[K]*list.Element
}

type mapTuple[K comparable, V any] struct {
	key   K
	value V
}

func (mq *MapQueue[K, V]) Count() int {
	return mq.list.Len()
}

func (mq *MapQueue[K, V]) Push(k K, v V) {
	if mq.hash == nil {
		mq.hash = make(map[K]*list.Element)
	}
	e := mq.list.PushBack(mapTuple[K, V]{k, v})
	mq.hash[k] = e
}

func (mq *MapQueue[K, V]) Pop() (V, bool) {
	f := mq.list.Front()
	if f == nil {
		return mq.zero()
	}

	t, ok := mq.list.Remove(f).(mapTuple[K, V])
	if ok {
		delete(mq.hash, t.key)
	}
	return t.value, ok
}

func (mq *MapQueue[K, V]) Remove(k K) (V, bool) {
	e, ok := mq.hash[k]
	if !ok {
		return mq.zero()
	}

	t, ok := mq.list.Remove(e).(mapTuple[K, V])
	if ok {
		delete(mq.hash, t.key)
	}
	return t.value, ok
}

func (mq *MapQueue[K, V]) zero() (V, bool) {
	var v V
	return v, false
}

// TODO: Test this
func (mq MapQueue[K, V]) Iterate() iter.Seq[V] {
	return func(yield func(V) bool) {
		e := mq.list.Front()
		for {
			if e == nil || !yield(e.Value.(mapTuple[K, V]).value) {
				return
			}
			e = e.Next()
		}
	}
}
