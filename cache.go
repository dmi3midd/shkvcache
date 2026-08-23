package shkvcache

import "hash/fnv"

// Cache is a partitioned map storing string key-value pairs across multiple shards.
type Cache[V any] struct {
	shardCount int
	shards     []*Shard[*Value[V]]
}

// NewCache creates a new Cache with default shard count.
func NewCache[V any](shardCount int) (*Cache[V], error) {
	if !isPowerOfTwo(shardCount) {
		return nil, ErrInvalidShardCount
	}
	shards := make([]*Shard[*Value[V]], shardCount)
	for i := range shards {
		shards[i] = NewShard[*Value[V]]()
	}
	return &Cache[V]{
		shardCount: shardCount,
		shards:     shards,
	}, nil
}

// getShard returns the shard corresponding to the given key based on FNV-1a hash.
func (s *Cache[V]) getShard(key string) *Shard[*Value[V]] {
	h := fnv.New32a()
	h.Write([]byte(key))
	idx := h.Sum32() & uint32(s.shardCount-1)
	return s.shards[idx]
}

// Get retrieves the value of the key from the sharded map.
// Returns zero value and false if the key does not exist or is expired.
func (s *Cache[V]) Get(key string) (V, bool) {
	shard := s.getShard(key)
	i, ok := shard.Get(key)
	var zero V
	if !ok {
		return zero, false
	}
	if i.isExpired() {
		shard.Delete(key)
		return zero, false
	}
	return i.getValue(), true
}

// Set sets or updates the value and optional TTL for the key in the sharded map.
func (s *Cache[V]) Set(key string, value V, ttl int64) bool {
	i := newValue(value, ttl)
	s.getShard(key).Set(key, i)
	return true
}

// Delete removes the key from the sharded map.
func (s *Cache[V]) Del(key string) {
	s.getShard(key).Delete(key)
}

// Flush removes all items across all shards.
func (s *Cache[V]) Flush() {
	for _, shard := range s.shards {
		shard.Flush(s.shardCount)
	}
}

// Expire sets the expiration time for the key.
// Returns true if the key exists and expiration was set.
// Returns false if the key does not exist or ttl is invalid.
func (s *Cache[V]) Expire(key string, ttl int64) bool {
	return s.getShard(key).Update(key, func(i *Value[V]) bool {
		// if i.isExpired() {
		// 	return false
		// }
		ok := i.expire(ttl)
		return ok
	})
}

// TTL returns the ttl of the key in seconds.
// Returns -1 if the key exists and has no expiration time.
// Returns -2 if the key does not exist or is expired.
func (s *Cache[V]) TTL(key string) int64 {
	i, ok := s.getShard(key).Get(key)
	if !ok || i.isExpired() {
		return -2
	}
	if i.ttl() == 0 {
		return -1
	}
	return i.ttl()
}

// CleanExpired removes all expired items across all shards.
func (s *Cache[V]) CleanExpired() {
	for _, shard := range s.shards {
		shard.Clean(func(key string, value *Value[V]) bool {
			return value.isExpired()
		})
	}
}
