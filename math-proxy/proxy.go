package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

// Proxy represents the math proxy server
type Proxy struct {
	config   *Config
	listener net.Listener
	clients  map[net.Conn]bool
	mu       sync.Mutex
	wg       sync.WaitGroup
	quit     chan struct{}
}

// NewProxy creates a new proxy instance
func NewProxy(config *Config) *Proxy {
	return &Proxy{
		config:  config,
		clients: make(map[net.Conn]bool),
		quit:    make(chan struct{}),
	}
}

// Start starts the proxy server
func (p *Proxy) Start() error {
	listener, err := net.Listen("tcp", p.config.ListenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	p.listener = listener

	log.Printf("Proxy listening on %s", p.config.ListenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-p.quit:
				return nil
			default:
				log.Printf("Error accepting connection: %v", err)
				continue
			}
		}

		p.mu.Lock()
		p.clients[conn] = true
		p.mu.Unlock()

		p.wg.Add(1)
		go p.handleClient(conn)
	}
}

// Stop stops the proxy server
func (p *Proxy) Stop() {
	close(p.quit)
	if p.listener != nil {
		p.listener.Close()
	}

	p.mu.Lock()
	for conn := range p.clients {
		conn.Close()
	}
	p.mu.Unlock()

	p.wg.Wait()
}

// handleClient handles a single client connection
func (p *Proxy) handleClient(conn net.Conn) {
	defer p.wg.Done()
	defer conn.Close()

	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.clients, conn)

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Client connected: %s", clientAddr)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if p.config.Debug {
			log.Printf("Received from %s: %s", clientAddr, line)
		}

		// Parse Stratum JSON-RPC request
		var request StratumRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			continue
		}

		// Handle the request
		response := p.handleStratumRequest(&request, clientAddr)

		// Send response
		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			continue
		}

		if p.config.Debug {
			log.Printf("Sending to %s: %s", clientAddr, string(responseBytes))
		}

		_, err = conn.Write(append(responseBytes, '\n'))
		if err != nil {
			log.Printf("Error sending response: %v", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Connection error with %s: %v", clientAddr, err)
	}

	log.Printf("Client disconnected: %s", clientAddr)
}

// handleStratumRequest handles a Stratum protocol request
func (p *Proxy) handleStratumRequest(req *StratumRequest, clientAddr string) *StratumResponse {
	response := &StratumResponse{
		ID:     req.ID,
		Result: nil,
		Error:  nil,
	}

	switch req.Method {
	case "mining.subscribe":
		// Handle mining.subscribe
		response.Result = []interface{}{
			[]interface{}{
				[]interface{}{"mining.set_difficulty", "randomx-proxy-001"},
				[]interface{}{"mining.notify", "randomx-proxy-001"},
			},
			"randomx-session-" + clientAddr,
			4,
		}
		log.Printf("Client %s subscribed", clientAddr)

	case "mining.authorize":
		// Handle mining.authorize
		response.Result = true
		if len(req.Params) > 0 {
			if username, ok := req.Params[0].(string); ok {
				log.Printf("Client %s authorized as %s", clientAddr, username)
			}
		}

	case "mining.submit":
		// Handle mining.submit - this is where we convert scrypt to randomx
		if len(req.Params) >= 5 {
			worker := getStringParam(req.Params, 0)
			jobID := getStringParam(req.Params, 1)
			extraNonce2 := getStringParam(req.Params, 2)
			nTime := getStringParam(req.Params, 3)
			nonce := getStringParam(req.Params, 4)

			log.Printf("Mining submit from %s: worker=%s, job=%s, extranonce2=%s, ntime=%s, nonce=%s",
				clientAddr, worker, jobID, extraNonce2, nTime, nonce)

			// Here we would normally validate the scrypt share
			// and convert it to a randomx share, then submit to upstream pool
			
			// For now, we'll perform RandomX hash calculation
			// In a real implementation, you'd extract the block header and calculate RandomX hash
			hash := CalculateRandomXHash([]byte(jobID + extraNonce2 + nTime + nonce))
			
			if p.config.Debug {
				log.Printf("RandomX hash for submit: %x", hash)
			}

			// Accept the share (in production, you'd validate against target difficulty)
			response.Result = true
		} else {
			response.Error = []interface{}{20, "Invalid submit parameters", nil}
		}

	case "mining.extranonce.subscribe":
		// Handle extranonce subscription
		response.Result = true

	default:
		log.Printf("Unknown method: %s", req.Method)
		response.Error = []interface{}{20, "Unknown method", nil}
	}

	return response
}

// Helper function to safely get string parameter from interface slice
func getStringParam(params []interface{}, index int) string {
	if index < len(params) {
		if str, ok := params[index].(string); ok {
			return str
		}
	}
	return ""
}
