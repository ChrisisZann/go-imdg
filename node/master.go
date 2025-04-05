package node

import (
	"context"
	"fmt"
	"go-imdg/comms"
	"go-imdg/config"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Master struct {
	id int

	// Comms
	config.Node
	comms.NetworkReader
	Receiver chan *comms.Message

	// Topology
	// TODO I need to change the map to hold keys of netSlave
	slaveTopology map[int]*netSlave
	topologyLock  sync.RWMutex

	// Context params
	ctx    context.Context
	cancel func()
	wg     sync.WaitGroup
}

// Instantiate new master node
func NewMaster(cfg config.Node) *Master {

	cfg.Logger.Println("Setting up new master...")
	cfg.Logger.Println("PID:", os.Getpid())
	cfg.Logger.Println("CFG:", cfg)

	newAddr, err := comms.NewNodeAddr("tcp", cfg.Hostname+":"+cfg.LPort)
	if err != nil {
		cfg.Logger.Fatal("Failed to create NodeAddr for master:", err)
	}

	m, cf := context.WithCancel(context.Background())

	return &Master{
		id:   os.Getpid(),
		Node: cfg,
		NetworkReader: *comms.NewNetworkReader(
			newAddr,
			cfg.RxLogger,
			10*time.Second,
		),
		Receiver:      make(chan *comms.Message, 10),
		slaveTopology: make(map[int]*netSlave),
		ctx:           m,
		cancel:        cf,
	}
}

// Start master node go-routines
func (m *Master) Start() {

	m.initMasterCommands()

	m.wg.Add(4)
	go func() {
		defer m.wg.Done()
		m.ReceiveHandler()
	}()
	go func() {
		defer m.wg.Done()
		m.checkHeartbeatLoop()
	}()
	go func() {
		defer m.wg.Done()
		m.userInput()
	}()
	go func() {
		defer m.wg.Done()
		m.Listen(m.ctx, m.Receiver)
	}()

	m.wg.Wait()
}

// Stops gracefully master node
func (m *Master) Stop() {
	m.Logger.Println("Stopping Master...")

	// Cancel the context to signal all goroutines to stop
	m.cancel()

	m.topologyLock.Lock()
	for _, slave := range m.slaveTopology {
		slave.connection.CloseConn()
	}
	m.topologyLock.Unlock()

	// Wait for all goroutines to finish
	m.wg.Wait()

	// Close the Receiver channel to unblock ReceiveHandler
	close(m.Receiver)

	m.Logger.Println("Master successfully shut down")
}

// Broadcast payload to all slaves
func (m *Master) BroadcastToSlaves(p *comms.Payload) {
	for i := range m.slaveTopology {
		m.slaveTopology[i].connection.SendPayload(p)
	}
}

// Decode incoming messages
func (m *Master) ReceiveHandler() {
	for {
		select {
		case msg, ok := <-m.Receiver:
			// m.RxLogger.Println("Received:", msg)
			if !ok {
				m.Logger.Println("master receive channel closed, exiting handler...")
				return
			}

			// m.Logger.Printf("Received Message: <%s>\n", msg)
			// m.Logger.Println("Sender:", msg.ReadSenderID())
			if msg.ReadSenderID() == 0 {
				// Internal message
				if strings.Compare(msg.ReadPayloadData(), "stop") == 0 {
					go m.Stop()

					//i think i need this return, because stop will return here and miss the signal??
					// return
				}
			} else if !m.exists_slave(msg.ReadSenderID()) {
				m.Logger.Println("Received message from new slave:", msg.ReadDest())
				// m.Logger.Println("I am:", m.Hostname+":"+m.LPort)
				m.addSlave(msg.ReadSender(), msg.ReadSenderID())
			} else {
				// m.Logger.Println("Received Payload Type:", msg.ReadPayloadType())
				// m.Logger.Println("Received Payload Data:", msg.ReadPayloadData())

				// Message Decoder
				if msg.GetPayloadType() == comms.StringToPayloadType("cmd") {
					if strings.Compare(msg.ReadPayloadData(), "alive") == 0 {
						m.Logger.Println("Received heartbeat from", msg.ReadSenderID())
						m.updateHeartbeat(msg.ReadSenderID(), true)
					}
				}
			}
		case <-m.ctx.Done():
			m.Logger.Println("ctx cancelled : stopping ReceiveHandler", m.ctx.Err())

			// Handle context cancellation gracefully
			//if err := m.ctx.Err(); err != nil {
			//	if err == context.Canceled {
			//		m.Logger.Println("ReceiveHandler: Context canceled, shutting down...")
			//	} else {
			//		m.Logger.Println("ReceiveHandler: Unexpected context error:", err)
			//	}
			//}
			return
		}
	}
}

// Receive direct messages from stdin
func (m *Master) userInput() {
	inputChan := make(chan string)
	go func() {
		var userInput string
		for {
			fmt.Print("Enter command:")
			fmt.Scan(&userInput)
			inputChan <- userInput

			// TBD : if i really need this check, why i need the goroutine??
			if strings.Compare(userInput, "stop") == 0 {
				break
			}
		}
	}()

	for {
		select {
		case <-m.ctx.Done():
			m.Logger.Println("ctx cancelled : stopping userInput", m.ctx.Err())
			close(inputChan)

			return
		case userInput := <-inputChan:
			p, err := comms.NewPayload(userInput, "cmd")
			if err != nil {
				m.Logger.Println("error - cant create payload from user")
				m.Logger.Println("Closing user input routine...")
				return
			}

			msg := comms.NewMessage(
				comms.NodeAddr{},
				0,
				m.GetAddr(),
				p,
			)
			m.Receiver <- msg
		}
	}
}

//========================================================================================================================

// internal slave representation for master node
type netSlave struct {
	connection *comms.NetworkWriter
	*heartbeat
}

// Returns true if slave already exists in topology
func (m *Master) exists_slave(cmpID int) bool {
	_, ok := m.slaveTopology[cmpID]
	return ok
}

// Adds new slave in topology and open communication
func (m *Master) addSlave(dest comms.NodeAddr, sid int) {
	m.topologyLock.Lock()

	addr, err := comms.NewNodeAddr("tcp", m.Hostname+":"+m.LPort)
	if err != nil {
		m.Logger.Fatal("Failed to create NodeAddr:", err)
	}

	tmp := netSlave{
		connection: comms.NewNetworkWriter(
			addr,
			dest,
			strconv.Itoa(m.id),
			m.Logger,
			m.TxLogger,
		),
		heartbeat: newHeartbeat(),
	}
	tmp.heartbeat.setCheckFrequency(30 * time.Second)

	m.slaveTopology[sid] = &tmp
	m.topologyLock.Unlock()

	tmp.connection.OpenSendChannel()

	p, err := comms.NewPayload("captured", "cmd")
	if err != nil {
		m.Logger.Fatal("Failed to create new payload")
	}
	m.BroadcastToSlaves(p)

	m.Logger.Println("Added slave to topology")
}
