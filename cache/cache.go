package cache

import (
	"awesome-cache/lru"
	"sync"
)

type cache struct {
	mtx       sync.Mutex
	lruCache  *lru.Cache
	cacheByte int64
}

func (c *cache) add(key string, value ByteView) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.lruCache == nil {
		c.lruCache = lru.NewCache(c.cacheByte, nil)
	}
	c.lruCache.Add(key, value)
}

func (c *cache) get(key string) (byteView ByteView, ok bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.lruCache == nil {
		return
	}

	if value, ok := c.lruCache.Get(key) ; ok {
		return value.(ByteView), true
	}

	return
}
