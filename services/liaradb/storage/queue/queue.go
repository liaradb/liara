package queue

import "container/list"

type Queue[K comparable, V any] struct {
	list list.List
	hash map[K]*list.Element
}

func (q *Queue[K, V]) Push(k K, v V) {
	if q.hash == nil {
		q.hash = make(map[K]*list.Element)
	}
	e := q.list.PushBack(v)
	q.hash[k] = e
}

func (q *Queue[K, V]) Pop() (V, bool) {
	f := q.list.Front()
	if f == nil {
		return q.zero()
	}

	v, ok := q.list.Remove(f).(V)
	return v, ok
}

func (q *Queue[K, V]) Remove(k K) (V, bool) {
	e, ok := q.hash[k]
	if !ok {
		return q.zero()
	}

	v, ok := q.list.Remove(e).(V)
	return v, ok
}

func (q *Queue[K, V]) zero() (V, bool) {
	var v V
	return v, false
}
