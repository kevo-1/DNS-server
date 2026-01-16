package models

import "time"

type CacheEntry struct {
	Domain    string
	IPAddress string
	TTL       time.Duration
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

type CacheStatistics struct {
	Hits          int64
	Misses        int64
	Evictions     int64
	TotalEntries  int
	TotalCapacity int
}

func (s *CacheStatistics) HitRate() float64 {
	total := s.Hits + s.Misses
	if total == 0 {
		return 0.0
	}
	return float64(s.Hits) / float64(total)
}

type CacheConfig struct {
	MaxEntries     int
	DefaultTTL     time.Duration
	CleanupInterval time.Duration
	EnableStats    bool
}

func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		MaxEntries:      1000,
		DefaultTTL:      5 * time.Minute,
		CleanupInterval: 1 * time.Minute,
		EnableStats:     true,
	}
}
