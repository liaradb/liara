package storage

import "iter"

type FreePool[K comparable, V any] interface {
	Count() int
	Iterate() iter.Seq[V]
	Pop() (V, bool)
	Push(k K, v V)
	Remove(k K) (V, bool)
}
