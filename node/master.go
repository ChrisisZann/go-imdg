package node

import (
	"go-imdg/comms"
	"go-imdg/config"
	"os"
	"strconv"
	"strings"
)

type Master struct {
	id int

	config.Node
	comms.MasterListener
	Receiver chan *comms.Message

	slaveTopology map[int]*netSlave
}

type netSlave struct {
	connection *comms.SlaveConnection
	*heartbeat
}

func (m Master) exists_slave(cmpID int) bool {
	_, ok := m.slaveTopology[cmpID]
	if ok {
		return true
	}
	return false
}

func (m *Master) addSlave(dest comms.NodeAddr, sid int) {

	tmp := netSlave{
		connection: comms.NewSlaveConnection(
			comms.NewNodeAddr("tcp", m.Hostname+":"+m.LPort),
			dest,
			strconv.Itoa(m.id),
			m.Logger,
		),
		heartbeat: newHeartbeat()}

	m.Logger.Println("Adding slave to topology")
	m.slaveTopology[sid] = &tmp
	m.slaveTopology[sid].connection.OpenSendChannel()
}

func (m Master) CompileHeader(dest string) string {
	return comms.CompileHeader(m.Hostname, strconv.Itoa(m.id), dest)
}

func NewMaster(cfg config.Node) *Master {

	cfg.Logger.Println("Setting up new master...")
	cfg.Logger.Println("PID:", os.Getpid())
	cfg.Logger.Println("CFG:", cfg)

	return &Master{
		id:   os.Getpid(),
		Node: cfg,
		MasterListener: *comms.NewMasterListener(
			comms.NewNodeAddr("tcp", cfg.Hostname+":"+cfg.LPort),
			cfg.Logger,
		),
		Receiver:      make(chan *comms.Message, 10),
		slaveTopology: make(map[int]*netSlave),
	}
}

func (m *Master) Start() {

	go m.ReceiveHandler()
	go m.checkHeartbeatLoop()
	m.Listen(m.Receiver)

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
					m.slaveTopology[msg.ReadSenderID()].alive = true
				}
			}
		}
	}
}
