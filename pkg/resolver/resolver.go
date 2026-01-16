package core

import (
	"DNS-server/models"
	"strings"
	"sync"
)

type Resolver struct {
	cache map[string]string
	mu    sync.Mutex
}

func (r *Resolver) lookupCache(domain string) (string, bool) {
	//lock the cache as it is a critical section
	r.mu.Lock()
	defer r.mu.Unlock()

	if ip, ok := r.cache[domain]; ok {
		return ip, true
	}
	return "", false
}

func (r *Resolver) updateCache(domain, ip string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache[domain] = ip
}

var (
	instance *Resolver
	once     sync.Once
)

func GetInstance() *Resolver {
	once.Do(func() {
		instance = &Resolver{
			cache: make(map[string]string),
		}
	})
	return instance
}

func retrieveIP(url models.URL) (string, error) {
	resolver := GetInstance()
	domain := strings.Split(url.Authority.Host, ".")

	if ip, ok := resolver.lookupCache(domain[len(domain)-1]); ok {
		return ip, nil
	}

	return "", nil
}
