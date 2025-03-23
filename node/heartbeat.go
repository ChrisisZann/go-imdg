package node

import "time"

type heartbeat struct {
	alive          bool
	pingInterval   time.Duration
	reviveAttempts int
}

func newHeartbeat() *heartbeat {
	return &heartbeat{
		alive:          false,
		pingInterval:   5 * time.Second,
		reviveAttempts: 3,
	}
}
