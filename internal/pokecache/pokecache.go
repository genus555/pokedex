package pokecache
import (
	"fmt"
	"time"
	"sync"
)

type cacheEntry struct {
	createdAt	time.Time
	val			[]byte
}

type Cache struct {
	cache		map[string]cacheEntry
	mu			sync.Mutex
	interval	time.Duration
}

func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{
		createdAt:		time.Now(),
		val:			val,
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, exists := c.cache[key]
	return entry.val, exists
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache:		make(map[string]cacheEntry),
		interval:	interval,
	}
	go c.reapLoop()
	return c
}

func (c *Cache)  reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.cache {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}

func ImportCheck() {
	fmt.Println("Imported Successfully")
}