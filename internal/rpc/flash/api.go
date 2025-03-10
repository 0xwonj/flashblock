package flash

import (
	"encoding/base64"
	"errors"
	"time"

	"flashblock/internal/mempool"
	"flashblock/internal/model"
	"flashblock/internal/processor"
)

// TransactionHook is a function called when a transaction is processed
type TransactionHook = mempool.TransactionHook

// API defines the Flash RPC methods
type API struct {
	mempool   *mempool.Mempool
	processor *processor.BlockProcessor
	startTime time.Time
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
	Exists      bool               `json:"exists"`
	Transaction *model.Transaction `json:"transaction,omitempty"`
}

// GetBlocksResult represents a list of blocks
type GetBlocksResult struct {
	Blocks []*model.Block `json:"blocks"`
	Count  int            `json:"count"`
}

// GetMempoolResult represents the current mempool state
type GetMempoolResult struct {
	Transactions []*model.Transaction `json:"transactions"`
	Count        int                  `json:"count"`
}

// StatusResult represents the system status
type StatusResult struct {
	Status          string `json:"status"`
	Uptime          string `json:"uptime"`
	Version         string `json:"version"`
	MempoolSize     int    `json:"mempool_size"`
	BlocksProcessed int    `json:"blocks_processed"`
}

// NewAPI creates a new Flash API instance
func NewAPI(mempool *mempool.Mempool, processor *processor.BlockProcessor, hooks []TransactionHook) *API {
	return &API{
		mempool:   mempool,
		processor: processor,
		startTime: time.Now(),
	}
}

// SubmitTransaction handles transaction submission
func (api *API) SubmitTransaction(args SubmitTransactionArgs) (*SubmitTransactionResult, error) {
	// Validate parameters
	if args.Data == "" {
		return nil, errors.New("data cannot be empty")
	}

	// Decode base64 data if necessary
	var data []byte
	var err error

	// Try to decode as base64, otherwise use as raw bytes
	data, err = base64.StdEncoding.DecodeString(args.Data)
	if err != nil {
		// If not base64, use the original string as bytes
		data = []byte(args.Data)
	}

	// Create transaction
	tx := model.NewTransaction(data, args.Priority)

	// Add to mempool
	added := api.mempool.AddTransaction(tx)

	// Return result
	return &SubmitTransactionResult{
		TransactionID: tx.ID,
		Added:         added,
	}, nil
}

// GetTransactionStatus checks the status of a transaction
func (api *API) GetTransactionStatus(args GetTransactionStatusArgs) (*GetTransactionStatusResult, error) {
	// Validate parameters
	if args.ID == "" {
		return nil, errors.New("transaction ID cannot be empty")
	}

	// Get transaction from mempool
	tx, exists := api.mempool.GetTransaction(args.ID)

	// Return result
	return &GetTransactionStatusResult{
		Exists:      exists,
		Transaction: tx,
	}, nil
}

// GetBlocks returns all processed blocks
func (api *API) GetBlocks() (*GetBlocksResult, error) {
	if api.processor == nil {
		return nil, errors.New("block processor not available")
	}

	blocks := api.processor.GetProcessedBlocks()
	return &GetBlocksResult{
		Blocks: blocks,
		Count:  len(blocks),
	}, nil
}

// GetMempool returns all transactions in the mempool
func (api *API) GetMempool() (*GetMempoolResult, error) {
	transactions := api.mempool.GetAllTransactions()
	return &GetMempoolResult{
		Transactions: transactions,
		Count:        len(transactions),
	}, nil
}

// GetStatus returns system status
func (api *API) GetStatus() (*StatusResult, error) {
	var blocksProcessed int
	if api.processor != nil {
		blocksProcessed = len(api.processor.GetProcessedBlocks())
	}

	return &StatusResult{
		Status:          "running",
		Uptime:          time.Since(api.startTime).String(),
		Version:         "1.0.0",
		MempoolSize:     api.mempool.Size(),
		BlocksProcessed: blocksProcessed,
	}, nil
}
