package cache

import "fmt"

type DerivedCache struct {
	cache *Cache
}

func NewDerivedCache(c *Cache) *DerivedCache {
	return &DerivedCache{cache: c}
}

func (d *DerivedCache) Get(userID int, key string) (any, bool) {
	return d.cache.Get(fmt.Sprintf("derived:%d:%s", userID, key))
}

func (d *DerivedCache) Set(userID int, key string, val any) {
	d.cache.Set(fmt.Sprintf("derived:%d:%s", userID, key), val, DataCacheTTL)
}

func (d *DerivedCache) Invalidate(userID int) {
	d.cache.DeletePrefix(fmt.Sprintf("derived:%d:", userID))
}
