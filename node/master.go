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
	slaveTopology map[int]*netSlave
	topologyLock  sync.RWMutex

	// Context params
	ctx    context.Context
	cancel func()
	wg     sync.WaitGroup
}

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
			cfg.Logger,
			10*time.Second,
		),
		Receiver:      make(chan *comms.Message, 10),
		slaveTopology: make(map[int]*netSlave),
		ctx:           m,
		cancel:        cf,
	}
}

func (m *Master) Start() {

	m.initMasterCommands()

	m.wg.Add(3)
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

	m.Listen(m.ctx, m.Receiver)
}

func (m *Master) Stop() {
	m.Logger.Println("Stopping Master...")
	m.cancel()
	m.wg.Wait()
	m.Logger.Println("Master successfully shut down")
}

func (m *Master) BroadcastToSlaves(p *comms.Payload) {
	for i := range m.slaveTopology {
		m.slaveTopology[i].connection.SendPayload(p)
	}
}

func (m *Master) ReceiveHandler() {
	for {
		select {
		case msg := <-m.Receiver:
			m.Logger.Printf("Received Message: <%s>\n", msg)
			// m.Logger.Println("Sender:", msg.ReadSenderID())
			if msg.ReadSenderID() == 0 {
				// Internal message
				if strings.Compare(msg.ReadPayloadData(), "stop") == 0 {
					m.Stop()
				}
			} else if !m.exists_slave(msg.ReadSenderID()) {
				m.Logger.Println("Received message from:", msg.ReadDest())
				m.Logger.Println("I am:", m.Hostname+":"+m.LPort)
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
			// Handle context cancellation gracefully
			if err := m.ctx.Err(); err != nil {
				if err == context.Canceled {
					m.Logger.Println("ReceiveHandler: Context canceled, shutting down...")
				} else {
					m.Logger.Println("ReceiveHandler: Unexpected context error:", err)
				}
			}
			return
		}
	}
}

func (m *Master) userInput() {
	for {
		select {
		case <-m.ctx.Done():
			m.Logger.Println("Stopping user input routine...")
			return
		default:
			var userInput string
			fmt.Print("Enter command:")
			fmt.Scan(&userInput)

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

type netSlave struct {
	connection *comms.NetworkWriter
	*heartbeat
}

func (m *Master) exists_slave(cmpID int) bool {
	_, ok := m.slaveTopology[cmpID]
	return ok
}

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
