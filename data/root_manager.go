package data

import (
	"sync"
)

type RootServerManager struct {
	mu 	sync.RWMutex
	servers	[]NameServer
	ipToServer 	map[string]*NameServer 
	nameToServer map[string]*NameServer
}

var (
	instance	*RootServerManager
	managerOnce	sync.Once
)

func GetRootServerManager() *RootServerManager {
	managerOnce.Do(func() {
		instance = newRootServerManager()
	})
	return instance
}

func newRootServerManager() *RootServerManager {
	manager := &RootServerManager{
		servers: RootServers,
		ipToServer: make(map[string]*NameServer),
		nameToServer: make(map[string]*NameServer),
	}

	for rootServer := range RootServers {
		server := &RootServers[rootServer]
		manager.ipToServer[server.IPv4] = server
		manager.ipToServer[server.IPv6] = server
		manager.nameToServer[server.Name] = server
	}

	return manager
}

func (m *RootServerManager) LookUpByIP(ip string) (*NameServer, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	server, found := m.ipToServer[ip]
	return server, found
}

func (m *RootServerManager) GetFirstRoot(ip string) *NameServer {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.servers) > 0 {
		return &m.servers[0]
	}
	return nil 
}