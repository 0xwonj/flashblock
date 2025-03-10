package eth

import (
	"encoding/hex"
	"fmt"
	"strings"

	"flashblock/internal/eth"
	"flashblock/internal/mempool"
)

// TransactionHook is a function called when a transaction is processed
type TransactionHook = mempool.TransactionHook

// API represents the Ethereum compatible JSON-RPC API
type API struct {
	mempool *mempool.Mempool
}

// SendRawTransactionArgs represents the arguments for eth_sendRawTransaction
type SendRawTransactionArgs struct {
	RawTx string
}

// SendRawTransactionResult represents the result of eth_sendRawTransaction
type SendRawTransactionResult struct {
	TransactionHash string
}

// NewAPI creates a new Ethereum API instance
func NewAPI(mempool *mempool.Mempool, hooks []TransactionHook) *API {
	return &API{
		mempool: mempool,
	}
}

// SendRawTransaction implements the eth_sendRawTransaction RPC method
func (api *API) SendRawTransaction(rawTx string) (string, error) {
	// Remove "0x" prefix if present
	rawTx = strings.TrimPrefix(rawTx, "0x")

	// Parse the raw transaction
	tx, err := eth.ParseRawTransaction(rawTx)
	if err != nil {
		return "", fmt.Errorf("invalid raw transaction: %w", err)
	}

	// Add transaction to mempool
	api.mempool.AddTransaction(tx)

	// Return the transaction hash (ID)
	return "0x" + tx.ID, nil
}

// GetTransactionByHash implements the eth_getTransactionByHash RPC method
func (api *API) GetTransactionByHash(hash string) (map[string]any, error) {
	// Remove "0x" prefix if present
	hash = strings.TrimPrefix(hash, "0x")

	// Get transaction from mempool
	tx, exists := api.mempool.GetTransaction(hash)
	if !exists {
		return nil, nil // Return null if transaction not found
	}

	// Convert to Ethereum format
	result := map[string]any{
		"hash":             "0x" + tx.ID,
		"from":             tx.From,
		"to":               nil,
		"value":            "0x0",
		"gas":              "0x0",
		"gasPrice":         "0x0",
		"nonce":            "0x0",
		"input":            "0x" + hex.EncodeToString(tx.Data),
		"blockHash":        nil,
		"blockNumber":      nil,
		"transactionIndex": nil,
	}

	// Add Ethereum-specific fields if available
	if tx.To != "" {
		result["to"] = tx.To
	}
	if tx.Value != nil && tx.Value.BitLen() > 0 {
		result["value"] = "0x" + tx.Value.Text(16)
	}
	if tx.GasPrice != nil && tx.GasPrice.BitLen() > 0 {
		result["gasPrice"] = "0x" + tx.GasPrice.Text(16)
	}
	if tx.GasLimit > 0 {
		result["gas"] = fmt.Sprintf("0x%x", tx.GasLimit)
	}
	if tx.Nonce > 0 {
		result["nonce"] = fmt.Sprintf("0x%x", tx.Nonce)
	}

	return result, nil
}

// GetTransactionReceipt implements the eth_getTransactionReceipt RPC method
func (api *API) GetTransactionReceipt(hash string) (map[string]any, error) {
	// This is a simplified version that will always return null
	// In a real implementation, you would check if the transaction is in a processed block
	return nil, nil
}
