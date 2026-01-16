package resolver

import (
	"DNS-server/internal/protocol"
	"DNS-server/models"
	"errors"
	"sync"
)

var (
	ErrResolutionFailed = errors.New("DNS resolution failed")
	ErrInvalidDomain    = errors.New("invalid domain name")
)

type Resolver struct {
	cache              *DNSCache
	iterativeResolver  *IterativeResolver
	mu                 sync.RWMutex
}

var (
	instance *Resolver
	once     sync.Once
)

func GetInstance() *Resolver {
	once.Do(func() {
		cache := NewDNSCache(models.DefaultCacheConfig())
		instance = &Resolver{
			cache:             cache,
			iterativeResolver: NewIterativeResolver(cache),
		}
	})
	return instance
}

func NewResolver(config *models.CacheConfig) *Resolver {
	cache := NewDNSCache(config)
	return &Resolver{
		cache:             cache,
		iterativeResolver: NewIterativeResolver(cache),
	}
}

func (r *Resolver) Resolve(domain string, recordType uint16) (string, error) {
	if domain == "" {
		return "", ErrInvalidDomain
	}

	if ip, found := r.cache.Get(domain); found {
		return ip, nil
	}

	ip, err := r.iterativeResolver.Resolve(domain, recordType)
	if err != nil {
		return "", ErrResolutionFailed
	}

	return ip, nil
}

func (r *Resolver) ResolveA(domain string) (string, error) {
	return r.Resolve(domain, protocol.TypeA)
}

func (r *Resolver) LookupCache(domain string) (string, bool) {
	return r.cache.Get(domain)
}

func (r *Resolver) UpdateCache(domain, ip string, ttl int) {
	r.cache.Set(domain, ip, protocol.ParseTTL(ttl))
}

func (r *Resolver) GetStats() models.CacheStatistics {
	return r.cache.GetStats()
}

func (r *Resolver) ClearCache() {
	r.cache.Clear()
}

func (r *Resolver) Close() {
	r.cache.Close()
}

func RetrieveIP(url models.URL) (string, error) {
	resolver := GetInstance()
	domain := url.Authority.Host
	return resolver.ResolveA(domain)
}
