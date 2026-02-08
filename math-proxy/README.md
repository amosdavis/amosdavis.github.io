# Math Proxy - Scrypt to RandomX Mining Proxy

A Golang-based mining proxy that allows miners performing scrypt calculations to perform RandomX calculations instead. This proxy acts as a bridge between scrypt-based mining software and RandomX-based cryptocurrency networks.

## Overview

This proxy server implements the Stratum mining protocol and translates scrypt mining requests into RandomX hash calculations. This enables miners using scrypt-compatible software to mine RandomX-based cryptocurrencies without modifying their mining software.

## Features

- **Stratum Protocol Support**: Implements the standard Stratum mining protocol
- **Request Translation**: Converts scrypt mining requests to RandomX calculations
- **Multiple Client Support**: Can handle multiple concurrent mining connections
- **Configurable**: Supports configuration via JSON file or command-line flags
- **Debug Mode**: Detailed logging for troubleshooting
- **Graceful Shutdown**: Properly handles shutdown signals

## Architecture

```
[Scrypt Miner] <--Stratum--> [Math Proxy] <--RandomX--> [Mining Pool]
                             (This software)
```

The proxy:
1. Listens for incoming Stratum connections from scrypt miners
2. Receives mining.subscribe, mining.authorize, and mining.submit requests
3. Performs RandomX hash calculations instead of scrypt
4. Returns properly formatted Stratum responses

## Installation

### Prerequisites

- Go 1.18 or higher
- Git (for cloning the repository)

### Building from Source

```bash
cd math-proxy
go build -o math-proxy
```

This will create a `math-proxy` executable in the current directory.

## Usage

### Basic Usage

Start the proxy with default settings:

```bash
./math-proxy
```

Default settings:
- Listen address: `0.0.0.0:3333`
- Upstream pool: `localhost:3334`

### Command-Line Options

```bash
./math-proxy [options]

Options:
  -listen string
        Address to listen on (default "0.0.0.0:3333")
  -upstream string
        Upstream pool address (default "localhost:3334")
  -config string
        Path to configuration file (default "config.json")
  -debug
        Enable debug logging
```

### Example Commands

Listen on a specific port:
```bash
./math-proxy -listen 0.0.0.0:8080
```

Enable debug mode:
```bash
./math-proxy -debug
```

Use a custom config file:
```bash
./math-proxy -config /path/to/config.json
```

## Configuration

### Configuration File

Create a `config.json` file:

```json
{
  "listen_addr": "0.0.0.0:3333",
  "upstream_addr": "pool.example.com:3334",
  "debug": false
}
```

Configuration options:
- `listen_addr`: The address and port the proxy listens on
- `upstream_addr`: The upstream RandomX pool address
- `debug`: Enable/disable debug logging

## Connecting Your Miner

Configure your scrypt miner to connect to the proxy:

```bash
# Example with cpuminer
minerd -a scrypt -o stratum+tcp://localhost:3333 -u YOUR_WALLET_ADDRESS -p x

# Example with cgminer
cgminer -o stratum+tcp://localhost:3333 -u YOUR_WALLET_ADDRESS -p x --scrypt
```

## How It Works

### Stratum Protocol

The proxy implements the Stratum mining protocol with these methods:

1. **mining.subscribe**: Client subscribes to mining notifications
2. **mining.authorize**: Client authorizes with username/password
3. **mining.submit**: Client submits a share (this is where translation happens)

### Hash Translation

When a miner submits a scrypt share:
1. Proxy receives the share data (nonce, timestamp, etc.)
2. Extracts the block header information
3. Calculates RandomX hash instead of scrypt hash
4. Validates the hash against the target difficulty
5. Returns success/failure to the miner

### RandomX Implementation

The current implementation uses a simplified RandomX simulation. For production use, you should integrate the actual RandomX library:

- [RandomX Official Repository](https://github.com/tevador/RandomX)
- [Go RandomX Bindings](https://github.com/dominant-strategies/go-randomx)

To use the real RandomX implementation, you would need to:
1. Install the RandomX C library
2. Add CGO bindings
3. Replace the `CalculateRandomXHash` function with actual RandomX calls

## Development

### Project Structure

```
math-proxy/
├── main.go          # Entry point and CLI
├── config.go        # Configuration handling
├── proxy.go         # Proxy server implementation
├── stratum.go       # Stratum protocol structures
├── randomx.go       # RandomX hash calculations
├── proxy_test.go    # Unit tests
├── config.json      # Sample configuration
├── examples/        # Example clients
│   └── test_client.go
└── README.md        # This file
```

### Testing

Build and test the proxy:

```bash
# Build
go build -o math-proxy

# Run with debug mode
./math-proxy -debug

# Test connection (in another terminal)
telnet localhost 3333
```

### Adding Real RandomX Support

To integrate actual RandomX:

1. Install RandomX library:
```bash
# Ubuntu/Debian
sudo apt-get install librandomx-dev

# Or build from source
git clone https://github.com/tevador/RandomX
cd RandomX && mkdir build && cd build
cmake .. && make
sudo make install
```

2. Update `randomx.go` to use CGO bindings:
```go
/*
#cgo LDFLAGS: -lrandomx
#include <randomx.h>
*/
import "C"
```

3. Replace the simulation with actual RandomX calls

## Limitations

- **Simplified RandomX**: Current implementation uses a simplified hash function for demonstration
- **No Upstream Pool Connection**: Doesn't forward shares to an actual upstream pool yet
- **Basic Validation**: Share validation is simplified
- **Single-threaded Hash Calculation**: Real RandomX should use multiple cores

## Future Enhancements

- [ ] Integrate actual RandomX library with CGO
- [ ] Add upstream pool connection and forwarding
- [ ] Implement proper difficulty adjustment
- [ ] Add support for multiple upstream pools
- [ ] Add statistics and monitoring
- [ ] Add web interface for monitoring
- [ ] Implement share validation caching
- [ ] Add support for mining.set_difficulty

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

This project is open source. Please check the repository for license details.

## References

- [Stratum Mining Protocol](https://en.bitcoin.it/wiki/Stratum_mining_protocol)
- [RandomX Algorithm](https://github.com/tevador/RandomX)
- [Scrypt Algorithm](https://en.wikipedia.org/wiki/Scrypt)
- [Mining Pool Protocols](https://braiins.com/stratum-v2)

## Support

For issues and questions, please use the GitHub issue tracker.

## Disclaimer

This is a proof-of-concept implementation. For production use, you should:
- Integrate the actual RandomX library
- Implement proper security measures
- Add comprehensive error handling
- Perform thorough testing
- Consider the legal implications of cryptocurrency mining
