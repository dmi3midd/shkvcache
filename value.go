package shkvcache

import "time"

// value represents a string value with optional expiration time.
// If expiresAt is 0, the value has no expiration.
type value[V any] struct {
	val       V
	expiresAt int64
}

// newValue creates a new value.
// Returns false if the value if ttl is negative.
func newValue[V any](val V, ttl int64) (*value[V], bool) {
	v := &value[V]{
		val: val,
	}
	ok := v.expire(ttl)
	if ok {
		return v, true
	}
	return nil, false
}

// getValue returns the string value of the value.
func (v *value[V]) getValue() V {
	return v.val
}

// isExpired returns true if the value has an expiration time and is expired.
// Returns false if the value has no expiration time or is not yet expired.
func (v *value[V]) isExpired() bool {
	if v.expiresAt == 0 {
		return false
	}
	return time.Now().Unix() > v.expiresAt
}

// expire sets or updates the expiration time for the value.
// Returns false if the value if ttl is negative.
func (v *value[V]) expire(ttl int64) bool {
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

// ttl returns the remaining time-to-live in seconds.
// Returns 0 if the value has no expiration time.
func (v *value[V]) ttl() int64 {
	if v.expiresAt == 0 {
		return 0
	}
	return v.expiresAt - time.Now().Unix()
}
