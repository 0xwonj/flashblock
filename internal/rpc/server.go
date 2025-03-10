package rpc

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"flashblock/internal/mempool"
	"flashblock/internal/processor"
	ethapi "flashblock/internal/rpc/eth"
	flashapi "flashblock/internal/rpc/flash"

	"github.com/ethereum/go-ethereum/rpc"
)

// TransactionHook is a function called when a transaction is processed
type TransactionHook = mempool.TransactionHook

// Server represents a JSON-RPC server
type Server struct {
	mempool   *mempool.Mempool
	processor *processor.BlockProcessor
	addr      string
	rpcServer *rpc.Server
}

// NewServer creates a new JSON-RPC server
func NewServer(mempool *mempool.Mempool, addr string) *Server {
	server := &Server{
		mempool: mempool,
		addr:    addr,
	}

	return server
}

// SetProcessor sets the block processor reference
func (s *Server) SetProcessor(bp *processor.BlockProcessor) {
	s.processor = bp
}

// AddTransactionHook adds a hook to be called when a transaction is processed
func (s *Server) AddTransactionHook(hook TransactionHook) {
	// Register hook with mempool directly
	s.mempool.AddTransactionHook(hook)
}

// Start starts the JSON-RPC server
func (s *Server) Start(ctx context.Context) error {
	// Create a new RPC server
	s.rpcServer = rpc.NewServer()

	// Create and register Flash API (empty hooks since we now register them with mempool)
	flashAPI := flashapi.NewAPI(s.mempool, s.processor, nil)
	if err := s.rpcServer.RegisterName("flash", flashAPI); err != nil {
		return err
	}

	// Create and register Ethereum API (empty hooks since we now register them with mempool)
	ethAPI := ethapi.NewAPI(s.mempool, nil)
	if err := s.rpcServer.RegisterName("eth", ethAPI); err != nil {
		return err
	}

	// Set up HTTP server with WebSocket support
	mux := http.NewServeMux()

	// Handle JSON-RPC requests via HTTP POST
	mux.Handle("/", s.rpcServer)

	// Handle Websocket requests
	mux.Handle("/ws", s.rpcServer.WebsocketHandler([]string{"*"}))

	// Create and configure HTTP server
	httpServer := &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	// Create TCP listener
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	// Start server in a goroutine
	go func() {
		log.Printf("JSON-RPC server listening on %s (HTTP and WebSocket)", s.addr)
		if err := httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("JSON-RPC server error: %v", err)
		}
	}()

	// Wait for context cancellation to stop server
	<-ctx.Done()
	log.Println("Shutting down JSON-RPC server...")

	// Create a timeout context for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return httpServer.Shutdown(shutdownCtx)
}
