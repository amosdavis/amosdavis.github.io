package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Parse command-line flags
	configFile := flag.String("config", "config.json", "Path to configuration file")
	listenAddr := flag.String("listen", "0.0.0.0:3333", "Address to listen on")
	upstreamAddr := flag.String("upstream", "localhost:3334", "Upstream pool address")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Load configuration
	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Printf("Warning: Could not load config file: %v. Using defaults.", err)
		config = DefaultConfig()
	}

	// Override with command-line flags if provided
	if *listenAddr != "0.0.0.0:3333" {
		config.ListenAddr = *listenAddr
	}
	if *upstreamAddr != "localhost:3334" {
		config.UpstreamAddr = *upstreamAddr
	}
	config.Debug = *debug

	log.Printf("Math Proxy starting...")
	log.Printf("Listen address: %s", config.ListenAddr)
	log.Printf("Upstream pool: %s", config.UpstreamAddr)
	log.Printf("Debug mode: %v", config.Debug)

	// Create and start the proxy server
	proxy := NewProxy(config)

	// Start the proxy server in a goroutine
	go func() {
		if err := proxy.Start(); err != nil {
			log.Fatalf("Failed to start proxy: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down proxy...")
	proxy.Stop()
	log.Println("Proxy stopped")
}
