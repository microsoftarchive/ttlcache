package ttlcache

import (
	"sync"
	"time"
)

// Cache is a synchronised map of items that auto-expire once stale
type Cache struct {
	sync.RWMutex
	ttl   time.Duration
	items map[string]*Item
}

// NewCache is a helper to create instance of the Cache struct
func New(duration time.Duration) *Cache {
	cache := &Cache{
		ttl:   duration,
		items: map[string]*Item{},
	}
	cache.startEvictionTimer()
	return cache
}

// Set is a thread-safe way to add new items to the map
func (cache *Cache) Set(key string, data []byte) {
	cache.Lock()
	defer cache.Unlock()

	item := &Item{data: data}
	item.touch(cache.ttl)
	cache.items[key] = item
}

// Get is a thread-safe way to lookup items
// Every lookup, also touches the item, hence extending it's life
func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.Lock()
	defer cache.Unlock()

	item, exists := cache.items[key]
	if !exists || item.expired() {
		return nil, false
	}

	item.touch(cache.ttl)
	return item.data, true
}

// Count returns the number of items in the cache
func (cache *Cache) Count() int {
	cache.RLock()
	defer cache.RUnlock()

	count := len(cache.items)
	return count
}

func (cache *Cache) startEvictionTimer() {
	duration := cache.ttl
	if duration < time.Second {
		duration = time.Second
	}

	ticker := time.Tick(duration)
	go (func() {
		for {
			select {
			case <-ticker:
				cache.evict()
			}
		}
	})()
}

func (cache *Cache) evict() {
	cache.Lock()
	defer cache.Unlock()

	for key, item := range cache.items {
		if item.expired() {
			delete(cache.items, key)
		}
	}
}
