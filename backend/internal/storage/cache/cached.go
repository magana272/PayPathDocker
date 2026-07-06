package cache

import (
	"time"

	"golang.org/x/sync/singleflight"
)

const (
	DataCacheTTL  = 2 * time.Minute
	TokenCacheTTL = 5 * time.Minute
)

func CachedList[T any](c *Cache, sf *singleflight.Group, key string, loader func() ([]T, error)) ([]T, error) {
	if v, ok := c.Get(key); ok {
		return v.([]T), nil
	}
	v, err, _ := sf.Do(key, func() (any, error) {
		if v, ok := c.Get(key); ok {
			return v, nil
		}
		list, err := loader()
		if err != nil {
			return nil, err
		}
		c.Set(key, list, DataCacheTTL)
		return list, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]T), nil
}
