package collection

import (
	"iter"
	"sync"
)

type SyncMap[K comparable, V any] struct {
	hash map[K]V
	lock sync.RWMutex
}

func (sh *SyncMap[K, V]) Get(key K) (V, bool) {
	sh.lock.RLock()
	defer sh.lock.RUnlock()

	value, ok := sh.hash[key]
	return value, ok
}

func (sh *SyncMap[K, V]) Set(key K, value V) {
	sh.lock.Lock()
	defer sh.lock.Unlock()

	if sh.hash == nil {
		sh.hash = make(map[K]V)
	}
	sh.hash[key] = value
}

func (sh *SyncMap[K, V]) Iterate() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		sh.lock.RLock()
		defer sh.lock.RUnlock()

		for key, value := range sh.hash {
			if !yield(key, value) {
				return
			}
		}
	}
}
