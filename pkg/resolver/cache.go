package resolver

import (
	"DNS-server/models"
	"container/list"
	"sync"
	"time"
)

type DNSCache struct {
	config     *models.CacheConfig
	entries    map[string]*cacheNode
	lruList    *list.List
	mu         sync.RWMutex
	stats      models.CacheStatistics
	stopCleanup chan bool
}

type cacheNode struct {
	entry   *models.CacheEntry
	element *list.Element
}

func NewDNSCache(config *models.CacheConfig) *DNSCache {
	if config == nil {
		config = models.DefaultCacheConfig()
	}

	cache := &DNSCache{
		config:      config,
		entries:     make(map[string]*cacheNode),
		lruList:     list.New(),
		stopCleanup: make(chan bool),
	}

	go cache.cleanupExpired()

	return cache
}

func (c *DNSCache) Get(domain string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.entries[domain]
	if !exists {
		if c.config.EnableStats {
			c.stats.Misses++
		}
		return "", false
	}

	if node.entry.IsExpired() {
		c.removeNode(domain)
		if c.config.EnableStats {
			c.stats.Misses++
		}
		return "", false
	}

	c.lruList.MoveToFront(node.element)
	
	if c.config.EnableStats {
		c.stats.Hits++
	}

	return node.entry.IPAddress, true
}

func (c *DNSCache) Set(domain, ipAddress string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.entries[domain]; exists {
		node.entry.IPAddress = ipAddress
		node.entry.TTL = ttl
		node.entry.ExpiresAt = time.Now().Add(ttl)
		c.lruList.MoveToFront(node.element)
		return
	}

	if c.lruList.Len() >= c.config.MaxEntries {
		c.evictLRU()
	}

	entry := &models.CacheEntry{
		Domain:    domain,
		IPAddress: ipAddress,
		TTL:       ttl,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
	}

	element := c.lruList.PushFront(domain)
	c.entries[domain] = &cacheNode{
		entry:   entry,
		element: element,
	}

	if c.config.EnableStats {
		c.stats.TotalEntries = len(c.entries)
		c.stats.TotalCapacity = c.config.MaxEntries
	}
}

func (c *DNSCache) evictLRU() {
	element := c.lruList.Back()
	if element != nil {
		domain := element.Value.(string)
		c.removeNode(domain)
		if c.config.EnableStats {
			c.stats.Evictions++
		}
	}
}

func (c *DNSCache) removeNode(domain string) {
	if node, exists := c.entries[domain]; exists {
		c.lruList.Remove(node.element)
		delete(c.entries, domain)
		if c.config.EnableStats {
			c.stats.TotalEntries = len(c.entries)
		}
	}
}

func (c *DNSCache) cleanupExpired() {
	ticker := time.NewTicker(c.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.removeExpiredEntries()
		case <-c.stopCleanup:
			return
		}
	}
}

func (c *DNSCache) removeExpiredEntries() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for domain, node := range c.entries {
		if node.entry.IsExpired() {
			c.lruList.Remove(node.element)
			delete(c.entries, domain)
		}
	}

	if c.config.EnableStats {
		c.stats.TotalEntries = len(c.entries)
	}
}

func (c *DNSCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*cacheNode)
	c.lruList = list.New()
	c.stats.TotalEntries = 0
}

func (c *DNSCache) GetStats() models.CacheStatistics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.stats
}

func (c *DNSCache) Close() {
	close(c.stopCleanup)
}
