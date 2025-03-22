package node

import (
	"go-imdg/comms"
	"go-imdg/config"
	"os"
	"strconv"
)

type Slave struct {
	id int

	config.Node
	comms.CommsBox
}

func (s Slave) CompileHeader(dest string) string {
	return comms.CompileHeader(s.Hostname, strconv.Itoa(s.id), dest)
}

func (s *Slave) NewCB(dest string, destPort string) {
	s.Logger.Println("Adding new connection...")
	s.CommsBox = *comms.NewCommsBox(
		comms.NewNodeAddr("tcp", s.Hostname+":"+s.LPort),
		comms.NewNodeAddr("tcp", dest+":"+destPort),
		strconv.Itoa(s.id),
		s.Logger,
	)
}

func NewSlave(cfg config.Node) *Slave {

	cfg.Logger.Println("Setting up new slave...")
	cfg.Logger.Println("PID:", os.Getpid())
	cfg.Logger.Println("CFG:", cfg)

	return &Slave{
		id:   os.Getpid(),
		Node: cfg,
	}
}
