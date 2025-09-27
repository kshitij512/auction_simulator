package types

import (
	"context"
	"time"
)

// Bidder defines the interface that auction package can use
type Bidder interface {
	EvaluateBid(ctx context.Context, request *BidRequest) (*BidResponse, error)
}

// BidRequest contains auction information sent to bidders
type BidRequest struct {
	AuctionID  string        `json:"auction_id"`
	Attributes []float64     `json:"attributes"`
	Timeout    time.Duration `json:"timeout"`
	Timestamp  time.Time     `json:"timestamp"`
}

// BidResponse contains a bidder's response
type BidResponse struct {
	BidderID  string    `json:"bidder_id"`
	AuctionID string    `json:"auction_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

// Attribute represents an auction object characteristic
type Attribute struct {
	ID    int     `json:"id"`
	Value float64 `json:"value"`
}

// Bid represents a bid from a bidder
type Bid struct {
	BidderID  string    `json:"bidder_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

// AuctionResult contains the final outcome of an auction
type AuctionResult struct {
	AuctionID string        `json:"auction_id"`
	Winner    *Bid          `json:"winner,omitempty"`
	TotalBids int           `json:"total_bids"`
	Duration  time.Duration `json:"duration"`
	Error     error         `json:"error,omitempty"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
}
