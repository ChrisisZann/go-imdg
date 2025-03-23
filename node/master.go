package node

import (
	"go-imdg/comms"
	"go-imdg/config"
	"os"
	"strconv"
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
	heartbeat
}

func (m *Master) addSlave(dest string, destPort string) {

	tmp := netSlave{
		connection: comms.NewSlaveConnection(
			comms.NewNodeAddr("tcp", m.Hostname+":"+m.LPort),
			comms.NewNodeAddr("tcp", dest+":"+destPort),
			strconv.Itoa(m.id),
			m.Logger,
		),
		heartbeat: heartbeat{
			alive: true,
		}}

	m.Logger.Println("Adding slave to topology")
	m.slaveTopology[tmp.connection.GetID()] = &tmp
	m.slaveTopology[tmp.connection.GetID()].connection.OpenSendChannel()
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
	m.Listen(m.Receiver)

}

func (m *Master) ReceiveHandler() {

	for msg := range m.Receiver {
		m.Logger.Printf("Received Message: <%s>\n", msg)
		m.Logger.Println("Received Payload Type:", msg.ReadPayloadType())
		m.Logger.Println("Received Payload Data:", msg.ReadPayloadData())
	}
}
