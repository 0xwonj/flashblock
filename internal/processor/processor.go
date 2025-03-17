package processor

import (
	"context"
	"log"
	"sort"
	"time"

	"flashblock/internal/attest"
	"flashblock/internal/mempool"
	"flashblock/internal/model"
)

// BlockProcessor processes transactions from the mempool and creates blocks
type BlockProcessor struct {
	mempool         *mempool.Mempool
	latestBlockID   string
	processedBlocks []*model.Block
	blockCallback   func(*model.Block, time.Duration)
	config          *Config
	tdxProvider     *attest.TDXProvider // TDX provider for quote generation
}

// Config holds configuration for the block processor
type Config struct {
	Interval        time.Duration
	BlockCallback   func(*model.Block, time.Duration)
	MaxStoredBlocks int  // Maximum number of recent blocks to keep in memory
	EnableTDXQuote  bool // Whether to generate TDX quotes for blocks
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Interval:        250 * time.Millisecond,
		MaxStoredBlocks: 100, // Default to storing the 100 most recent blocks
		EnableTDXQuote:  false,
	}
}

// New creates a new block processor
func New(mempool *mempool.Mempool, config *Config) *BlockProcessor {
	if config == nil {
		config = DefaultConfig()
	}

	bp := &BlockProcessor{
		mempool:         mempool,
		latestBlockID:   "",
		processedBlocks: make([]*model.Block, 0),
		blockCallback:   config.BlockCallback,
		config:          config,
	}

	// Initialize TDX provider if quote generation is enabled
	if config.EnableTDXQuote {
		provider, err := attest.NewTDXProvider()
		if err != nil {
			log.Printf("Warning: Failed to initialize TDX provider: %v. TDX quotes will be disabled.", err)
			// Disable TDX quote generation if not supported
			bp.config.EnableTDXQuote = false
		} else {
			bp.tdxProvider = provider
			log.Println("TDX quote provider initialized successfully")
		}
	}

	return bp
}

// Start begins the block processing loop
func (bp *BlockProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(bp.config.Interval)
	defer ticker.Stop()

	log.Printf("Block processor started with interval: %v", bp.config.Interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Block processor stopped")
			return
		case <-ticker.C:
			go bp.processNextBlock()
		}
	}
}

// processNextBlock creates a new block from the mempool transactions
func (bp *BlockProcessor) processNextBlock() {
	// Start measuring block creation time
	startTime := time.Now()

	// Get all transactions from mempool
	transactions := bp.mempool.GetAllTransactions()

	// Skip if there are no transactions
	if len(transactions) == 0 {
		return
	}

	// Sort transactions by priority fee (high to low)
	sort.Slice(transactions, func(i, j int) bool {
		// Compare transactions by priority (higher priority first)
		return transactions[i].Priority > transactions[j].Priority
	})

	// Create a new block
	block := model.NewBlock(transactions, bp.latestBlockID)

	// Generate TDX quote if enabled
	if bp.config.EnableTDXQuote && bp.tdxProvider != nil {
		bp.generateTDXQuoteForBlock(block)
	}

	// Update latest block ID
	bp.latestBlockID = block.ID

	// Add block to processed blocks
	bp.processedBlocks = append(bp.processedBlocks, block)

	// Limit the number of stored blocks to prevent memory growth
	if len(bp.processedBlocks) > bp.config.MaxStoredBlocks {
		// Remove oldest blocks to maintain the limit
		excess := len(bp.processedBlocks) - bp.config.MaxStoredBlocks
		bp.processedBlocks = bp.processedBlocks[excess:]
	}

	// Remove processed transactions from mempool
	txIDs := make([]string, len(transactions))
	for i, tx := range transactions {
		txIDs[i] = tx.ID
	}
	bp.mempool.RemoveTransactions(txIDs)

	// Calculate block creation time
	blockCreationTime := time.Since(startTime)

	// Call the callback if set
	if bp.blockCallback != nil {
		bp.blockCallback(block, blockCreationTime)
	}
}

// generateTDXQuoteForBlock generates a TDX quote for the given block
func (bp *BlockProcessor) generateTDXQuoteForBlock(block *model.Block) {
	// Use block ID as user data for the quote
	var quoteData []byte
	var err error

	quoteData, err = bp.tdxProvider.GetQuote([]byte(block.ID))
	if err != nil {
		log.Printf("Failed to generate TDX quote for block %s: %v", block.ID, err)
		return
	}

	block.TDXQuote = quoteData
	log.Printf("Generated TDX quote for block %s (%d bytes)", block.ID, len(quoteData))
}

// GetProcessedBlocks returns all blocks that have been processed
func (bp *BlockProcessor) GetProcessedBlocks() []*model.Block {
	return bp.processedBlocks
}
