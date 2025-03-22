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
	comms.CommsBox
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
		id:   os.Getpid(),
		Node: cfg,
	}
}
