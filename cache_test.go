package shkvcache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/dmi3midd/shkvcache"
	"github.com/stretchr/testify/assert"
)

var (
	validOpts = &shkvcache.Options{
		ShardCount:      8,
		CleanerInterval: 15,
	}
)

func TestNewCache(t *testing.T) {
	c, err := shkvcache.NewCache[string](context.Background(), validOpts)
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestNewCacheDefaultShardCount(t *testing.T) {
	c, err := shkvcache.NewCache[string](context.Background(), validOpts)
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestFailNewCache(t *testing.T) {
	_, err := shkvcache.NewCache[string](context.Background(), &shkvcache.Options{
		ShardCount:      7,
		CleanerInterval: 15,
	})
	assert.Error(t, err)
}

func TestSetGet(t *testing.T) {
	c, _ := shkvcache.NewCache[string](context.Background(), validOpts)
	c.Set("key1", "value1", 0)
	v, ok := c.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", v)
}

func TestFailGet(t *testing.T) {
	c, _ := shkvcache.NewCache[string](context.Background(), validOpts)
	_, ok := c.Get("key1")
	assert.False(t, ok)
}

func TestUpsert(t *testing.T) {
	c, _ := shkvcache.NewCache[string](context.Background(), validOpts)
	c.Set("key1", "value1", 0)
	c.Set("key1", "value2", 0)
	v, ok := c.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value2", v)
}

func TestSetNegativeTTL(t *testing.T) {
	c, _ := shkvcache.NewCache[string](context.Background(), validOpts)
	ok := c.Set("key1", "value1", -1)
	assert.False(t, ok)
}

func TestDel(t *testing.T) {
	c, _ := shkvcache.NewCache[string](context.Background(), validOpts)
	c.Set("key1", "value1", 0)
	c.Del("key1")
	_, ok := c.Get("key1")
	assert.False(t, ok)
}

func TestExpireAndTTL(t *testing.T) {
	c, _ := shkvcache.NewCache[string](context.Background(), validOpts)
	c.Set("key1", "value1", 0)
	ok := c.Expire("key1", 10)
	assert.True(t, ok)
	ttl := c.TTL("key1")
	assert.Equal(t, int64(10), ttl)
}

func TestTTLWithNoExpiration(t *testing.T) {
	c, _ := shkvcache.NewCache[string](context.Background(), validOpts)
	c.Set("key1", "value1", 0)
	ttl := c.TTL("key1")
	assert.Equal(t, int64(-1), ttl)
}

func TestTTLForNonExistentKey(t *testing.T) {
	c, _ := shkvcache.NewCache[string](context.Background(), validOpts)
	ttl := c.TTL("key1")
	assert.Equal(t, int64(-2), ttl)
}

func TestFlush(t *testing.T) {
	c, _ := shkvcache.NewCache[int](context.Background(), validOpts)
	for i := range 4 {
		c.Set(fmt.Sprintf("key%d", i), i, 0)
	}
	c.Flush()
	_, ok := c.Get("key1")
	assert.False(t, ok)
}
