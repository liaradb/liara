package queue

import "container/list"

type MapQueue[K comparable, V any] struct {
	list list.List
	hash map[K]*list.Element
}

func (mq *MapQueue[K, V]) Count() int {
	return mq.list.Len()
}

func (mq *MapQueue[K, V]) Push(k K, v V) {
	if mq.hash == nil {
		mq.hash = make(map[K]*list.Element)
	}
	e := mq.list.PushBack(v)
	mq.hash[k] = e
}

func (mq *MapQueue[K, V]) Pop() (V, bool) {
	f := mq.list.Front()
	if f == nil {
		return mq.zero()
	}

	v, ok := mq.list.Remove(f).(V)
	return v, ok
}

func (mq *MapQueue[K, V]) Remove(k K) (V, bool) {
	e, ok := mq.hash[k]
	if !ok {
		return mq.zero()
	}

	v, ok := mq.list.Remove(e).(V)
	return v, ok
}

func (mq *MapQueue[K, V]) zero() (V, bool) {
	var v V
	return v, false
}
