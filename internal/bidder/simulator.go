package bidder

import (
	"auction-simulator/internal/types"
	"context"
	"math"
	"math/rand"
	"time"
)

// Simulator handles bidder behavior simulation
type Simulator struct {
	bidder *Bidder
}

// NewSimulator creates a new bidder simulator
func NewSimulator(bidder *Bidder) *Simulator {
	return &Simulator{
		bidder: bidder,
	}
}

// EvaluateBid implements the types.Bidder interface
func (s *Simulator) EvaluateBid(ctx context.Context, request *types.BidRequest) (*types.BidResponse, error) {
	// Check if context is already cancelled
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Simulate response time with context awareness
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(time.Duration(s.bidder.SpeedMS) * time.Millisecond):
		// Continue after delay
	}

	// Simple bid decision: use bid chance directly
	if rand.Float64() > s.bidder.BidChance {
		return nil, nil // No bid
	}

	// Calculate bid amount with some randomness
	baseAmount := s.bidder.BaseBid
	variation := (rand.Float64() - 0.5) * s.bidder.BidRange
	bidAmount := math.Max(1.0, baseAmount+variation)
	bidAmount = math.Round(bidAmount*100) / 100 // Round to 2 decimal places

	response := &types.BidResponse{
		BidderID:  s.bidder.ID,
		AuctionID: request.AuctionID,
		Amount:    bidAmount,
		Timestamp: time.Now(),
	}

	return response, nil
}
