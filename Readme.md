## TTLCache - an in-memory cache with expiration

TTLCache is a minimal wrapper over a map of custom in golang, entries of which are

1. Thread-safe
2. Auto-Expiring after a certain time
3. Managed auto-extending expiration on `Get`s

[![Build Status](https://travis-ci.org/ikoroteev/ttlcache.svg)](https://travis-ci.org/ikoroteev/ttlcache)

#### Usage
```go
import (
  "time"
  "github.com/ikoroteev/ttlcache"
)

func main () {
  cache := ttlcache.NewCache()
  cache.Set("key", "value", time.Second)
  cache.Set("key1", 24, time.Duration(500) * time.Millisecond)
  cache.Set("key3", time.Second, time.Second)
  value, exists := cache.Get("key", true) // true - extend cache ttl, otherwise false
  count := cache.Count()
}
```