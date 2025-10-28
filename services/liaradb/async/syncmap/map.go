package syncmap

import (
	"iter"
	"sync"
)

type Map[K comparable, V any] struct {
	hash map[K]V
	lock sync.RWMutex
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	value, ok := m.hash[key]
	return value, ok
}

func (m *Map[K, V]) Set(key K, value V) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.hash == nil {
		m.hash = make(map[K]V)
	}
	m.hash[key] = value
}

func (m *Map[K, V]) Iterate() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.lock.RLock()
		defer m.lock.RUnlock()

		for key, value := range m.hash {
			if !yield(key, value) {
				return
			}
		}
	}
}
