package utils

import (
	"math/rand"
	"time"
)

// InitRandom seeds the random number generator
func InitRandom() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomFloat returns a random float between min and max
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandomInt returns a random integer between min and max
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// RandomDuration returns a random duration between min and max
func RandomDuration(min, max time.Duration) time.Duration {
	nanos := rand.Int63n(int64(max-min) + int64(min))
	return time.Duration(nanos)
}

// RandomChance returns true with the given probability (0.0 to 1.0)
func RandomChance(probability float64) bool {
	return rand.Float64() < probability
}
