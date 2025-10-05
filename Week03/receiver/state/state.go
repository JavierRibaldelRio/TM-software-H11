package state

import "sync"

var (
	mu      sync.RWMutex
	PowerOn bool = true
)

// Set Power
func SetPower(on bool) {
	mu.Lock()
	defer mu.Unlock()
	PowerOn = on
}

// Get power
func IsPowerOn() bool {
	mu.RLock()
	defer mu.RUnlock()
	return PowerOn
}
