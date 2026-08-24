package shkvcache

// Options holds the configuration options for the Cache.
type Options struct {
	ShardCount      int  // number of shards to use
	CleanerInterval int  // interval in seconds for cleaner to run
	RunCleaner      bool // run cleaner goroutine on startup
}

// Validate validates the options.
func (o *Options) Validate() error {
	if !isPowerOfTwo(o.ShardCount) {
		return ErrInvalidShardCount
	}
	if o.RunCleaner && o.CleanerInterval <= 0 {
		return ErrInvalidCleanerInterval
	}
	return nil
}

// DefaultOptions returns the default options for the Cache.
func DefaultOptions() *Options {
	return &Options{
		ShardCount:      8,
		CleanerInterval: 15,
		RunCleaner:      true,
	}
}
