🏆 Auction Simulator
A high-performance, concurrent auction simulation system built in Go that runs 40 auctions with 100 bidders each, featuring real-time bidding, timeout handling, and comprehensive metrics collection.
📋 Project Overview
This project implements a sophisticated auction simulation system that demonstrates advanced Go concurrency patterns, resource management, and real-time system design. The simulator efficiently handles 40 concurrent auctions with 100 bidders each while maintaining strict resource constraints.
🎯 Assignment Requirements Met
Requirement	Status	Implementation
100 bidders	✅	Configurable bidder management
40 concurrent auctions	✅	Goroutine-based parallel execution
20 attributes per auction	✅	Random attribute generation
Timeout handling	✅	Context-based cancellation
Winner declaration	✅	Highest bid selection
Time measurement	✅	Precise nanosecond timing
Resource standardization	✅	CPU & memory limits
Sample output files	✅	JSON format results
Code clarity	✅	Well-documented, modular
🚀 Quick Start
Prerequisites
Go 1.25 or later
Git
Installation & Running
bash

Clone the repository
git clone 
cd auction-simulator

Run the simulation
go run cmd/simulator/main.go

Or build and run
go build -o simulator cmd/simulator/main.go
./simulator
Expected Output
text
Auction Simulator - Go go1.25.1
Available CPUs: 8, GOMAXPROCS: 8
Simulation Configuration:
Auctions: 40 (concurrent)
Bidders: 100
Attributes: 20 per auction
Auction Timeout: 2s
Resource Standardization:
Max vCPUs: 2
Max Memory: 1024 MB
Max Concurrent Bidders: 200

============================================================
Starting auction simulation...
🏗️ System Architecture
Package Structure
text
auction-simulator/
├── cmd/simulator/          # Application entry point
├── internal/
│   ├── auction/           # Auction management & processing
│   ├── bidder/           # Bidder simulation & behavior
│   ├── config/           # Configuration constants
│   ├── metrics/          # Metrics collection & reporting
│   └── types/            # Shared data types
├── pkg/utils/            # Utility functions
└── output/               # Generated sample files
Key Components
Auction Manager: Orchestrates 40 concurrent auctions
Bidder Manager: Manages 100 simulated bidders
Orchestrator: Handles concurrent execution with resource limits
Metrics Collector: Tracks performance and timing data
Reporter: Generates JSON output files
⚡ Concurrency Design
Resource Control
go
// CPU limiting
runtime.GOMAXPROCS(2)
// Concurrency limiting
semaphore := make(chan struct{}, 200) // Max 200 concurrent bidders
// Timeout handling
ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
Concurrency Patterns Used
Goroutines: 40 auctions × 100 bidders = 4000 concurrent operations
Channels: Buffered channels for bid collection
Context: Propagation for cancellation and timeouts
WaitGroups: Coordination between goroutines
Mutexes: Thread-safe shared state access
🎮 Simulation Flow
Initialization
Create 40 auctions with 20 random attributes each
Initialize 100 bidders with unique characteristics
Set up resource limits and metrics collection
Concurrent Execution
Launch 40 auction goroutines simultaneously
Each auction processes bids from 100 bidders
Semaphore limits to 200 concurrent bidders globally
Bidding Process
Bidders evaluate auction attributes (70-95% bid probability)
Attribute matching influences bid decisions
Context timeouts ensure 2-second completion
Result Collection
Determine highest bidder for each auction
Record comprehensive metrics and timing
Generate JSON output files

📊 Sample Output
Console Summary
text

AUCTION SIMULATION SUMMARY

Total Duration: 2.008s
Successful Auctions: 40/40 (100.0%)
Total Bids Received: 1078 (avg: 26.9 per auction)
Max Goroutines: 4045
Peak Memory Usage: 3.52 MB
Generated Files
Metrics: output/simulation_metrics_20250928_125050.json
Auction Results: output/auction_auction-1.json to auction-40.json
Example Auction Result
json
{
"auction_id": "auction-1",
"winner": {
"bidder_id": "bidder-35",
"amount": 91.47,
"timestamp": "2025-09-28T12:50:50.123456Z"
},
"total_bids": 60,
"duration": "493ms",
"start_time": "2025-09-28T12:50:50.123456Z",
"end_time": "2025-09-28T12:50:50.616456Z"
}
⚙️ Configuration
Resource Limits (internal/config/config.go)
go
const (
DefaultMaxVCPUs    = 2
DefaultMaxMemoryMB = 1024
MaxConcurrentBiddersPerCPU = 100
)
Simulation Parameters
go
const (
TotalBidders         = 100
TotalAuctions        = 40
AttributesPerAuction = 20
DefaultTimeout       = 2 * time.Second
)
Bidder Behavior
go
BidderConfig{
MinBidChance: 0.7,   // 70% minimum bid probability
MaxBidChance: 0.95,  // 95% maximum bid probability
MinBaseBid:   50.0,  // Minimum bid amount
MaxBaseBid:   150.0, // Maximum bid amount
}
🔧 Development
Building from Source
bash

Format code
go fmt ./...

Run tests
go test ./...

Check code quality
go vet ./...

Build for production
go build -o auction-simulator cmd/simulator/main.go
Project Structure Overview
cmd/simulator/main.go: Application entry point and orchestration
internal/auction/: Auction management and processing logic
internal/bidder/: Bidder simulation and behavior algorithms
internal/metrics/: Performance tracking and reporting
internal/config/: Centralized configuration management
internal/types/: Shared data structures and interfaces
📈 Performance Characteristics
Concurrency: Efficiently handles 4000 concurrent operations
Memory Usage: ~3.5MB peak (well within 1GB limit)
Execution Time: ~2 seconds for complete simulation
Success Rate: 100% auction completion
Bid Distribution: 80%+ of auctions receive multiple bids
🎯 Design Highlights
Architecture Decisions
Modular Design: Clear separation of concerns with dedicated packages
Interface-based Communication: Reduced coupling between components
Resource Awareness: Strict CPU and memory constraints
Error Resilience: Graceful handling of timeouts and failures
Go Best Practices
Proper error handling and propagation
Efficient goroutine management
Thread-safe concurrent operations
Clean, documented code structure
🔮 Future Enhancements
Potential Improvements
Dynamic configuration via command-line flags
Advanced bidder strategies and algorithms
Real-time visualization of auction progress
Database integration for result persistence
Comprehensive test suite with benchmarks
Production Features
Structured logging implementation
Health check endpoints
Performance monitoring integration
👨‍💻 Author
Developed with ❤️ using Go, featuring advanced concurrency patterns and resource-aware design.
