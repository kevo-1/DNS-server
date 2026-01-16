package main

import (
	"DNS-server/internal/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("DNS Server starting...")

	config := server.DefaultConfig()

	srv, err := server.NewServer(config)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal: %v", sig)

	log.Println("Shutting down server...")
	if err := srv.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
		os.Exit(1)
	}

	stats := srv.GetStats()
	fmt.Printf("\nServer Statistics:\n")
	fmt.Printf("  Cache Hits: %d\n", stats.Hits)
	fmt.Printf("  Cache Misses: %d\n", stats.Misses)
	fmt.Printf("  Cache Hit Rate: %.2f%%\n", stats.HitRate()*100)
	fmt.Printf("  Cache Evictions: %d\n", stats.Evictions)
	fmt.Printf("  Total Entries: %d/%d\n", stats.TotalEntries, stats.TotalCapacity)

	log.Println("Server stopped successfully")
}
