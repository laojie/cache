// Package memcached use the "github.com/bradfitz/gomemcache/memcache" as cache driver
package memcached

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gookit/cache"
)

// MemCached definition
type MemCached struct {
	servers []string
	client  *memcache.Client
}

// New a MemCached instance
func New(servers ...string) *MemCached {
	return &MemCached{
		servers: servers,
	}
}

// Connect new a MemCached instance and connect to memcached servers.
func Connect(servers ...string) *MemCached {
	c := &MemCached{
		servers: servers,
	}

	return c.Connect()
}

// Connect to servers
func (c *MemCached) Connect() *MemCached {
	c.client = memcache.New(c.servers...)
	return c
}

// Has cache key
func (c *MemCached) Has(key string) bool {
	_, err := c.client.Get(key)
	return err == nil
}

// Get value by key
func (c *MemCached) Get(key string) (val interface{}) {
	item, err := c.client.Get(key)
	if err != nil {
		return
	}

	err = cache.GobDecode(item.Value, val)
	if err != nil {
		return nil
	}

	return
}

// Set value by key
func (c *MemCached) Set(key string, val interface{}, ttl time.Duration) (err error) {
	bts, err := cache.GobEncode(val)
	if err != nil {
		return err
	}

	return c.client.Set(&memcache.Item{
		Key:   key,
		Value: bts,
		// expire time. 0 is never
		Expiration: int32(ttl / time.Second),
	})
}

// Del value by key
func (c *MemCached) Del(key string) error {
	return c.client.Delete(key)
}

// GetMulti values by multi key
func (c *MemCached) GetMulti(keys []string) map[string]interface{} {
	items, err := c.client.GetMulti(keys)
	if err != nil {
		return nil
	}

	values := make(map[string]interface{}, len(keys))

	for key, item := range items {
		var val interface{}
		if err := cache.GobDecode(item.Value, val); err != nil {
			return nil
		}

		values[key] = val
	}

	return values
}

// SetMulti values by multi key
func (c *MemCached) SetMulti(values map[string]interface{}, ttl time.Duration) (err error) {
	for key, val := range values {
		if err = c.Set(key, val, ttl); err != nil {
			return
		}
	}

	return
}

// DelMulti values by multi key
func (c *MemCached) DelMulti(keys []string) error {
	for _, key := range keys {
		if err := c.client.Delete(key); err != nil {
			return err
		}
	}

	return nil
}

// Clear all caches
func (c *MemCached) Clear() error {
	return c.client.DeleteAll()
}

// Close driver
func (*MemCached) Close() error {
	return nil
}

// Client get
func (c *MemCached) Client() *memcache.Client {
	return c.client
}
