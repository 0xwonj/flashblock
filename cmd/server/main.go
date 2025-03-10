package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"flashblock/internal/mempool"
	"flashblock/internal/metrics"
	"flashblock/internal/model"
	"flashblock/internal/processor"
	"flashblock/internal/rpc"
)

func main() {
	// Parse command line flags
	var (
		rpcAddr        = flag.String("rpc-addr", ":8080", "JSON-RPC server address")
		blockInterval  = flag.Duration("block-interval", 250*time.Millisecond, "Block creation interval")
		logBlockEvents = flag.Bool("log-blocks", true, "Log block creation events")
		logFile        = flag.String("log-file", "flashblock.log", "Log file path")
	)
	flag.Parse()

	// Set up logger to write to both file and stdout
	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer f.Close()

	// Create a multi writer for both stdout and log file
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("Starting FlashBlock server...")

	// Create metrics
	m := metrics.New()
	log.Println("Metrics initialized")

	// Create mempool
	mp := mempool.New()
	log.Println("Mempool initialized")

	// Create block processor
	processorConfig := &processor.Config{
		Interval: *blockInterval,
	}

	// Add block callback if logging is enabled
	if *logBlockEvents {
		processorConfig.BlockCallback = func(block *model.Block, blockCreationTime time.Duration) {
			m.IncrementBlocksCreated()
			m.IncrementTransactionsProcessed(uint64(len(block.Transactions)))
			m.RecordBlockCreationTime(blockCreationTime)
			m.CalculateMetrics()
			log.Printf("Block created: ID=%s, Transactions=%d, Creation Time=%v", block.ID, len(block.Transactions), blockCreationTime)
		}
	}

	bp := processor.New(mp, processorConfig)
	log.Printf("Block processor initialized with interval: %v", *blockInterval)

	// Create JSON-RPC server with metrics
	rpcServer := rpc.NewServer(mp, *rpcAddr)
	log.Printf("JSON-RPC server initialized with address: %s", *rpcAddr)

	// Set the processor reference in the RPC server
	rpcServer.SetProcessor(bp)

	// Add transaction hook to track metrics
	rpcServer.AddTransactionHook(func(tx *model.Transaction, added bool) {
		m.IncrementTransactionsReceived()
		if !added {
			m.IncrementTransactionsRejected()
		}
	})

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start block processor in a goroutine
	go bp.Start(ctx)

	// Start JSON-RPC server in a goroutine
	go func() {
		if err := rpcServer.Start(ctx); err != nil {
			log.Fatalf("JSON-RPC server error: %v", err)
		}
	}()

	log.Println("System is ready. Press Ctrl+C to stop.")

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Shutdown gracefully
	log.Println("Shutting down...")
	cancel()

	// Give some time for goroutines to finish
	time.Sleep(1 * time.Second)
	log.Println("Server stopped")
}
