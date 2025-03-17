# FlashBlock

A high-performance, lightweight TEE-based secure block builder for Ethereum.

## Overview

FlashBlock is a lightweight block builder that uses a TEE to securely build blocks. It is designed to be used in conjunction with an existing Ethereum client to build blocks at a configurable interval.

## Features

- **Ultra-fast Blocks**: Generate blocks at configurable millisecond intervals (default: 250ms)
- **JSON-RPC API**: Compatible interface for submitting transactions and querying blocks
- **WebSocket Support**: Real-time block notifications
- **Mempool Management**: Efficient transaction queuing and processing
- **Metrics**: Built-in performance tracking
- **Low Latency**: Minimize transaction confirmation times for testing
- **Configurable**: Adjust block intervals and other parameters to suit testing needs

## Getting Started

### Prerequisites

- Go 1.24 or higher

### Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/flashblock.git
cd flashblock
```

### Building

Build the server and client:

```bash
make build
```

### Running the Server

```bash
make run
```

With custom configuration:

```bash
make run-custom
```

Or directly:

```bash
./build/flashblock --rpc-addr=:8888 --block-interval=500ms
```

### Running the Client

```bash
make run-client
```

## Configuration

FlashBlock can be configured via command line flags:

- `--rpc-addr`: JSON-RPC server address (default: `:8080`)
- `--block-interval`: Block creation interval (default: `250ms`)
- `--log-blocks`: Enable block creation event logging (default: `true`)
- `--log-file`: Log file path (default: `flashblock.log`)

A sample configuration file (`config.yaml`) is also available for client workload testing:

```yaml
# Number of concurrent clients
num_clients: 5

# Requests per second per client
requests_per_second: 10

# Total duration of the test in seconds
duration_seconds: 60

# Server URL
server_url: "http://localhost:8080"
```

## Usage Examples

Check the `examples/` directory for sample code showing how to interact with FlashBlock:

- Basic client for submitting transactions
- WebSocket client for real-time block notifications

## Development

### Project Structure

- `cmd/`: Application entry points
  - `server/`: FlashBlock server
  - `client/`: Test client implementation
- `internal/`: Internal packages
  - `mempool/`: Transaction queue management
  - `processor/`: Block creation and transaction processing
  - `rpc/`: JSON-RPC API implementation
  - `model/`: Data structures
  - `metrics/`: Performance measurement
  - `eth/`: Ethereum compatibility

### Running Tests

```bash
make test
```

With coverage:

```bash
make test-coverage
```

Run benchmarks:

```bash
make bench
```

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
