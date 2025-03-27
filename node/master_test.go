// filepath: /home/yippee/Documents/fedoraWorkspace/go-imdg/node/master_test.go
package node

import (
	"go-imdg/comms"
	"go-imdg/config"
	"log"
	"os"
	"testing"
	"time"
)

func TestMasterWithTwoSlaves(t *testing.T) {
	// Setup logger
	logFile, _ := os.CreateTemp("", "master_test.log")
	defer os.Remove(logFile.Name())
	logger := log.New(logFile, "", log.LstdFlags)

	// Create master configuration
	masterCfg := config.Node{
		Logger:   logger,
		NodeType: "master",
		Hostname: "localhost",
		LPort:    "3333",
		Name:     "m1",
	}

	// Initialize master
	master := NewMaster(masterCfg)
	go master.Start()

	// Simulate two slaves
	slave1Addr, err := comms.NewNodeAddr("tcp", "localhost:3334")
	if err != nil {
		t.Fatalf("Failed to create address for slave 1: %v", err)
	}
	slave2Addr, err := comms.NewNodeAddr("tcp", "localhost:3335")
	if err != nil {
		t.Fatalf("Failed to create address for slave 2: %v", err)
	}

	// Simulate slave 1 sending a heartbeat
	slave1Payload, err := comms.NewPayload("alive", "cmd")
	if err != nil {
		t.Fatalf("Failed to create payload for slave 1: %v", err)
	}
	destinationAddr, err := comms.NewNodeAddr("tcp", "localhost:3333")
	if err != nil {
		t.Fatalf("Failed to create destination address for slave 1: %v", err)
	}
	slave1Msg := comms.NewMessage(slave1Addr, 1, destinationAddr, slave1Payload)
	master.Receiver <- slave1Msg

	// Simulate slave 2 sending a heartbeat
	slave2Payload, err := comms.NewPayload("alive", "cmd")
	if err != nil {
		t.Fatalf("Failed to create payload for slave 2: %v", err)
	}
	destinationAddr, err = comms.NewNodeAddr("tcp", "localhost:3333")
	if err != nil {
		t.Fatalf("Failed to create destination address for slave 2: %v", err)
	}
	slave2Msg := comms.NewMessage(slave2Addr, 2, destinationAddr, slave2Payload)
	master.Receiver <- slave2Msg

	// Allow some time for the master to process the messages
	time.Sleep(1 * time.Second)

	// Verify that both slaves were added to the master's topology
	if len(master.slaveTopology) != 2 {
		t.Fatalf("Expected 2 slaves in topology, but got %d", len(master.slaveTopology))
	}

	// Verify that the slaves are marked as alive
	if !master.slaveTopology[1].alive {
		t.Errorf("Slave 1 should be marked as alive")
	}
	if !master.slaveTopology[2].alive {
		t.Errorf("Slave 2 should be marked as alive")
	}
}
