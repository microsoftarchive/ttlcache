package ttlcache

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	cache := &Cache{
		ttl:   time.Second,
		items: map[string]*Item{},
	}

	data, exists := cache.Get("hello")
	if exists || data != nil {
		t.Errorf("Expected empty cache to return no data")
	}

	cache.Set("hello", []byte("world"))
	data, exists = cache.Get("hello")
	if !exists {
		t.Errorf("Expected cache to return data for `hello`")
	}
	if string(data) != "world" {
		t.Errorf("Expected cache to return `world` for `hello`")
	}
}

func TestExpiration(t *testing.T) {
	cache := &Cache{
		ttl:   time.Second,
		items: map[string]*Item{},
	}

	cache.Set("x", []byte("1"))
	cache.Set("y", []byte("z"))
	cache.Set("z", []byte("3"))
	cache.startEvictionTimer()

	count := cache.Count()
	if count != 3 {
		t.Errorf("Expected cache to contain 3 items")
	}

	<-time.After(500 * time.Millisecond)
	cache.Lock()
	cache.items["y"].touch(time.Second)
	item, exists := cache.items["x"]
	cache.Unlock()
	if !exists || string(item.data) != "1" || item.expired() {
		t.Errorf("Expected `x` to not have expired after 200ms")
	}

	<-time.After(time.Second)
	cache.RLock()
	_, exists = cache.items["x"]
	if exists {
		t.Errorf("Expected `x` to have expired")
	}
	_, exists = cache.items["z"]
	if exists {
		t.Errorf("Expected `z` to have expired")
	}
	_, exists = cache.items["y"]
	if !exists {
		t.Errorf("Expected `y` to not have expired")
	}
	cache.RUnlock()

	count = cache.Count()
	if count != 1 {
		t.Errorf("Expected cache to contain 1 item")
	}

	<-time.After(600 * time.Millisecond)
	cache.RLock()
	_, exists = cache.items["y"]
	if exists {
		t.Errorf("Expected `y` to have expired")
	}
	cache.RUnlock()

	count = cache.Count()
	if count != 0 {
		t.Errorf("Expected cache to be empty")
	}
}
