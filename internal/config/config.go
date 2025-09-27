package config

import (
	"runtime"
	"time"
)

// Simulation constants
const (
	TotalBidders         = 100
	TotalAuctions        = 40
	AttributesPerAuction = 20
	DefaultTimeout       = 2 * time.Second
)

// Resource constraints based on our specifications
const (
	DefaultMaxVCPUs            = 2
	DefaultMaxMemoryMB         = 1024 // 1GB
	MaxConcurrentBiddersPerCPU = 100
)

// ResourceLimits holds the standardized resource constraints
type ResourceLimits struct {
	MaxVCPUs             int
	MaxMemoryMB          int
	MaxConcurrentBidders int
}

// Config holds the complete simulation configuration
type Config struct {
	TotalAuctions        int
	TotalBidders         int
	AttributesPerAuction int
	AuctionTimeout       time.Duration
	ResourceLimits       ResourceLimits
}

// DefaultConfig returns the default configuration with resource standardization
func DefaultConfig() *Config {
	limits := CalculateResourceLimits()

	return &Config{
		TotalAuctions:        TotalAuctions,
		TotalBidders:         TotalBidders,
		AttributesPerAuction: AttributesPerAuction,
		AuctionTimeout:       DefaultTimeout,
		ResourceLimits:       limits,
	}
}

// CalculateResourceLimits determines optimal resource constraints
func CalculateResourceLimits() ResourceLimits {
	availableCPUs := runtime.NumCPU()
	if availableCPUs > DefaultMaxVCPUs {
		availableCPUs = DefaultMaxVCPUs
	}

	maxConcurrent := availableCPUs * MaxConcurrentBiddersPerCPU
	if maxConcurrent > TotalBidders*TotalAuctions {
		maxConcurrent = TotalBidders * TotalAuctions
	}

	return ResourceLimits{
		MaxVCPUs:             availableCPUs,
		MaxMemoryMB:          DefaultMaxMemoryMB,
		MaxConcurrentBidders: maxConcurrent,
	}
}
