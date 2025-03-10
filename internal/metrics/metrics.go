package metrics

import (
	"sync/atomic"
	"time"
	"unsafe"
)

// Metrics tracks system metrics
type Metrics struct {
	// Transaction metrics
	TransactionsReceived  uint64
	TransactionsProcessed uint64
	TransactionsRejected  uint64

	// Block metrics
	BlocksCreated  uint64
	TotalBlockTime time.Duration
	LastBlockTime  time.Time

	// Performance metrics
	StartTime      time.Time
	ProcessedTPS   float64 // Transactions Per Second
	AverageLatency time.Duration
}

// New creates a new metrics instance
func New() *Metrics {
	return &Metrics{
		StartTime:     time.Now(),
		LastBlockTime: time.Now(),
	}
}

// IncrementTransactionsReceived increments the received transactions counter
func (m *Metrics) IncrementTransactionsReceived() {
	atomic.AddUint64(&m.TransactionsReceived, 1)
}

// IncrementTransactionsProcessed increments the processed transactions counter
func (m *Metrics) IncrementTransactionsProcessed(count uint64) {
	atomic.AddUint64(&m.TransactionsProcessed, count)
}

// IncrementTransactionsRejected increments the rejected transactions counter
func (m *Metrics) IncrementTransactionsRejected() {
	atomic.AddUint64(&m.TransactionsRejected, 1)
}

// IncrementBlocksCreated increments the created blocks counter
func (m *Metrics) IncrementBlocksCreated() {
	atomic.AddUint64(&m.BlocksCreated, 1)
}

// RecordBlockCreationTime records the time taken to create a block
func (m *Metrics) RecordBlockCreationTime(duration time.Duration) {
	// Add duration to total time (using nanoseconds for atomic operations)
	atomic.AddUint64((*uint64)(unsafe.Pointer(&m.TotalBlockTime)), uint64(duration.Nanoseconds()))
	m.LastBlockTime = time.Now()
}

// CalculateMetrics calculates derived metrics like TPS and average latency
func (m *Metrics) CalculateMetrics() {
	uptime := time.Since(m.StartTime).Seconds()
	if uptime > 0 {
		m.ProcessedTPS = float64(m.TransactionsProcessed) / uptime
	}

	if m.BlocksCreated > 0 {
		m.AverageLatency = time.Duration(int64(m.TotalBlockTime) / int64(m.BlocksCreated))
	}
}

// GetSnapshot returns a snapshot of the current metrics
func (m *Metrics) GetSnapshot() *Metrics {
	m.CalculateMetrics()

	// Create a copy of the metrics
	snapshot := &Metrics{
		TransactionsReceived:  atomic.LoadUint64(&m.TransactionsReceived),
		TransactionsProcessed: atomic.LoadUint64(&m.TransactionsProcessed),
		TransactionsRejected:  atomic.LoadUint64(&m.TransactionsRejected),
		BlocksCreated:         atomic.LoadUint64(&m.BlocksCreated),
		TotalBlockTime:        m.TotalBlockTime,
		LastBlockTime:         m.LastBlockTime,
		StartTime:             m.StartTime,
		ProcessedTPS:          m.ProcessedTPS,
		AverageLatency:        m.AverageLatency,
	}

	return snapshot
}
