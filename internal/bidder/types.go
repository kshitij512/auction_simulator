package bidder

// Bidder represents a simulated bidder
type Bidder struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	BidChance  float64 `json:"bid_chance"`
	BaseBid    float64 `json:"base_bid"`
	BidRange   float64 `json:"bid_range"`
	SpeedMS    int     `json:"speed_ms"`
	Attributes []int   `json:"attributes"`
}

// BidderConfig holds configuration for bidder behavior
type BidderConfig struct {
	MinBidChance float64 `json:"min_bid_chance"`
	MaxBidChance float64 `json:"max_bid_chance"`
	MinBaseBid   float64 `json:"min_base_bid"`
	MaxBaseBid   float64 `json:"max_base_bid"`
	MinSpeedMS   int     `json:"min_speed_ms"`
	MaxSpeedMS   int     `json:"max_speed_ms"`
}

// DefaultBidderConfig returns defaults for bidder behavior
func DefaultBidderConfig() *BidderConfig {
	return &BidderConfig{
		MinBidChance: 0.6,
		MaxBidChance: 0.8,
		MinBaseBid:   50.0,
		MaxBaseBid:   150.0,
		MinSpeedMS:   5,
		MaxSpeedMS:   250,
	}
}
