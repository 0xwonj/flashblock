package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"gopkg.in/yaml.v2"
)

// WorkloadConfig represents the configuration for the workload
type WorkloadConfig struct {
	NumClients        int    `yaml:"num_clients"`
	RequestsPerSecond int    `yaml:"requests_per_second"`
	DurationSeconds   int    `yaml:"duration_seconds"`
	ServerURL         string `yaml:"server_url"`
}

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
	// Parse command-line flags
	configFile := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting workload with %d clients, %d requests/sec per client, for %d seconds",
		config.NumClients, config.RequestsPerSecond, config.DurationSeconds)

	// Create a WaitGroup to wait for all clients to complete
	var wg sync.WaitGroup

	// Start the specified number of clients
	for i := range config.NumClients {
		wg.Add(1)
		go runClient(i, config, &wg)
	}

	// Wait for all clients to complete
	wg.Wait()
	log.Println("Workload completed")
}

// loadConfig loads the workload configuration from a YAML file
func loadConfig(filePath string) (*WorkloadConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config WorkloadConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Validate configuration
	if config.NumClients <= 0 {
		return nil, fmt.Errorf("num_clients must be greater than 0")
	}
	if config.RequestsPerSecond <= 0 {
		return nil, fmt.Errorf("requests_per_second must be greater than 0")
	}
	if config.DurationSeconds <= 0 {
		return nil, fmt.Errorf("duration_seconds must be greater than 0")
	}
	if config.ServerURL == "" {
		return nil, fmt.Errorf("server_url cannot be empty")
	}

	return &config, nil
}

// runClient runs a single client that generates the specified workload
func runClient(clientID int, config *WorkloadConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	// Create a new random source with current time and client ID as seed
	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(clientID)))

	// Connect to the server
	client, err := rpc.Dial(config.ServerURL)
	if err != nil {
		log.Printf("Client %d: Failed to connect to the server: %v", clientID, err)
		return
	}
	defer client.Close()

	log.Printf("Client %d: Connected to server %s", clientID, config.ServerURL)

	// Calculate interval between requests to achieve the desired rate
	interval := time.Second / time.Duration(config.RequestsPerSecond)

	// Create a timer to control the request rate
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Create a timer for the overall duration
	timeout := time.After(time.Duration(config.DurationSeconds) * time.Second)

	// Track transactions for status checking
	var txIDs []string
	var txIDsMutex sync.Mutex

	// Run the workload
	txCounter := 0
	for {
		select {
		case <-timeout:
			// Duration complete
			log.Printf("Client %d: Completed workload (%d transactions sent)", clientID, txCounter)

			// Check status of transactions (sample up to 10)
			checkTransactionStatuses(client, txIDs, clientID)
			return

		case <-ticker.C:
			// Time to send another transaction
			data := fmt.Sprintf("Client %d transaction %d", clientID, txCounter)
			priority := r.Intn(100)

			// Submit transaction
			txID, err := submitTransaction(client, data, priority)
			if err != nil {
				log.Printf("Client %d: Failed to submit transaction: %v", clientID, err)
				continue
			}

			// Store the transaction ID
			txIDsMutex.Lock()
			txIDs = append(txIDs, txID)
			txIDsMutex.Unlock()

			if txCounter%100 == 0 {
				log.Printf("Client %d: Submitted %d transactions", clientID, txCounter)
			}

			txCounter++
		}
	}
}

// checkTransactionStatuses checks the status of a sampling of transactions
func checkTransactionStatuses(client *rpc.Client, txIDs []string, clientID int) {
	// Sample up to 10 transactions to check
	sampleSize := min(10, len(txIDs))

	// Check status of sampled transactions
	for i := range sampleSize {
		// Get a transaction from evenly distributed positions in the array
		idx := i * len(txIDs) / sampleSize
		txID := txIDs[idx]

		exists, err := checkTransactionStatus(client, txID)
		if err != nil {
			log.Printf("Client %d: Failed to check transaction status: %v", clientID, err)
			continue
		}

		if exists {
			log.Printf("Client %d: Transaction (ID: %s) is still in the mempool", clientID, txID)
		} else {
			log.Printf("Client %d: Transaction (ID: %s) has been processed", clientID, txID)
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
