package shkvcache

import "errors"

var (
	ErrInvalidShardCount      = errors.New("invalid shard count")
	ErrInvalidCleanerInterval = errors.New("invalid cleaner interval")
)
