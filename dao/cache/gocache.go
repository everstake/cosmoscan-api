package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type Cache struct {
	cache *cache.Cache
}

func New() *Cache {
	return &Cache{
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (c *Cache) CacheSet(key string, data interface{}, duration time.Duration) {
	c.cache.Set(key, data, duration)
}

func (c *Cache) CacheGet(key string) (data interface{}, found bool) {
	return c.cache.Get(key)
}
