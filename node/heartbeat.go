package node

import (
	"time"
)

const DEFAULT_FREQUENCY = 30 * time.Second

type heartbeat struct {
	alive          bool
	lastPing       time.Time
	checkFrequency time.Duration
}

func newHeartbeat() *heartbeat {
	return &heartbeat{
		alive:          false,
		lastPing:       time.Now(),
		checkFrequency: DEFAULT_FREQUENCY,
	}
}

func (h *heartbeat) setCheckFrequency(f time.Duration) {
	h.checkFrequency = f
}

func (m *Master) checkHeartbeatLoop() {
	ticker := time.NewTicker(15 * time.Second) // Check every 15 seconds
	defer ticker.Stop()

	for range ticker.C {
		for id, slave := range m.slaveTopology {
			if time.Since(slave.lastPing) > DEFAULT_FREQUENCY { // 30-second timeout
				m.Logger.Printf("Slave %d is unresponsive. Marking as inactive.\n", id)
				m.updateHeartbeat(id, false)
			}
		}
	}
}

func (m *Master) updateHeartbeat(senderID int, flag bool) {
	m.topologyLock.Lock()
	defer m.topologyLock.Unlock()

	if slave, ok := m.slaveTopology[senderID]; ok {
		slave.alive = flag
		slave.lastPing = time.Now()
	}
}
