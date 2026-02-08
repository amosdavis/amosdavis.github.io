package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// TestClient is a simple Stratum client for testing
func main() {
	// Connect to the proxy
	conn, err := net.Dial("tcp", "localhost:3333")
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to proxy")

	// Send mining.subscribe
	subscribe := StratumRequest{
		ID:     1,
		Method: "mining.subscribe",
		Params: []interface{}{"test-miner/1.0"},
	}
	
	if err := sendRequest(conn, &subscribe); err != nil {
		fmt.Printf("Failed to send subscribe: %v\n", err)
		return
	}
	
	response, err := readResponse(conn)
	if err != nil {
		fmt.Printf("Failed to read subscribe response: %v\n", err)
		return
	}
	fmt.Printf("Subscribe response: %+v\n", response)

	// Send mining.authorize
	authorize := StratumRequest{
		ID:     2,
		Method: "mining.authorize",
		Params: []interface{}{"test.worker", "password"},
	}
	
	if err := sendRequest(conn, &authorize); err != nil {
		fmt.Printf("Failed to send authorize: %v\n", err)
		return
	}
	
	response, err = readResponse(conn)
	if err != nil {
		fmt.Printf("Failed to read authorize response: %v\n", err)
		return
	}
	fmt.Printf("Authorize response: %+v\n", response)

	// Send mining.submit
	submit := StratumRequest{
		ID:     3,
		Method: "mining.submit",
		Params: []interface{}{
			"test.worker",           // worker name
			"job123",                // job id
			"00000000",              // extranonce2
			"507c0000",              // ntime
			"00000000",              // nonce
		},
	}
	
	if err := sendRequest(conn, &submit); err != nil {
		fmt.Printf("Failed to send submit: %v\n", err)
		return
	}
	
	response, err = readResponse(conn)
	if err != nil {
		fmt.Printf("Failed to read submit response: %v\n", err)
		return
	}
	fmt.Printf("Submit response: %+v\n", response)

	fmt.Println("\nAll tests passed!")
}

func sendRequest(conn net.Conn, req *StratumRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	
	_, err = conn.Write(append(data, '\n'))
	return err
}

func readResponse(conn net.Conn) (*StratumResponse, error) {
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	
	var response StratumResponse
	err = json.Unmarshal([]byte(line), &response)
	if err != nil {
		return nil, err
	}
	
	return &response, nil
}
