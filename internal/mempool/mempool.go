package mempool

import (
	"sort"
	"sync"

	"flashblock/internal/model"
)

// TransactionHook is a function called when a transaction is processed
type TransactionHook func(*model.Transaction, bool)

// Mempool stores pending transactions in memory
type Mempool struct {
	transactions map[string]*model.Transaction
	hooks        []TransactionHook
	mu           sync.RWMutex
}

// New creates a new empty mempool
func New() *Mempool {
	return &Mempool{
		transactions: make(map[string]*model.Transaction),
		hooks:        make([]TransactionHook, 0),
	}
}

// AddTransactionHook adds a hook to be called when a transaction is added to the mempool
func (mp *Mempool) AddTransactionHook(hook TransactionHook) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.hooks = append(mp.hooks, hook)
}

// AddTransaction adds a new transaction to the mempool
func (mp *Mempool) AddTransaction(tx *model.Transaction) bool {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	// Check if transaction already exists
	if _, exists := mp.transactions[tx.ID]; exists {
		return false
	}

	// Add transaction to mempool
	mp.transactions[tx.ID] = tx

	// Execute transaction hooks outside the lock
	added := true
	go mp.executeHooks(tx, added)

	return added
}

// executeHooks runs all registered hooks for a transaction
func (mp *Mempool) executeHooks(tx *model.Transaction, added bool) {
	mp.mu.RLock()
	hooks := make([]TransactionHook, len(mp.hooks))
	copy(hooks, mp.hooks)
	mp.mu.RUnlock()

	for _, hook := range hooks {
		hook(tx, added)
	}
}

// GetTransaction retrieves a transaction by ID
func (mp *Mempool) GetTransaction(id string) (*model.Transaction, bool) {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	tx, exists := mp.transactions[id]
	return tx, exists
}

// GetAllTransactions returns all transactions currently in the mempool
func (mp *Mempool) GetAllTransactions() []*model.Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	// Create a slice to hold transactions
	txs := make([]*model.Transaction, 0, len(mp.transactions))

	// Add all transactions to the slice
	for _, tx := range mp.transactions {
		txs = append(txs, tx)
	}

	return txs
}

// GetSortedTransactions returns all transactions sorted by priority (high to low)
func (mp *Mempool) GetSortedTransactions() []*model.Transaction {
	transactions := mp.GetAllTransactions()

	// Sort transactions by priority (high to low)
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Priority > transactions[j].Priority
	})

	return transactions
}

// RemoveTransactions removes transactions with the given IDs from the mempool
func (mp *Mempool) RemoveTransactions(ids []string) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	for _, id := range ids {
		delete(mp.transactions, id)
	}
}

// Clear removes all transactions from the mempool
func (mp *Mempool) Clear() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.transactions = make(map[string]*model.Transaction)
}

// Size returns the number of transactions in the mempool
func (mp *Mempool) Size() int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return len(mp.transactions)
}
