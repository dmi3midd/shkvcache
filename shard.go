package shkvcache

import (
	"sync"
)

// Shard is a thread-safe unit of storage containing a map of items protected by an RWMutex.
type shard[V any] struct {
	mu    sync.RWMutex
	items map[string]V
}

// NewShard creates a new shard.
func newShard[V any]() *shard[V] {
	return &shard[V]{
		mu:    sync.RWMutex{},
		items: make(map[string]V, 8),
	}
}

// get retrieves an item for the given key from the shard.
// Returns zero value and false if the key does not exist.
func (s *shard[V]) get(key string) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	return value, true
}

// set sets or updates the item for the given key in the shard.
func (s *shard[V]) set(key string, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = value
}

// getOrSet retrieves an item for the given key, or sets and stores a new item if it does not exist.
// Returns the item and true if it existed, or the set item and false if it was newly created.
func (s *shard[V]) getOrSet(key string, initFn func() V) (V, bool) {
	s.mu.RLock()
	val, ok := s.items[key]
	s.mu.RUnlock()
	if ok {
		return val, true
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if val, ok = s.items[key]; ok {
		return val, true
	}
	val = initFn()
	s.items[key] = val
	return val, false
}

// update updates an item in the Shard under write lock if it exists.
// Returns true if the item was found and the update function returned true.
// Returns false if the item does not exist.
func (s *shard[V]) update(key string, fn func(val V) bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.items[key]
	if !ok {
		return false
	}
	if fn(val) {
		s.items[key] = val
		return true
	}
	return false
}

// delete removes the item for the given key from the shard.
func (s *shard[V]) delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, key)
}

// flush removes all items from the shard.
func (s *shard[V]) flush(shardCount int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make(map[string]V, shardCount)
}

// clean removes items from the shard for which the predicate returns true.
func (s *shard[V]) clean(predicate func(key string, value V) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range s.items {
		if predicate(k, v) {
			delete(s.items, k)
		}
	}
}
