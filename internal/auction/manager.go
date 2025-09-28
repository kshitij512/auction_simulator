package auction

import (
	"auction-simulator/internal/config"
	"auction-simulator/internal/types"
	"fmt"
	"log"
	"sync"
	"time"
)

// Manager orchestrates all auctions
type Manager struct {
	config    *config.Config
	auctions  []*Auction
	results   chan *types.AuctionResult // Changed to types.AuctionResult
	metrics   *Metrics
	startTime time.Time
	endTime   time.Time
}

// Metrics tracks auction performance
type Metrics struct {
	mu                sync.RWMutex
	totalAuctions     int
	completedAuctions int
	failedAuctions    int
	totalBids         int
}

// NewManager creates a new auction manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config:   cfg,
		auctions: make([]*Auction, 0, cfg.TotalAuctions),
		results:  make(chan *types.AuctionResult, cfg.TotalAuctions), // Changed to types.AuctionResult
		metrics:  &Metrics{},
	}
}

// InitializeAuctions creates all auction instances
func (m *Manager) InitializeAuctions() error {
	log.Printf("Initializing %d auctions...", m.config.TotalAuctions)

	for i := 0; i < m.config.TotalAuctions; i++ {
		auction := m.createAuction(i)
		m.auctions = append(m.auctions, auction)
	}

	log.Printf("✅ Successfully initialized %d auctions", len(m.auctions))
	return nil
}

// createAuction generates a single auction with random attributes
func (m *Manager) createAuction(id int) *Auction {
	attributes := make([]types.Attribute, m.config.AttributesPerAuction) // Changed to types.Attribute
	for j := 0; j < m.config.AttributesPerAuction; j++ {
		attributes[j] = types.Attribute{ // Changed to types.Attribute
			ID:    j,
			Value: generateAttributeValue(),
		}
	}

	return &Auction{
		ID:         fmt.Sprintf("auction-%d", id+1),
		Attributes: attributes,
		Timeout:    m.config.AuctionTimeout,
		Bids:       make([]types.Bid, 0), // Changed to types.Bid
		IsComplete: false,
	}
}

// GetAuctions returns all auctions (thread-safe)
func (m *Manager) GetAuctions() []*Auction {
	return m.auctions
}

// GetMetrics returns current metrics
func (m *Manager) GetMetrics() *Metrics {
	return m.metrics
}

// RecordResult processes auction results
func (m *Manager) RecordResult(result *types.AuctionResult) {
	m.metrics.mu.Lock()
	defer m.metrics.mu.Unlock()

	m.metrics.totalAuctions++
	if result.Error != nil {
		m.metrics.failedAuctions++
	} else {
		m.metrics.completedAuctions++
	}
	m.metrics.totalBids += result.TotalBids
}

// generateAttributeValue creates random attribute values
func generateAttributeValue() float64 {
	return float64(int(time.Now().UnixNano())%10000) / 100.0
}
