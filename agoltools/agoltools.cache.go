package agoltools

import (
	"appengine"
	"appengine/memcache"
	"time"
)

type Cache struct {
	context appengine.Context
}

func (c *Cache) Set(key string, value interface{}) (err error) {
	item := &memcache.Item{
		Key:    key,
		Object: value,
	}
	return memcache.Gob.Set(c.context, item)
}

func (c *Cache) SetWithExpiration(key string, value interface{}, expiration time.Duration) (err error) {
	item := &memcache.Item{
		Key:        key,
		Object:     value,
		Expiration: expiration,
	}
	return memcache.Gob.Set(c.context, item)
}

func (c *Cache) Get(key string, data interface{}) (err error) {
	_, err = memcache.Gob.Get(c.context, key, data)
	if err != nil {
		return err
	}
	return nil
}

func NewCache(context appengine.Context) *Cache {
	return &Cache{context: context}
}
