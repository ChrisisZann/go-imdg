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
	comms.CommsBox

	Receiver chan *comms.Payload
}

func (m Master) CompileHeader(dest string) string {
	return comms.CompileHeader(m.Hostname, strconv.Itoa(m.id), dest)
}

func (m *Master) NewCB(dest string, destPort string) {

	m.CommsBox = *comms.NewCommsBox(
		comms.NewNodeAddr("tcp", m.Hostname+":"+m.LPort),
		comms.NewNodeAddr("tcp", dest+":"+destPort),
		strconv.Itoa(m.id),
		m.Logger,
	)
}

func NewMaster(cfg config.Node) *Master {

	cfg.Logger.Println("Setting up new master...")
	cfg.Logger.Println("PID:", os.Getpid())
	cfg.Logger.Println("CFG:", cfg)

	return &Master{
		id:       os.Getpid(),
		Node:     cfg,
		Receiver: make(chan *comms.Payload, 10),
	}
}

func (m *Master) ReceiveHandler() {

	for p := range m.Receiver {
		trimmed := strings.Trim(p.ReadData(), "\x00")
		m.Logger.Printf("Received message: <%s>\n", trimmed)
		m.Logger.Println("Received:", p)
	}
}
