package go_lru_cache

import (
	"container/list"
	"sync"
)

// Cache is a thread-safe LRU cache with a fixed capacity.
// When the cache is full, adding a new entry evicts the entry that was accessed
// least recently. The zero value is not usable; use NewCache to create one.
type Cache[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	items    map[K]*list.Element
	order    *list.List
}

// entry holds a key-value pair stored in the list.
type entry[K comparable, V any] struct {
	key   K
	value V
}

// NewCache creates a new LRU cache with the given capacity.
// Capacity must be positive; if it is not, NewCache panics.
func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity <= 0 {
		panic("capacity must be positive")
	}
	return &Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]*list.Element),
		order:    list.New(),
	}
}

// Get returns the value associated with key and whether it was found.
// If found, the entry is marked as most recently used.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		return elem.Value.(*entry[K, V]).value, true
	}
	var zero V
	return zero, false
}

// Put adds or updates the value for key.
// If the key already exists, its value is updated and it becomes the most
// recently used entry. If the cache is at capacity, the least recently used
// entry is evicted before inserting the new one.
func (c *Cache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		elem.Value.(*entry[K, V]).value = value
		c.order.MoveToFront(elem)
		return
	}

	if c.order.Len() >= c.capacity {
		c.evict()
	}

	elem := c.order.PushFront(&entry[K, V]{key: key, value: value})
	c.items[key] = elem
}

// Len returns the number of entries currently in the cache.
func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.order.Len()
}

// evict removes the least recently used entry.
// It must be called with the lock held and only when the cache is non-empty.
func (c *Cache[K, V]) evict() {
	back := c.order.Back()
	if back == nil {
		return
	}
	c.order.Remove(back)
	delete(c.items, back.Value.(*entry[K, V]).key)
}
