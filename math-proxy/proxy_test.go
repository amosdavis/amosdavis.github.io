package main

import (
	"testing"
)

func TestCalculateRandomXHash(t *testing.T) {
	// Test basic hash calculation
	data := []byte("test data")
	hash := CalculateRandomXHash(data)
	
	if len(hash) != 32 {
		t.Errorf("Expected hash length 32, got %d", len(hash))
	}
	
	// Hash should be deterministic
	hash2 := CalculateRandomXHash(data)
	if string(hash) != string(hash2) {
		t.Error("Hash calculation is not deterministic")
	}
	
	// Different data should produce different hash
	differentData := []byte("different data")
	differentHash := CalculateRandomXHash(differentData)
	if string(hash) == string(differentHash) {
		t.Error("Different data produced same hash")
	}
}

func TestValidateRandomXHash(t *testing.T) {
	tests := []struct {
		name     string
		hash     []byte
		target   []byte
		expected bool
	}{
		{
			name:     "hash below target",
			hash:     make([]byte, 32),
			target:   []byte{0xFF, 0xFF, 0xFF, 0xFF, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			expected: true,
		},
		{
			name:     "invalid hash length",
			hash:     make([]byte, 16),
			target:   make([]byte, 32),
			expected: false,
		},
		{
			name:     "invalid target length",
			hash:     make([]byte, 32),
			target:   make([]byte, 16),
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateRandomXHash(tt.hash, tt.target)
			if result != tt.expected {
				t.Errorf("ValidateRandomXHash() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDifficultyToTarget(t *testing.T) {
	target := DifficultyToTarget(1.0)
	
	if len(target) != 32 {
		t.Errorf("Expected target length 32, got %d", len(target))
	}
	
	// Higher difficulty should produce lower target
	target1 := DifficultyToTarget(1.0)
	target2 := DifficultyToTarget(2.0)
	
	// At least one byte should be smaller for higher difficulty
	different := false
	for i := 0; i < 32; i++ {
		if target2[i] < target1[i] {
			different = true
			break
		}
	}
	
	if !different {
		t.Error("Higher difficulty should produce lower target")
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.ListenAddr == "" {
		t.Error("ListenAddr should not be empty")
	}
	
	if config.UpstreamAddr == "" {
		t.Error("UpstreamAddr should not be empty")
	}
}

func TestStratumRequest(t *testing.T) {
	// Test that we can create stratum structures
	req := &StratumRequest{
		ID:     1,
		Method: "mining.subscribe",
		Params: []interface{}{},
	}
	
	if req.Method != "mining.subscribe" {
		t.Error("Failed to create StratumRequest")
	}
}

func TestStratumResponse(t *testing.T) {
	// Test that we can create stratum response structures
	resp := &StratumResponse{
		ID:     1,
		Result: true,
		Error:  nil,
	}
	
	if resp.Result != true {
		t.Error("Failed to create StratumResponse")
	}
}
