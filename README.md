# go_lru_cache

A small, thread-safe, fixed-capacity LRU cache for Go.

## Usage

```go
package main

import (
	"fmt"

	lru "github.com/yourname/go_lru_cache"
)

func main() {
	c := lru.NewCache[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Get("a") // "a" is now the most recently used
	c.Put("c", 3) // evicts "b"

	if v, ok := c.Get("b"); ok {
		fmt.Println("unexpected:", v)
	} else {
		fmt.Println("b was evicted")
	}

	fmt.Println(c.Len()) // 2
}
```

## Why this exists

Go's standard library has no LRU cache. A map alone cannot enforce a size bound, and
naively clearing the map when it fills up throws away entries that are still hot. This
library keeps insertion and lookup O(1) by pairing a map with a doubly linked list:
the map gives direct access to entries, and the list tracks recency order.

The trade-off is memory overhead: every entry stores a list element and a map entry,
roughly doubling the bookkeeping cost compared to a plain map. For workloads where the
working set fits comfortably in memory and eviction must be cheap, that is a good deal.

## Concurrency

All public methods are safe for concurrent use. A single mutex guards both the map and
the list, so a `Get` briefly blocks other operations. This is fine for typical cache
sizes; if you need higher throughput under heavy read contention, a sharded cache would
be the next step, but that is deliberately out of scope here.

## Edge cases

- `NewCache` panics if capacity is not positive. A cache with zero capacity cannot
  hold anything, and silently accepting one would hide a configuration mistake.
- `Get` on a missing key returns the zero value of `V` and `false`. It does not
  insert the key.
- Updating an existing key does not change the cache size and makes that key the
  most recently used.
- The zero value of `Cache` is not usable; always construct with `NewCache`.
