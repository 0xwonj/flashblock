package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

// SubmitTransactionArgs represents parameters for the submitTransaction method
type SubmitTransactionArgs struct {
	Data     string `json:"data"`
	Priority int    `json:"priority"`
}

// SubmitTransactionResult represents the result of the submitTransaction method
type SubmitTransactionResult struct {
	TransactionID string `json:"transaction_id"`
	Added         bool   `json:"added"`
}

// StatusResult represents the system status
type StatusResult struct {
	Status          string `json:"status"`
	Uptime          string `json:"uptime"`
	Version         string `json:"version"`
	MempoolSize     int    `json:"mempool_size"`
	BlocksProcessed int    `json:"blocks_processed"`
}

// TransactionInfo represents transaction info returned by the API
type TransactionInfo struct {
	ID        string    `json:"id"`
	Data      []byte    `json:"data"`
	Priority  int       `json:"priority"`
	Timestamp time.Time `json:"timestamp"`
}

// SubscribeRequest represents a JSON-RPC request for subscription
type SubscribeRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []any         `json:"params"`
}

func main() {
	// Create a new random source with current time as seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Server URL for WebSocket connection
	wsURL := "ws://localhost:8080/ws"

	// Connect to the server via WebSocket
	client, err := rpc.DialWebsocket(context.Background(), wsURL, "*")
	if err != nil {
		log.Fatalf("Failed to connect to the server: %v", err)
	}
	defer client.Close()

	log.Printf("Connected to %s via WebSocket", wsURL)

	// First, get system status
	var status StatusResult
	err = client.Call(&status, "flash_getStatus")
	if err != nil {
		log.Fatalf("Failed to get status: %v", err)
	}
	log.Printf("System status: %+v", status)

	// Submit 5 transactions with random priorities
	var lastTxID string
	for i := 0; i < 5; i++ {
		// Generate random data and priority
		data := fmt.Sprintf("WebSocket Transaction %d", i)
		priority := r.Intn(100)

		args := SubmitTransactionArgs{
			Data:     data,
			Priority: priority,
		}

		var result SubmitTransactionResult
		err := client.Call(&result, "flash_submitTransaction", args)
		if err != nil {
			log.Printf("Failed to submit transaction: %v", err)
			continue
		}

		lastTxID = result.TransactionID
		log.Printf("Submitted transaction %d with ID: %s, Priority: %d", i, result.TransactionID, priority)

		// Sleep briefly to avoid overwhelming the server
		time.Sleep(50 * time.Millisecond)
	}

	// Wait for processing
	log.Println("Waiting for transactions to be processed...")
	time.Sleep(1 * time.Second)

	// Check mempool
	var mempoolResult struct {
		Transactions []TransactionInfo `json:"transactions"`
		Count        int               `json:"count"`
	}
	err = client.Call(&mempoolResult, "flash_getMempool")
	if err != nil {
		log.Fatalf("Failed to get mempool: %v", err)
	}
	log.Printf("Mempool has %d transactions", mempoolResult.Count)

	// Check blocks
	var blocksResult struct {
		Blocks []json.RawMessage `json:"blocks"`
		Count  int               `json:"count"`
	}
	err = client.Call(&blocksResult, "flash_getBlocks")
	if err != nil {
		log.Fatalf("Failed to get blocks: %v", err)
	}
	log.Printf("System has %d blocks", blocksResult.Count)

	// Demonstrate batch requests
	doBatchRequests(client, lastTxID)

	log.Println("WebSocket client example completed")
}

// doBatchRequests demonstrates batch requests
func doBatchRequests(client *rpc.Client, txID string) {
	log.Println("Performing batch RPC requests...")

	var statusResult StatusResult
	var txStatusResult struct {
		Exists      bool             `json:"exists"`
		Transaction *TransactionInfo `json:"transaction,omitempty"`
	}

	batch := []rpc.BatchElem{
		{
			Method: "flash_getStatus",
			Result: &statusResult,
		},
		{
			Method: "flash_getTransactionStatus",
			Args:   []any{map[string]string{"id": txID}},
			Result: &txStatusResult,
		},
	}

	err := client.BatchCall(batch)
	if err != nil {
		log.Fatalf("Batch request failed: %v", err)
	}

	for i, elem := range batch {
		if elem.Error != nil {
			log.Printf("Request %d failed: %v", i, elem.Error)
			continue
		}
	}

	log.Printf("Batch result - Status: %s, Transaction exists: %v",
		statusResult.Status, txStatusResult.Exists)
}
