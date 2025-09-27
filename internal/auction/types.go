package auction

import (
	"auction-simulator/internal/types"
	"time"
)

// Auction represents a single auction instance
type Auction struct {
	ID         string            `json:"id"`
	Attributes []types.Attribute `json:"attributes"`
	StartTime  time.Time         `json:"start_time"`
	EndTime    time.Time         `json:"end_time"`
	Timeout    time.Duration     `json:"timeout"`
	Winner     *types.Bid        `json:"winner,omitempty"`
	Bids       []types.Bid       `json:"bids"`
	IsComplete bool              `json:"is_complete"`
}
