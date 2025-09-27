package metrics

import (
	"runtime"
	"time"
)

// ResourceMonitor provides detailed resource tracking
type ResourceMonitor struct {
	startTime time.Time
	readings  []ResourceUsage
	stopChan  chan struct{}
}

// NewResourceMonitor creates a new resource monitor
func NewResourceMonitor() *ResourceMonitor {
	return &ResourceMonitor{
		readings: make([]ResourceUsage, 0),
		stopChan: make(chan struct{}),
	}
}

// Start begins resource monitoring
func (rm *ResourceMonitor) Start() {
	rm.startTime = time.Now()

	go rm.monitorLoop()
}

// Stop ends resource monitoring and returns collected data
func (rm *ResourceMonitor) Stop() []ResourceUsage {
	close(rm.stopChan)
	return rm.readings
}

// monitorLoop collects resource usage at regular intervals
func (rm *ResourceMonitor) monitorLoop() {
	ticker := time.NewTicker(50 * time.Millisecond) // High frequency for accuracy
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rm.captureReading()
		case <-rm.stopChan:
			return
		}
	}
}

// captureReading captures current resource usage
func (rm *ResourceMonitor) captureReading() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	reading := ResourceUsage{
		Timestamp:      time.Now(),
		GoroutineCount: runtime.NumGoroutine(),
		MemoryMB:       float64(m.Alloc) / 1024 / 1024,
	}

	rm.readings = append(rm.readings, reading)
}

// GetPeakUsage returns the maximum resource usage observed
func (rm *ResourceMonitor) GetPeakUsage() ResourceUsage {
	if len(rm.readings) == 0 {
		return ResourceUsage{}
	}

	peak := rm.readings[0]
	for _, reading := range rm.readings {
		if reading.MemoryMB > peak.MemoryMB {
			peak.MemoryMB = reading.MemoryMB
		}
		if reading.GoroutineCount > peak.GoroutineCount {
			peak.GoroutineCount = reading.GoroutineCount
		}
	}

	return peak
}
