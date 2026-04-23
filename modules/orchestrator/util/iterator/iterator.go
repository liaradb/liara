package iterator

import (
	"container/list"
	"iter"
)

func Forward[T any](l *list.List) iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := l.Front(); e != nil; e = e.Next() {
			if !yield(e.Value.(T)) {
				return
			}
		}
	}
}

func Reverse[T any](l *list.List) iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := l.Back(); e != nil; e = e.Prev() {
			if !yield(e.Value.(T)) {
				return
			}
		}
	}
}
