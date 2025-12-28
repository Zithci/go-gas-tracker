# Go Gas Tracker

Live multi-chain gas price tracker with real-time updates.

## Features

- **5 EVM chains**: Ethereum, Arbitrum, Base, Optimism, Polygon
- **Parallel fetching**: Goroutines + channels for concurrent RPC calls
- **Live updates**: 3 update cycles with timestamps
- **Statistics**: Min/Max/Avg gas prices across updates
- **Fast**: ~12s total execution time

## Tech Stack

- Go 1.21+
- Goroutines for concurrency
- Channels for communication
- Context for timeout handling
- JSON-RPC for blockchain data

## Concepts Demonstrated

- Concurrent programming (goroutines, channels, WaitGroup)
- HTTP requests with context
- JSON marshaling/unmarshaling
- Hex to decimal conversion
- Statistics calculation


## Usage
```bash
go run gasfaster.go
```

## Output
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       LIVE GAS PRICE TRACKER
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📊 Update #1 [HH:MM:SS]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Arbitrum     0.01 gwei
Base         0.00 gwei
Optimism     0.00 gwei
Ethereum     0.03 gwei
Polygon      63.83 gwei

📊 Update #2 [HH:MM:SS]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Base         0.00 gwei
Arbitrum     0.01 gwei
Optimism     0.00 gwei
Polygon      63.79 gwei
Ethereum     0.03 gwei

📊 Update #3 [HH:MM:SS]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Arbitrum     0.01 gwei
Base         0.00 gwei
Ethereum     0.03 gwei
Optimism     0.00 gwei
Polygon      63.90 gwei

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       STATISTICS (3 updates)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Polygon      Min:  63.79 | Max:  63.90 | Avg:  63.84
Ethereum     Min:   0.03 | Max:   0.03 | Avg:   0.03
Arbitrum     Min:   0.01 | Max:   0.01 | Avg:   0.01
Base         Min:   0.00 | Max:   0.00 | Avg:   0.00
Optimism     Min:   0.00 | Max:   0.00 | Avg:   0.00
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⏱️  Total Time:  11.6s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

