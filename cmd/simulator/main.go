package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"auction-simulator/internal/auction"
	"auction-simulator/internal/bidder"
	"auction-simulator/internal/config"
	"auction-simulator/internal/metrics"
	"auction-simulator/internal/types"
	"auction-simulator/pkg/utils"
)

func main() {
	// Initialize random seed
	utils.InitRandom()

	// Display environment information
	fmt.Printf("Auction Simulator - Go %s\n", runtime.Version())
	fmt.Printf("Available CPUs: %d, GOMAXPROCS: %d\n",
		runtime.NumCPU(), runtime.GOMAXPROCS(0))

	// Load configuration with resource standardization
	cfg := config.DefaultConfig()

	fmt.Printf("\nSimulation Configuration:\n")
	fmt.Printf("   Auctions: %d (concurrent)\n", cfg.TotalAuctions)
	fmt.Printf("   Bidders: %d\n", cfg.TotalBidders)
	fmt.Printf("   Attributes: %d per auction\n", cfg.AttributesPerAuction)
	fmt.Printf("   Auction Timeout: %v\n", cfg.AuctionTimeout)

	fmt.Printf("\nResource Standardization:\n")
	fmt.Printf("   Max vCPUs: %d\n", cfg.ResourceLimits.MaxVCPUs)
	fmt.Printf("   Max Memory: %d MB\n", cfg.ResourceLimits.MaxMemoryMB)
	fmt.Printf("   Max Concurrent Bidders: %d\n", cfg.ResourceLimits.MaxConcurrentBidders)

	// Validate environment
	if err := validateEnvironment(cfg); err != nil {
		log.Fatalf("Environment validation failed: %v", err)
	}

	// Set CPU limit
	runtime.GOMAXPROCS(cfg.ResourceLimits.MaxVCPUs)
	fmt.Printf("\nResource limits applied: GOMAXPROCS=%d\n", runtime.GOMAXPROCS(0))

	// Use constant strings for the separators
	separator := strings.Repeat("=", 60)
	fmt.Printf("\n%s\n", separator)
	fmt.Println("Starting auction simulation...")
	fmt.Printf("%s\n", separator)

	// Run the simulation
	if err := runSimulation(cfg); err != nil {
		log.Fatalf("Simulation failed: %v", err)
	}

	log.Println("Simulation completed successfully")
	os.Exit(0)
}

func validateEnvironment(cfg *config.Config) error {
	if runtime.Version() < "go1.25" {
		return fmt.Errorf("requires Go 1.25 or later, current: %s", runtime.Version())
	}

	if cfg.ResourceLimits.MaxVCPUs < 1 {
		return fmt.Errorf("invalid CPU limit: %d", cfg.ResourceLimits.MaxVCPUs)
	}

	if cfg.ResourceLimits.MaxMemoryMB < 100 {
		return fmt.Errorf("insufficient memory limit: %d MB", cfg.ResourceLimits.MaxMemoryMB)
	}

	return nil
}

func runSimulation(cfg *config.Config) error {
	// Record overall simulation start time
	simulationStart := time.Now()

	// Initialize components
	auctionManager := auction.NewManager(cfg)
	bidderManager := bidder.NewManager(cfg)
	metricsCollector := metrics.NewCollector()
	reporter := metrics.NewReporter("output")

	// Start metrics collection
	metricsCollector.Start()

	fmt.Println("Initializing simulation components...")

	// Initialize auctions and bidders
	if err := auctionManager.InitializeAuctions(); err != nil {
		return fmt.Errorf("failed to initialize auctions: %w", err)
	}

	if err := bidderManager.InitializeBidders(); err != nil {
		return fmt.Errorf("failed to initialize bidders: %w", err)
	}

	fmt.Println("Initialization complete:")
	fmt.Printf("   Auctions: %d\n", len(auctionManager.GetAuctions()))
	fmt.Printf("   Bidders: %d\n", len(bidderManager.GetBidders()))

	// Create orchestrator for concurrent auction execution
	orchestrator := auction.NewOrchestrator(cfg, auctionManager, bidderManager)

	fmt.Printf("\nStarting %d concurrent auctions with %d bidders each...\n",
		cfg.TotalAuctions, cfg.TotalBidders)
	fmt.Printf("Auction timeout: %v\n", cfg.AuctionTimeout)
	fmt.Printf("Resource limit: %d concurrent bidders\n", cfg.ResourceLimits.MaxConcurrentBidders)

	// Run all auctions concurrently
	auctionResults, err := runConcurrentAuctions(orchestrator, cfg)
	if err != nil {
		return fmt.Errorf("auction execution failed: %w", err)
	}

	// Calculate total simulation duration
	totalDuration := time.Since(simulationStart)

	// Stop metrics and report results
	simulationMetrics := metricsCollector.Stop()
	simulationMetrics.TotalDuration = totalDuration
	simulationMetrics.TotalAuctions = cfg.TotalAuctions
	simulationMetrics.TotalBidders = cfg.TotalBidders

	// Calculate auction statistics
	successfulAuctions := 0
	totalBidsReceived := 0
	for _, result := range auctionResults {
		if result.Error == nil {
			successfulAuctions++
		}
		totalBidsReceived += result.TotalBids
	}

	simulationMetrics.SuccessfulAuctions = successfulAuctions
	simulationMetrics.FailedAuctions = len(auctionResults) - successfulAuctions
	simulationMetrics.TotalBidsReceived = totalBidsReceived
	if len(auctionResults) > 0 {
		simulationMetrics.AverageBidsPerAuction = float64(totalBidsReceived) / float64(len(auctionResults))
	}

	// Report results using constant strings
	separator := strings.Repeat("=", 60)
	fmt.Printf("\n%s\n", separator)
	reporter.ReportSummary(simulationMetrics)

	// Save results to files
	if err := reporter.SaveMetrics(simulationMetrics); err != nil {
		log.Printf("Warning: Could not save metrics: %v", err)
	}

	if err := reporter.SaveAuctionResults(auctionResults); err != nil {
		log.Printf("Warning: Could not save auction results: %v", err)
	}

	// Print detailed auction results
	printAuctionDetails(auctionResults)

	fmt.Printf("\nSimulation completed in %v\n", totalDuration)
	return nil
}

func runConcurrentAuctions(orchestrator *auction.Orchestrator, cfg *config.Config) ([]*types.AuctionResult, error) {
	// Create context for the entire simulation
	ctx := context.Background()

	// Run all auctions concurrently
	results, err := orchestrator.RunAllAuctions(ctx)
	if err != nil {
		return nil, fmt.Errorf("error running auctions: %w", err)
	}

	return results, nil
}

func printAuctionDetails(results []*types.AuctionResult) {
	fmt.Printf("\nDetailed Auction Results:\n")

	// Use constant strings for formatting
	lineSeparator := strings.Repeat("-", 80)
	fmt.Printf("%s\n", lineSeparator)
	fmt.Printf("%-12s %-8s %-12s %-15s %-20s\n",
		"Auction ID", "Bids", "Duration", "Winner", "Status")
	fmt.Printf("%s\n", lineSeparator)

	successful := 0
	for _, result := range results {
		status := "Success"
		if result.Error != nil {
			status = fmt.Sprintf("Error: %v", result.Error)
		} else {
			successful++
		}

		winnerInfo := "None"
		if result.Winner != nil {
			winnerInfo = fmt.Sprintf("%s ($%.2f)", result.Winner.BidderID, result.Winner.Amount)
		}

		fmt.Printf("%-12s %-8d %-12v %-15s %-20s\n",
			result.AuctionID,
			result.TotalBids,
			result.Duration.Round(time.Millisecond),
			winnerInfo,
			status)
	}

	fmt.Printf("%s\n", lineSeparator)
	fmt.Printf("Summary: %d/%d auctions successful (%.1f%%)\n",
		successful, len(results), float64(successful)/float64(len(results))*100)
}
