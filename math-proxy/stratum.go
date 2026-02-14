package main

// StratumRequest represents a JSON-RPC request from a mining client
type StratumRequest struct {
	ID     interface{}   `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

// StratumResponse represents a JSON-RPC response to a mining client
type StratumResponse struct {
	ID     interface{}   `json:"id"`
	Result interface{}   `json:"result"`
	Error  []interface{} `json:"error"`
}

// StratumJob represents a mining job
type StratumJob struct {
	JobID          string
	PrevHash       string
	CoinBase1      string
	CoinBase2      string
	MerkleBranches []string
	Version        string
	NBits          string
	NTime          string
	CleanJobs      bool
}
