package di

import "sync"

type syncMap[K comparable, V any] struct {
	mu     sync.RWMutex
	values map[K]V
}

func (m *syncMap[K, V]) Load(key K) (value V, ok bool) {
	m.mu.RLock()
	value, ok = m.values[key]
	m.mu.RUnlock()
	return
}

func (m *syncMap[K, V]) Has(key K) (ok bool) {
	m.mu.RLock()
	_, ok = m.values[key]
	m.mu.RUnlock()
	return
}

func (m *syncMap[K, V]) Store(key K, value V) {
	m.mu.Lock()
	if m.values == nil {
		m.values = make(map[K]V)
	}
	m.values[key] = value
	m.mu.Unlock()
}

func (m *syncMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	m.mu.Lock()
	if m.values == nil {
		m.values = make(map[K]V)
	}

	previous, loaded = m.values[key]
	m.values[key] = value
	m.mu.Unlock()
	return
}

func (m *syncMap[K, V]) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.values)
}
