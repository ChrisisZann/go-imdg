package node

import (
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

	config.Node
	comms.MasterListener
	Receiver chan *comms.Message

	slaveTopology map[int]*netSlave
	topologyLock  sync.RWMutex
}

type netSlave struct {
	connection *comms.SlaveConnection
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
		connection: comms.NewSlaveConnection(
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

func NewMaster(cfg config.Node) *Master {

	cfg.Logger.Println("Setting up new master...")
	cfg.Logger.Println("PID:", os.Getpid())
	cfg.Logger.Println("CFG:", cfg)

	newAddr, err := comms.NewNodeAddr("tcp", cfg.Hostname+":"+cfg.LPort)
	if err != nil {
		cfg.Logger.Fatal("Failed to create NodeAddr for master:", err)
	}

	return &Master{
		id:   os.Getpid(),
		Node: cfg,
		MasterListener: *comms.NewMasterListener(
			newAddr,
			cfg.Logger,
			10*time.Second,
		),
		Receiver:      make(chan *comms.Message, 10),
		slaveTopology: make(map[int]*netSlave),
	}
}

func (m *Master) Start() {

	m.initMasterCommands()

	go m.ReceiveHandler()
	go m.checkHeartbeatLoop()
	m.Listen(m.Receiver)
}

func (m *Master) BroadcastToSlaves(p *comms.Payload) {
	for i := range m.slaveTopology {
		m.slaveTopology[i].connection.SendPayload(p)
	}
}

func (m *Master) ReceiveHandler() {

	for msg := range m.Receiver {
		m.Logger.Printf("Received Message: <%s>\n", msg)
		m.Logger.Println("Sender:", msg.ReadSenderID())

		if !m.exists_slave(msg.ReadSenderID()) {
			m.Logger.Println("Received message from:", msg.ReadDest())
			m.Logger.Println("i am :", m.Hostname+":"+m.LPort)
			m.addSlave(msg.ReadSender(), msg.ReadSenderID())
		} else {
			m.Logger.Println("Received Payload Type:", msg.ReadPayloadType())
			m.Logger.Println("Received Payload Data:", msg.ReadPayloadData())

			// Message Decoder
			if msg.GetPayloadType() == comms.StringToPayloadType("cmd") {

				if strings.Compare(msg.ReadPayloadData(), "alive") == 0 {
					m.Logger.Println("Received heartbeat from", msg.ReadSenderID())
					m.updateHeartbeat(msg.ReadSenderID(), true)
				}
			}
		}
	}
}
