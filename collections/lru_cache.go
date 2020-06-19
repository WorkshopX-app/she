package collections

import (
	"container/list"
	"sync"
)

//LRUCacheStats stats about the cache itself
type LRUCacheStats struct {
	Size int64
	Miss int64
	Hit  int64
}

//LRUCache lru cache but safe for concurrent access
//https://github.com/golang/groupcache/blob/master/lru/lru.go
//https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU)
type LRUCache struct {
	// MaxEntries is the maximum number of cache entries before
	// an item is evicted. Zero means no limit.
	MaxEntries int
	mux        sync.Mutex
	stats      LRUCacheStats
	// OnEvicted optionally specifies a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(key Key, value interface{})

	ll    *list.List
	cache map[interface{}]*list.Element
}

// A Key may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators
type Key interface{}

type entry struct {
	key   Key
	value interface{}
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func New(maxEntries int) *LRUCache {
	return &LRUCache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[interface{}]*list.Element),
	}
}

// Add adds a value to the cache.
func (c *LRUCache) Add(key Key, value interface{}) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		return
	}
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		c.RemoveOldest()
	}
}

//Stats return cache stats
func (c *LRUCache) Stats() LRUCacheStats {
	return c.stats
}

// Get looks up a key's value from the cache.
func (c *LRUCache) Get(key Key) (value interface{}, ok bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *LRUCache) Remove(key Key) {
	defer c.mux.Unlock()
	c.mux.Lock()

	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRUCache) RemoveOldest() {
	defer c.mux.Unlock()
	c.mux.Lock()

	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *LRUCache) removeElement(e *list.Element) {
	defer c.mux.Unlock()
	c.mux.Lock()

	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}

// Len returns the number of items in the cache.
func (c *LRUCache) Len() int {
	defer c.mux.Unlock()
	c.mux.Lock()

	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Clear purges all stored items from the cache.
func (c *LRUCache) Clear() {
	if c.OnEvicted != nil {
		for _, e := range c.cache {
			kv := e.Value.(*entry)
			c.OnEvicted(kv.key, kv.value)
		}
	}
	c.ll = nil
	c.cache = nil
}
