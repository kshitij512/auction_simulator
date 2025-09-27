package metrics

import (
	"runtime"
	"sync"
	"time"
)

// Collector gathers simulation metrics
type Collector struct {
	mu            sync.RWMutex
	metrics       *SimulationMetrics
	startTime     time.Time
	maxMemory     uint64
	maxGoroutines int
}

// NewCollector creates a new metrics collector
func NewCollector() *Collector {
	return &Collector{
		metrics: &SimulationMetrics{},
	}
}

// Start begins metrics collection
func (c *Collector) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.startTime = time.Now()
	c.metrics.StartTime = c.startTime
	c.maxMemory = 0
	c.maxGoroutines = 0

	// Start background monitoring
	go c.monitorResources()
}

// Stop ends metrics collection and finalizes results
func (c *Collector) Stop() *SimulationMetrics {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics.EndTime = time.Now()
	c.metrics.TotalDuration = c.metrics.EndTime.Sub(c.metrics.StartTime)
	c.metrics.MemoryUsageMB = float64(c.maxMemory) / 1024 / 1024
	c.metrics.MaxGoroutines = c.maxGoroutines

	return c.metrics
}

// RecordAuctionResult updates metrics with auction results
func (c *Collector) RecordAuctionResult(successful bool, bidCount int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if successful {
		c.metrics.SuccessfulAuctions++
	} else {
		c.metrics.FailedAuctions++
	}

	c.metrics.TotalBidsReceived += bidCount
	c.metrics.TotalAuctions = c.metrics.SuccessfulAuctions + c.metrics.FailedAuctions

	if c.metrics.TotalAuctions > 0 {
		c.metrics.AverageBidsPerAuction = float64(c.metrics.TotalBidsReceived) / float64(c.metrics.TotalAuctions)
	}
}

// monitorResources periodically checks resource usage
func (c *Collector) monitorResources() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.updateResourceStats()
		}
	}
}

// updateResourceStats updates current resource usage
func (c *Collector) updateResourceStats() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Get memory statistics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	if m.Alloc > c.maxMemory {
		c.maxMemory = m.Alloc
	}

	// Get goroutine count
	goroutineCount := runtime.NumGoroutine()
	if goroutineCount > c.maxGoroutines {
		c.maxGoroutines = goroutineCount
	}
}
