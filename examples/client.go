package main

import (
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

// GetTransactionStatusArgs represents parameters for the getTransactionStatus method
type GetTransactionStatusArgs struct {
	ID string `json:"id"`
}

// GetTransactionStatusResult represents the result of the getTransactionStatus method
type GetTransactionStatusResult struct {
	Exists      bool             `json:"exists"`
	Transaction *TransactionInfo `json:"transaction,omitempty"`
}

// TransactionInfo represents transaction info returned by the API
type TransactionInfo struct {
	ID        string    `json:"id"`
	Data      []byte    `json:"data"`
	Priority  int       `json:"priority"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	// Create a new random source with current time as seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Server URL
	serverURL := "http://localhost:8080"

	// Connect to the server
	client, err := rpc.Dial(serverURL)
	if err != nil {
		log.Fatalf("Failed to connect to the server: %v", err)
	}
	defer client.Close()

	// Submit 10 transactions with random priorities
	var txIDs []string
	for i := 0; i < 10; i++ {
		// Generate random data and priority
		data := fmt.Sprintf("Transaction data %d", i)
		priority := r.Intn(100)

		// Submit transaction
		txID, err := submitTransaction(client, data, priority)
		if err != nil {
			log.Fatalf("Failed to submit transaction: %v", err)
		}

		txIDs = append(txIDs, txID)
		log.Printf("Submitted transaction %d with ID: %s, Priority: %d", i, txID, priority)

		// Sleep briefly to avoid overwhelming the server
		time.Sleep(50 * time.Millisecond)
	}

	// Wait for transactions to be processed
	log.Println("Waiting for transactions to be processed...")
	time.Sleep(1 * time.Second)

	// Check status of each transaction
	for i, txID := range txIDs {
		exists, err := checkTransactionStatus(client, txID)
		if err != nil {
			log.Printf("Failed to check transaction status: %v", err)
			continue
		}

		if exists {
			log.Printf("Transaction %d (ID: %s) is still in the mempool", i, txID)
		} else {
			log.Printf("Transaction %d (ID: %s) has been processed into a block", i, txID)
		}
	}
}

// submitTransaction submits a transaction to the server
func submitTransaction(client *rpc.Client, data string, priority int) (string, error) {
	args := SubmitTransactionArgs{
		Data:     data,
		Priority: priority,
	}

	var result SubmitTransactionResult
	err := client.Call(&result, "flash_submitTransaction", args)
	if err != nil {
		return "", fmt.Errorf("RPC error: %v", err)
	}

	return result.TransactionID, nil
}

// checkTransactionStatus checks the status of a transaction
func checkTransactionStatus(client *rpc.Client, txID string) (bool, error) {
	args := GetTransactionStatusArgs{
		ID: txID,
	}

	var result GetTransactionStatusResult
	err := client.Call(&result, "flash_getTransactionStatus", args)
	if err != nil {
		return false, fmt.Errorf("RPC error: %v", err)
	}

	return result.Exists, nil
}
