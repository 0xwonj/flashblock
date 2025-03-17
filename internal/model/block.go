package model

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Block represents a collection of transactions
type Block struct {
	ID           string         `json:"id"`
	Transactions []*Transaction `json:"transactions"`
	Timestamp    time.Time      `json:"timestamp"`
	PrevBlockID  string         `json:"prev_block_id"`
	TDXQuote     []byte         `json:"tdx_quote,omitempty"`
}

// NewBlock creates a new block with the given transactions and previous block ID
func NewBlock(transactions []*Transaction, prevBlockID string) *Block {
	timestamp := time.Now()

	// Create a new block
	block := &Block{
		Transactions: transactions,
		Timestamp:    timestamp,
		PrevBlockID:  prevBlockID,
	}

	// Generate block ID by hashing its contents
	block.generateID()

	return block
}

// generateID creates a unique ID for the block based on its contents
func (b *Block) generateID() {
	// Concatenate transaction IDs, timestamp, and previous block ID
	var data []byte
	for _, tx := range b.Transactions {
		data = append(data, []byte(tx.ID)...)
	}
	data = append(data, []byte(b.Timestamp.String())...)
	data = append(data, []byte(b.PrevBlockID)...)

	// Hash the data to generate block ID
	hash := sha256.Sum256(data)
	b.ID = hex.EncodeToString(hash[:])
}
