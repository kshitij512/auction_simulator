package metrics

import "time"

// SimulationMetrics holds overall simulation metrics
type SimulationMetrics struct {
	TotalAuctions int           `json:"total_auctions"`
	TotalBidders  int           `json:"total_bidders"`
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	TotalDuration time.Duration `json:"total_duration"`

	// Resource usage
	MaxGoroutines int     `json:"max_goroutines"`
	MemoryUsageMB float64 `json:"memory_usage_mb"`

	// Auction statistics
	SuccessfulAuctions    int     `json:"successful_auctions"`
	FailedAuctions        int     `json:"failed_auctions"`
	TotalBidsReceived     int     `json:"total_bids_received"`
	AverageBidsPerAuction float64 `json:"average_bids_per_auction"`
}

// ResourceUsage tracks system resource consumption
type ResourceUsage struct {
	Timestamp      time.Time `json:"timestamp"`
	GoroutineCount int       `json:"goroutine_count"`
	MemoryMB       float64   `json:"memory_mb"`
	CPUPercent     float64   `json:"cpu_percent"`
}
