package ttlcache

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	cache := &Cache{
		items: map[string]*Item{},
	}

	data, exists := cache.Get("hello", true)
	if exists || data != "" {
		t.Errorf("Expected empty cache to return no data")
	}

	cache.Set("hello", "world", time.Second)
	data, exists = cache.Get("hello", true)
	if !exists {
		t.Errorf("Expected cache to return data for `hello`")
	}
	if data != "world" {
		t.Errorf("Expected cache to return `world` for `hello`")
	}
	if cache.GetCounter() != 1 {
		t.Errorf("Expected cache get counter is equal to 1")
	}
}

func TestDelete(t *testing.T) {
	cache := &Cache{
		items: map[string]*Item{},
	}
	cache.Set("Test", "Delete", time.Second)
	_, exists := cache.Get("Test", true)
	if !exists {
		t.Errorf("Expected cache to return data for `Test`")
	}
	cache.Delete("Test")
	_, exists = cache.Get("Test", true)
	if exists {
		t.Errorf("Expected cache to delete data for `Test`")
	}

}

func TestExpiration(t *testing.T) {
	cache := &Cache{
		items: map[string]*Item{},
	}

	cache.Set("x", "1", time.Second)
	cache.Set("y", 123, time.Second)
	cache.Set("z", time.Second, time.Second)
	cache.startCleanupTimer()

	count := cache.Count()
	if count != 3 {
		t.Errorf("Expected cache to contain 3 items")
	}

	<-time.After(500 * time.Millisecond)
	cache.mutex.Lock()
	cache.items["y"].touch()
	item, exists := cache.items["x"]
	cache.mutex.Unlock()
	if !exists || item.data != "1" || item.expired() {
		t.Errorf("Expected `x` to not have expired after 200ms")
	}

	<-time.After(time.Second)
	cache.mutex.RLock()
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
	cache.mutex.RUnlock()

	count = cache.Count()
	if count != 1 {
		t.Errorf("Expected cache to contain 1 item")
	}

	<-time.After(600 * time.Millisecond)
	cache.mutex.RLock()
	_, exists = cache.items["y"]
	if exists {
		t.Errorf("Expected `y` to have expired")
	}
	cache.mutex.RUnlock()

	count = cache.Count()
	if count != 0 {
		t.Errorf("Expected cache to be empty")
	}
}
