package main

import (
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	value      interface{} //meaning it can hold any type of value (e.g., string, int, struct, etc.).
	expiration time.Time
}

type Cache struct {
	sync.RWMutex
	items   map[string]CacheItem
	maxSize int
	janitor *time.Ticker //triggers periodic cleanup of expired cache items.

}

//On each tick, the cleanup method iterates through all items in the cache
// and deletes those whose expiration time has passed.
func (c *Cache) cleanup() {
	for range c.janitor.C {
		c.Lock()
		for key, item := range c.items {
			if time.Now().After(item.expiration) {
				delete(c.items, key)
			}
		}
	}
}

// Initializes and returns a new Cache instance.
func NewCache(maxSize int, cleanupInterval time.Duration) *Cache {
	cache := &Cache{
		items:   make(map[string]CacheItem),
		maxSize: maxSize,
		janitor: time.NewTicker(cleanupInterval),
	}

	go cache.cleanup()
	return cache
}

// Adds a new item to the cache or updates an existing item.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	c.Lock()
	defer c.Unlock()

	if len(c.items) >= c.maxSize {
		return fmt.Errorf("cache is full")
	}

	c.items[key] = CacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}

	return nil

}

// Retrieves an item from the cache.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	item, exists := c.items[key]

	if !exists {
		return nil, false
	}
	if time.Now().After(item.expiration) {
		return nil, false
	}
	return item.value, true
}

// Stops the janitor's cleanup process
func (c *Cache) Close() {
	c.janitor.Stop()
}

func main() {

	cache := NewCache(100, 5*time.Second)
	defer cache.Close()

	cache.Set("key1", "value1", 10*time.Second)
	cache.Set("key2", "value2", 2*time.Second)

	val, exists := cache.Get("key1")

	if exists {
		fmt.Printf("Key1: %v\n", val)

	} else {
		fmt.Println("Key1 expired")
	}
	time.Sleep(3 * time.Second)

	if val, exists := cache.Get("key2"); exists {
		fmt.Printf("Key2:%v\n", val)
	} else {
		fmt.Println("Key2 expired")
	}

}
