# FlashBlock with Ethereum JSON-RPC Compatibility

This project implements a lightweight transaction processing system with Ethereum JSON-RPC compatibility.

## Features

- Ethereum JSON-RPC compatibility (`eth_sendRawTransaction`, `eth_getTransactionByHash`)
- Transaction mempool with priority-based ordering
- Block processor that creates blocks from transactions based on gas price (priority fee)

## Ethereum RPC Methods Supported

- `eth_sendRawTransaction`: Submit a raw Ethereum transaction to the mempool
- `eth_getTransactionByHash`: Get transaction details by hash
- `eth_getTransactionReceipt`: Get transaction receipt (simplified implementation)

## Transaction Processing

1. Raw Ethereum transactions are submitted via `eth_sendRawTransaction`
2. Transactions are parsed and stored in the mempool with their Ethereum-specific fields
3. The block processor sorts transactions by gas price (highest first) and creates blocks
4. Processed transactions are removed from the mempool

## Gas Price Priority

Transactions in the mempool are sorted by:
1. Gas price (highest first) for Ethereum transactions
2. Legacy priority value for non-Ethereum transactions

## Usage

Start the server:

```bash
go run cmd/server/main.go
```

Submit a raw transaction using the Ethereum JSON-RPC:

```bash
curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0x...raw transaction hex..."],"id":1}' http://localhost:8545
```

Get transaction details:

```bash
curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0x...transaction hash..."],"id":1}' http://localhost:8545
```

## Dependencies

- github.com/ethereum/go-ethereum: For Ethereum transaction handling and RPC compatibility

## Getting Started

### Prerequisites

- Go 1.24 or higher

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/flashblock.git
cd flashblock

# Build the server
go build -o flashblock ./cmd/server
```

### Running the Server

```bash
# Run with default settings
./flashblock

# Run with custom settings
./flashblock --rpc-addr=:8888 --api-addr=:8889 --block-interval=500ms --log-blocks=true
```

## Usage

### Command-line Options

- `--rpc-addr`: JSON-RPC server address (default: ":8080")
- `--api-addr`: API server address (default: ":8081")
- `--block-interval`: Block creation interval (default: 250ms)
- `--log-blocks`: Log block creation events (default: true)

### JSON-RPC API

All RPC methods are accessible via HTTP POST to the root URL or via WebSocket. Methods are namespaced with the `flash_` prefix.

#### Submit a Transaction

```json
// Request
{
  "jsonrpc": "2.0",
  "method": "flash_submitTransaction",
  "params": {
    "data": "Your transaction data here",
    "priority": 10
  },
  "id": 1
}

// Response
{
  "jsonrpc": "2.0",
  "result": {
    "transaction_id": "f7c8731e8c3a2dc7a009a24e0c71fde0d3486d4a5b3cdb00d4f4bd48757ad2b2",
    "added": true
  },
  "id": 1
}
```

#### Check Transaction Status

```json
// Request
{
  "jsonrpc": "2.0",
  "method": "flash_getTransactionStatus",
  "params": {
    "id": "f7c8731e8c3a2dc7a009a24e0c71fde0d3486d4a5b3cdb00d4f4bd48757ad2b2"
  },
  "id": 2
}

// Response
{
  "jsonrpc": "2.0",
  "result": {
    "exists": true,
    "transaction": {
      "id": "f7c8731e8c3a2dc7a009a24e0c71fde0d3486d4a5b3cdb00d4f4bd48757ad2b2",
      "data": "Your transaction data here",
      "priority": 10,
      "timestamp": "2023-04-01T12:34:56Z"
    }
  },
  "id": 2
}
```

#### Get All Blocks

```json
// Request
{
  "jsonrpc": "2.0",
  "method": "flash_getBlocks",
  "params": {},
  "id": 3
}

// Response
{
  "jsonrpc": "2.0",
  "result": {
    "blocks": [
      {
        "id": "a1b2c3d4e5f6...",
        "transactions": [...],
        "timestamp": "2023-04-01T15:32:45Z",
        "prev_block_id": "f6e5d4c3b2a1..."
      },
      ...
    ],
    "count": 42
  },
  "id": 3
}
```

#### Get Mempool Contents

```json
// Request
{
  "jsonrpc": "2.0",
  "method": "flash_getMempool",
  "params": {},
  "id": 4
}

// Response
{
  "jsonrpc": "2.0",
  "result": {
    "transactions": [
      {
        "id": "tx1",
        "data": "...",
        "priority": 10,
        "timestamp": "2023-04-01T15:32:40Z"
      },
      ...
    ],
    "count": 5
  },
  "id": 4
}
```

#### Get System Status

```json
// Request
{
  "jsonrpc": "2.0",
  "method": "flash_getStatus",
  "params": {},
  "id": 5
}

// Response
{
  "jsonrpc": "2.0",
  "result": {
    "status": "running",
    "uptime": "10m32s",
    "version": "1.0.0",
    "mempool_size": 5,
    "blocks_processed": 42
  },
  "id": 5
}
```

### WebSocket Connection

You can connect to the JSON-RPC server via WebSocket for more efficient communication:

```javascript
// Browser example
const ws = new WebSocket('ws://localhost:8080');

ws.onopen = () => {
  ws.send(JSON.stringify({
    jsonrpc: '2.0',
    method: 'flash_getStatus',
    params: {},
    id: 1
  }));
};

ws.onmessage = (event) => {
  const response = JSON.parse(event.data);
  console.log('Received:', response);
};
```

### HTTP Metrics API

The system provides an HTTP API for monitoring system status and metrics.

#### System Status

```
GET /api/status
```

Response:
```json
{
  "status": "running",
  "uptime": "10m32s",
  "version": "1.0.0",
  "mempool_size": 5,
  "blocks_processed": 42
}
```

#### System Metrics

```
GET /api/metrics
```

Response:
```json
{
  "transactions": {
    "received": 1025,
    "processed": 1020,
    "rejected": 5
  },
  "blocks": {
    "created": 42,
    "avg_block_time": "5.2ms",
    "last_block_at": "2023-04-01T15:32:45Z"
  },
  "performance": {
    "tps": "95.24",
    "uptime_seconds": 632.5,
    "uptime": "10m32s"
  }
}
```

#### Mempool Contents

```
GET /api/mempool
```

Response:
```json
{
  "count": 5,
  "transactions": [
    {
      "id": "f7c8731e8c3a2dc7a009a24e0c71fde0d3486d4a5b3cdb00d4f4bd48757ad2b2",
      "priority": 10,
      "size": 24,
      "time": "2023-04-01T15:32:40Z"
    },
    ...
  ]
}
```

#### Block History

```
GET /api/blocks
```

Response:
```json
{
  "count": 42,
  "blocks": [
    {
      "id": "a1b2c3d4e5f6...",
      "timestamp": "2023-04-01T15:32:45Z",
      "prev_block_id": "f6e5d4c3b2a1...",
      "tx_count": 25,
      "transactions": ["tx1", "tx2", ...]
    },
    ...
  ]
}
```

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make bench
```

### Building Examples

```bash
# Build and run the example client
make run-examples
```

## Architecture

The system consists of several components:

1. **Transaction Model**: Represents a single transaction with data and priority
2. **Block Model**: Contains a collection of transactions
3. **Mempool**: Thread-safe storage for pending transactions
4. **Block Processor**: Periodically creates blocks from transactions in the mempool
5. **JSON-RPC Server**: Handles transaction submissions and status queries using Ethereum RPC implementation
6. **Metrics**: Tracks system performance and statistics
7. **API Server**: Provides HTTP endpoints for monitoring and management

## License

This project is licensed under the MIT License - see the LICENSE file for details. 