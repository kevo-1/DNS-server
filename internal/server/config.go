package server

import (
	"fmt"
	"time"
)

type Config struct {
	// Server settings
	UDPPort int
	TCPPort int
	Host    string

	// Timeouts
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// Limits
	MaxConnections int
	MaxUDPSize     int

	// Features
	EnableUDP       bool
	EnableTCP       bool
	EnableRecursion bool
	EnableCaching   bool

	// Cache settings
	CacheMaxEntries     int
	CacheTTL            time.Duration
	CacheCleanupInterval time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		// Server settings
		UDPPort: 53,
		TCPPort: 53,
		Host:    "0.0.0.0",

		// Timeouts
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,

		// Limits
		MaxConnections: 100,
		MaxUDPSize:     512,

		// Features
		EnableUDP:       true,
		EnableTCP:       true,
		EnableRecursion: true,
		EnableCaching:   true,

		// Cache settings
		CacheMaxEntries:      1000,
		CacheTTL:             5 * time.Minute,
		CacheCleanupInterval: 1 * time.Minute,
	}
}

func (c *Config) Validate() error {
	if c.UDPPort < 0 || c.UDPPort > 65535 {
		return &ConfigError{"invalid UDP port"}
	}

	if c.TCPPort < 0 || c.TCPPort > 65535 {
		return &ConfigError{"invalid TCP port"}
	}

	if !c.EnableUDP && !c.EnableTCP {
		return &ConfigError{"at least one transport (UDP or TCP) must be enabled"}
	}

	if c.MaxUDPSize < 512 {
		return &ConfigError{"max UDP size must be at least 512 bytes"}
	}

	if c.MaxConnections < 1 {
		return &ConfigError{"max connections must be at least 1"}
	}

	return nil
}

type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return "config error: " + e.Message
}

func (c *Config) GetUDPAddress() string {
	return formatAddress(c.Host, c.UDPPort)
}

func (c *Config) GetTCPAddress() string {
	return formatAddress(c.Host, c.TCPPort)
}

func formatAddress(host string, port int) string {
	if host == "" {
		host = "0.0.0.0"
	}
	return fmt.Sprintf("%s:%d", host, port)
}
