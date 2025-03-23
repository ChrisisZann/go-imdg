package node

import "time"

type heartbeat struct {
	alive    bool
	lastPing time.Time
}

func newHeartbeat() *heartbeat {
	return &heartbeat{
		alive:    false,
		lastPing: time.Now(),
	}
}

func (m Master) checkHeartbeatLoop() {
	pingInterval := 15 * time.Second
	// reviveAttempts := 3
	for {

		for ns := range m.slaveTopology {

			if !m.slaveTopology[ns].heartbeat.alive {
				m.Logger.Println("slave is dead:", m.slaveTopology[ns])
				// Manage failure
			} else {
				m.Logger.Println("slave is alive:", m.slaveTopology[ns])
			}
		}

		time.Sleep(pingInterval)
	}
}
