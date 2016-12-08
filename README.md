## TTLCache - an in-memory LRU cache with expiration

TTLCache is a minimal wrapper over a string map in golang, entries of which are

1. Thread-safe
2. Auto-Expiring after a certain time
3. Auto-Extending expiration on `Get`s

[![Build Status](https://travis-ci.org/gospackler/ttlcache.svg)](https://travis-ci.org/gospackler/ttlcache)

#### Usage
```go
import (
  "time"
  "github.com/wunderlist/ttlcache"
)

func main () {
  cache := ttlcache.New(time.Second)
  // ttlcache stores byte slice values.
  cache.Set("key", []byte("value"))
  value, exists := cache.Get("key")
  count := cache.Count()
}
```
