package cache

import (
	"strings"
	"sync"
	"time"
)

type entry struct {
	value     any
	expiresAt time.Time
}

type Cache struct {
	data sync.Map
}

func New() *Cache {
	c := &Cache{}
	go c.reap()
	return c
}

func (c *Cache) Get(key string) (any, bool) {
	v, ok := c.data.Load(key)
	if !ok {
		return nil, false
	}
	e := v.(*entry)
	if time.Now().After(e.expiresAt) {
		c.data.Delete(key)
		return nil, false
	}
	return e.value, true
}

func (c *Cache) Set(key string, val any, ttl time.Duration) {
	c.data.Store(key, &entry{value: val, expiresAt: time.Now().Add(ttl)})
}

func (c *Cache) Delete(key string) {
	c.data.Delete(key)
}

func (c *Cache) DeletePrefix(prefix string) {
	c.data.Range(func(key, _ any) bool {
		if k, ok := key.(string); ok && strings.HasPrefix(k, prefix) {
			c.data.Delete(key)
		}
		return true
	})
}

func (c *Cache) reap() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		now := time.Now()
		c.data.Range(func(key, value any) bool {
			if e, ok := value.(*entry); ok && now.After(e.expiresAt) {
				c.data.Delete(key)
			}
			return true
		})
	}
}
