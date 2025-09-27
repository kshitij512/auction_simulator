package auction

import (
	"auction-simulator/internal/types"
	"context"
	"time"
)

// Processor handles individual auction execution
type Processor struct {
	auction *Auction
}

// NewProcessor creates a new auction processor
func NewProcessor(auction *Auction) *Processor {
	return &Processor{
		auction: auction,
	}
}

// Run executes the auction and returns the result
func (p *Processor) Run(ctx context.Context) *types.AuctionResult { // Changed to types.AuctionResult
	startTime := time.Now()
	result := &types.AuctionResult{ // Changed to types.AuctionResult
		AuctionID: p.auction.ID,
		StartTime: startTime,
	}

	defer func() {
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(startTime)
	}()

	// Set auction start time
	p.auction.StartTime = startTime

	// Create auction-specific context with timeout
	auctionCtx, cancel := context.WithTimeout(ctx, p.auction.Timeout)
	defer cancel()

	// Simulate auction completion
	select {
	case <-auctionCtx.Done():
		result.Error = auctionCtx.Err()
	case <-time.After(100 * time.Millisecond): // Simulate work
		p.auction.EndTime = time.Now()
		p.auction.IsComplete = true
		result.TotalBids = len(p.auction.Bids)

		if len(p.auction.Bids) > 0 {
			winner := p.determineWinner()
			p.auction.Winner = winner
			result.Winner = winner
		}
	}

	return result
}

// determineWinner selects the highest bidder
func (p *Processor) determineWinner() *types.Bid { // Changed to types.Bid
	if len(p.auction.Bids) == 0 {
		return nil
	}

	winner := &p.auction.Bids[0]
	for i := 1; i < len(p.auction.Bids); i++ {
		if p.auction.Bids[i].Amount > winner.Amount {
			winner = &p.auction.Bids[i]
		}
	}

	return winner
}

// AddBid safely adds a bid to the auction
func (p *Processor) AddBid(bid types.Bid) { // Changed to types.Bid
	p.auction.Bids = append(p.auction.Bids, bid)
}
