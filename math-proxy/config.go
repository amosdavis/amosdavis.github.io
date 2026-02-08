package main

import (
	"encoding/json"
	"os"
)

// Config holds the configuration for the proxy
type Config struct {
	ListenAddr   string `json:"listen_addr"`
	UpstreamAddr string `json:"upstream_addr"`
	Debug        bool   `json:"debug"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		ListenAddr:   "0.0.0.0:3333",
		UpstreamAddr: "localhost:3334",
		Debug:        false,
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := DefaultConfig()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
