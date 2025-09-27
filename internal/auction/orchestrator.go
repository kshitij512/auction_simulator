package auction

import (
	"auction-simulator/internal/types"
	"context"
	"log"
	"sync"
	"time"

	"auction-simulator/internal/bidder"
	"auction-simulator/internal/config"
)

// Orchestrator manages concurrent auction execution
type Orchestrator struct {
	config         *config.Config
	auctionManager *Manager
	bidderManager  *bidder.Manager
	semaphore      chan struct{}
}

// NewOrchestrator creates a new auction orchestrator
func NewOrchestrator(cfg *config.Config, auctionMgr *Manager, bidderMgr *bidder.Manager) *Orchestrator {
	return &Orchestrator{
		config:         cfg,
		auctionManager: auctionMgr,
		bidderManager:  bidderMgr,
		semaphore:      make(chan struct{}, cfg.ResourceLimits.MaxConcurrentBidders),
	}
}

// RunAllAuctions executes all auctions concurrently
func (o *Orchestrator) RunAllAuctions(ctx context.Context) ([]*types.AuctionResult, error) {
	auctions := o.auctionManager.GetAuctions()
	results := make([]*types.AuctionResult, len(auctions))

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstError error

	batchCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	startTime := time.Now()
	log.Printf("🏁 Starting %d auctions concurrently", len(auctions))

	for i, auction := range auctions {
		wg.Add(1)

		go func(idx int, auct *Auction) {
			defer wg.Done()
			result := o.runSingleAuction(batchCtx, auct, idx)

			mu.Lock()
			results[idx] = result
			if result.Error != nil && firstError == nil {
				firstError = result.Error
			}
			mu.Unlock()
		}(i, auction)
	}

	wg.Wait()
	log.Printf("✅ All auctions completed in %v", time.Since(startTime))
	return results, firstError
}

// runSingleAuction executes a single auction
func (o *Orchestrator) runSingleAuction(ctx context.Context, auct *Auction, auctionIndex int) *types.AuctionResult {
	processor := NewProcessor(auct)

	auctionCtx, cancel := context.WithTimeout(ctx, auct.Timeout)
	defer cancel()

	log.Printf("🎯 Starting auction %s (timeout: %v)", auct.ID, auct.Timeout)

	result := &types.AuctionResult{
		AuctionID: auct.ID,
		StartTime: time.Now(),
	}

	// Collect bids
	o.collectBids(auctionCtx, processor, auct)

	// Determine winner
	if len(auct.Bids) > 0 {
		winner := processor.determineWinner()
		auct.Winner = winner
		result.Winner = winner
	}

	result.TotalBids = len(auct.Bids)
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	auct.IsComplete = true

	log.Printf("📊 Auction %s completed: %d bids", auct.ID, result.TotalBids)
	return result
}

// collectBids collects bids from all bidders using simulators
func (o *Orchestrator) collectBids(ctx context.Context, processor *Processor, auct *Auction) {
	var wg sync.WaitGroup
	bidCh := make(chan types.Bid, o.config.TotalBidders)

	attributeValues := make([]float64, len(auct.Attributes))
	for i, attr := range auct.Attributes {
		attributeValues[i] = attr.Value
	}

	bidRequest := &types.BidRequest{
		AuctionID:  auct.ID,
		Attributes: attributeValues,
		Timeout:    auct.Timeout,
		Timestamp:  time.Now(),
	}

	// Get bidder simulators instead of direct bidders
	simulators := o.bidderManager.GetBidderSimulators()

	for _, simulator := range simulators {
		wg.Add(1)

		go func(sim *bidder.Simulator) {
			defer wg.Done()

			select {
			case o.semaphore <- struct{}{}:
				defer func() { <-o.semaphore }()
			case <-ctx.Done():
				return
			}

			bidResponse, err := sim.EvaluateBid(ctx, bidRequest)
			if err != nil || bidResponse == nil {
				return
			}

			bid := types.Bid{
				BidderID:  bidResponse.BidderID,
				Amount:    bidResponse.Amount,
				Timestamp: bidResponse.Timestamp,
			}

			select {
			case bidCh <- bid:
			case <-ctx.Done():
			}
		}(simulator)
	}

	go func() {
		wg.Wait()
		close(bidCh)
	}()

	for bid := range bidCh {
		processor.AddBid(bid)
	}
}
