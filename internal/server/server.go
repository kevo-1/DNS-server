package server

import (
	"DNS-server/internal/transport"
	"DNS-server/models"
	"DNS-server/pkg/resolver"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type Server struct {
	config   *Config
	handler  *Handler
	udp      *transport.UDPTransport
	tcp      *transport.TCPTransport
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	resolver *resolver.Resolver
}

func NewServer(config *Config) (*Server, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	cacheConfig := &models.CacheConfig{
		MaxEntries:      config.CacheMaxEntries,
		DefaultTTL:      config.CacheTTL,
		CleanupInterval: config.CacheCleanupInterval,
		EnableStats:     true,
	}

	res := resolver.NewResolver(cacheConfig)

	handler := NewHandler(config)

	ctx, cancel := context.WithCancel(context.Background())

	server := &Server{
		config:   config,
		handler:  handler,
		ctx:      ctx,
		cancel:   cancel,
		resolver: res,
	}

	return server, nil
}

func (s *Server) Start() error {
	log.Println("Starting DNS server...")

	if s.config.EnableUDP {
		udp := transport.NewUDPTransport(s.config.GetUDPAddress(), s.handler.HandleRequest)
		s.udp = udp

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if err := s.udp.Start(s.ctx); err != nil {
				log.Printf("UDP transport error: %v", err)
			}
		}()

		log.Printf("UDP server listening on %s", s.config.GetUDPAddress())
	}

	if s.config.EnableTCP {
		tcp := transport.NewTCPTransport(s.config.GetTCPAddress(), s.handler.HandleRequest)
		s.tcp = tcp

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if err := s.tcp.Start(s.ctx); err != nil {
				log.Printf("TCP transport error: %v", err)
			}
		}()

		log.Printf("TCP server listening on %s", s.config.GetTCPAddress())
	}

	log.Println("DNS server started successfully")
	return nil
}

func (s *Server) Stop() error {
	log.Println("Stopping DNS server...")

	s.cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("DNS server stopped gracefully")
	case <-shutdownCtx.Done():
		log.Println("DNS server shutdown timed out")
		return fmt.Errorf("shutdown timeout")
	}

	s.resolver.Close()

	return nil
}

func (s *Server) Wait() {
	s.wg.Wait()
}

func (s *Server) GetStats() models.CacheStatistics {
	return s.handler.GetStats()
}

func (s *Server) GetConfig() *Config {
	return s.config
}

func (s *Server) IsRunning() bool {
	select {
	case <-s.ctx.Done():
		return false
	default:
		return true
	}
}
