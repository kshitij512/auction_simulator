package bidder

import (
	"auction-simulator/internal/config"
	"auction-simulator/pkg/utils"
	"fmt"
	"log"
	"math/rand"
)

// Manager handles all bidders
type Manager struct {
	config  *config.Config
	bidders []*Bidder
}

// NewManager creates a new bidder manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config:  cfg,
		bidders: make([]*Bidder, 0, cfg.TotalBidders),
	}
}

// InitializeBidders creates all bidder instances
func (m *Manager) InitializeBidders() error {
	log.Printf("Initializing %d bidders...", m.config.TotalBidders)
	bidderConfig := DefaultBidderConfig()

	for i := 0; i < m.config.TotalBidders; i++ {
		bidder := m.createBidder(i, bidderConfig)
		m.bidders = append(m.bidders, bidder)
	}

	log.Printf("✅ Successfully initialized %d bidders", len(m.bidders))
	return nil
}

// createBidder generates a single bidder
func (m *Manager) createBidder(id int, config *BidderConfig) *Bidder {
	return &Bidder{
		ID:         fmt.Sprintf("bidder-%d", id+1),
		Name:       fmt.Sprintf("Bidder %d", id+1),
		BidChance:  utils.RandomFloat(config.MinBidChance, config.MaxBidChance),
		BaseBid:    utils.RandomFloat(config.MinBaseBid, config.MaxBaseBid),
		BidRange:   utils.RandomFloat(5.0, 20.0),
		SpeedMS:    utils.RandomInt(config.MinSpeedMS, config.MaxSpeedMS),
		Attributes: m.generatePreferredAttributes(),
	}
}

// GetBidders returns all bidders
func (m *Manager) GetBidders() []*Bidder {
	return m.bidders
}

// GetBidder returns a specific bidder by ID
func (m *Manager) GetBidder(id string) (*Bidder, error) {
	for _, bidder := range m.bidders {
		if bidder.ID == id {
			return bidder, nil
		}
	}
	return nil, fmt.Errorf("bidder not found: %s", id)
}

// GetBidderSimulators returns simulator instances for all bidders
func (m *Manager) GetBidderSimulators() []*Simulator {
	bidders := m.GetBidders()
	simulators := make([]*Simulator, len(bidders))
	for i, bidder := range bidders {
		simulators[i] = NewSimulator(bidder)
	}
	return simulators
}

// generatePreferredAttributes creates random attribute preferences
func (m *Manager) generatePreferredAttributes() []int {
	numPreferences := rand.Intn(6) + 3 // 3-8 preferred attributes
	preferences := make([]int, numPreferences)
	for i := 0; i < numPreferences; i++ {
		preferences[i] = rand.Intn(20) // 0-19 attribute IDs
	}
	return preferences
}
