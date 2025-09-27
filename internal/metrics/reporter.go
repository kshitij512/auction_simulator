package metrics

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"auction-simulator/internal/types"
)

// Reporter handles output of simulation results
type Reporter struct {
	outputDir string
}

// NewReporter creates a new results reporter
func NewReporter(outputDir string) *Reporter {
	if outputDir == "" {
		outputDir = "output"
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Printf("Warning: Could not create output directory: %v", err)
	}

	return &Reporter{
		outputDir: outputDir,
	}
}

// ReportSummary prints a summary of simulation results
func (r *Reporter) ReportSummary(metrics *SimulationMetrics) {
	separator := strings.Repeat("=", 60)
	fmt.Printf("\n%s\n", separator)
	fmt.Println("AUCTION SIMULATION SUMMARY")
	fmt.Printf("%s\n", separator)

	fmt.Printf("Total Duration: %v\n", metrics.TotalDuration)
	fmt.Printf("Successful Auctions: %d/%d (%.1f%%)\n",
		metrics.SuccessfulAuctions, metrics.TotalAuctions,
		float64(metrics.SuccessfulAuctions)/float64(metrics.TotalAuctions)*100)
	fmt.Printf("Total Bids Received: %d (avg: %.1f per auction)\n",
		metrics.TotalBidsReceived, metrics.AverageBidsPerAuction)
	fmt.Printf("Max Goroutines: %d\n", metrics.MaxGoroutines)
	fmt.Printf("Peak Memory Usage: %.2f MB\n", metrics.MemoryUsageMB)

	fmt.Printf("%s\n", separator)
}

// SaveMetrics writes metrics to a JSON file
func (r *Reporter) SaveMetrics(metrics *SimulationMetrics) error {
	filename := fmt.Sprintf("%s/simulation_metrics_%s.json",
		r.outputDir, time.Now().Format("20060102_150405"))

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create metrics file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(metrics); err != nil {
		return fmt.Errorf("could not encode metrics: %w", err)
	}

	log.Printf("Metrics saved to: %s", filename)
	return nil
}

// SaveAuctionResults writes individual auction results to files
func (r *Reporter) SaveAuctionResults(results []*types.AuctionResult) error {
	for _, result := range results {
		filename := fmt.Sprintf("%s/auction_%s.json", r.outputDir, result.AuctionID)

		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("could not create auction result file: %w", err)
		}

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")

		if err := encoder.Encode(result); err != nil {
			file.Close()
			return fmt.Errorf("could not encode auction result: %w", err)
		}

		file.Close()
	}

	log.Printf("Auction results saved to: %s/", r.outputDir)
	return nil
}
