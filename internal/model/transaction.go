package model

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"time"
)

// Transaction represents a single transaction in the system with Ethereum-compatible fields
type Transaction struct {
	ID        string    `json:"id"`
	Data      []byte    `json:"data"`     // Transaction payload data
	Priority  int       `json:"priority"` // Legacy priority (will be replaced by gas price)
	Timestamp time.Time `json:"timestamp"`

	// Ethereum transaction fields
	From     string   `json:"from"`      // Sender address
	To       string   `json:"to"`        // Recipient address
	Value    *big.Int `json:"value"`     // Transaction value in wei
	GasPrice *big.Int `json:"gas_price"` // Gas price in wei
	GasLimit uint64   `json:"gas_limit"` // Gas limit
	Nonce    uint64   `json:"nonce"`     // Transaction nonce
	RawData  string   `json:"raw_data"`  // Original raw transaction data
}

// NewTransaction creates a new transaction with the given data and priority
func NewTransaction(data []byte, priority int) *Transaction {
	// Generate a transaction ID by hashing the data and timestamp
	hash := sha256.Sum256(append(data, []byte(time.Now().String())...))

	return &Transaction{
		ID:        hex.EncodeToString(hash[:]),
		Data:      data,
		Priority:  priority,
		Timestamp: time.Now(),
		Value:     new(big.Int),
		GasPrice:  new(big.Int),
	}
}

// NewEthereumTransaction creates a new transaction from Ethereum transaction data
func NewEthereumTransaction(
	from string,
	to string,
	value *big.Int,
	gasPrice *big.Int,
	gasLimit uint64,
	nonce uint64,
	data []byte,
	rawData string,
) *Transaction {
	// Generate a transaction ID by hashing the data and timestamp
	hashInput := append(data, []byte(time.Now().String())...)
	hashInput = append(hashInput, []byte(from)...)
	hashInput = append(hashInput, []byte(to)...)
	hash := sha256.Sum256(hashInput)

	// Set priority based on gas price
	priority := 0
	if gasPrice != nil && gasPrice.BitLen() > 0 {
		// Convert gas price to a priority value
		// Higher gas price = higher priority
		// This is a simplified conversion, might need adjustment
		priority = int(new(big.Int).Div(gasPrice, big.NewInt(1000000000)).Int64())
	}

	return &Transaction{
		ID:        hex.EncodeToString(hash[:]),
		Data:      data,
		Priority:  priority,
		Timestamp: time.Now(),
		From:      from,
		To:        to,
		Value:     value,
		GasPrice:  gasPrice,
		GasLimit:  gasLimit,
		Nonce:     nonce,
		RawData:   rawData,
	}
}
