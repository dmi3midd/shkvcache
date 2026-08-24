package shkvcache

import (
	"context"
	"time"
)

// Cleaner executes a periodic cleanup function on a fixed interval until stopped or canceled.
type cleaner struct {
	interval  int // in seconds
	cleanFunc func()
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewCleaner creates a new cleaner with the given interval and cleanup function.
// If parentCtx is nil, context.Background() is used.
func NewCleaner(parentCtx context.Context, interval int, cleanFunc func()) *cleaner {
	if parentCtx == nil {
		parentCtx = context.Background()
	}
	ctx, cancel := context.WithCancel(parentCtx)
	return &cleaner{
		interval:  interval,
		cleanFunc: cleanFunc,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start runs the cleaner background loop in a new goroutine.
func (c *cleaner) Start() {
	go c.run()
}

// Stop stops the background cleaner loop.
func (c *cleaner) Stop() {
	c.cancel()
}

// run executes cleanFunc periodically on every ticker interval until context is canceled.
func (c *cleaner) run() {
	ticker := time.NewTicker(time.Duration(c.interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanFunc()
		case <-c.ctx.Done():
			return
		}
	}
}
