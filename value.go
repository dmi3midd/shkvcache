package shkvcache

import "time"

// Value represents a string value with optional expiration time.
// If expiresAt is 0, the Value has no expiration.
type Value[V any] struct {
	val       V
	expiresAt int64
}

// NewValue creates a new Value.
// Returns false if the Value if ttl is negative.
func newValue[V any](value V, ttl int64) (*Value[V], bool) {
	v := &Value[V]{
		val: value,
	}
	ok := v.expire(ttl)
	if ok {
		return v, true
	}
	return nil, false
}

// GetValue returns the string value of the Value.
func (v *Value[V]) getValue() V {
	return v.val
}

// IsExpired returns true if the Value has an expiration time and is expired.
// Returns false if the Value has no expiration time or is not yet expired.
func (v *Value[V]) isExpired() bool {
	if v.expiresAt == 0 {
		return false
	}
	return time.Now().Unix() > v.expiresAt
}

// Expire sets or updates the expiration time for the Value.
// Returns false if the Value if ttl is negative.
func (v *Value[V]) expire(ttl int64) bool {
	if ttl < 0 {
		return false
	}
	if ttl == 0 {
		v.expiresAt = 0
		return true
	}
	v.expiresAt = time.Now().Unix() + ttl
	return true
}

// TTL returns the remaining time-to-live in seconds.
// Returns 0 if the Value has no expiration time.
func (v *Value[V]) ttl() int64 {
	if v.expiresAt == 0 {
		return 0
	}
	return v.expiresAt - time.Now().Unix()
}
