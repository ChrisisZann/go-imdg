package node

import (
	"go-imdg/comms"
	"go-imdg/config"
	"os"
	"strconv"
	"time"
)

type Slave struct {
	id int

	config.Node
	comms.MasterConnection

	Receiver chan *comms.Payload
}

func (s Slave) CompileHeader(dest string) string {
	return comms.CompileHeader(s.Hostname, strconv.Itoa(s.id), dest)
}

func (s *Slave) NewMasterConnection(dest string, destPort string) {

	s.Logger.Println("Creating new connection...")

	srcAddr, err := comms.NewNodeAddr("tcp", s.Hostname+":"+s.LPort)
	if err != nil {
		s.Logger.Println("Error creating source address:", err)
		return
	}
	desAddr, err := comms.NewNodeAddr("tcp", dest+":"+destPort)
	if err != nil {
		s.Logger.Println("Error creating destination address:", err)
		return
	}

	s.MasterConnection = *comms.NewMasterConnection(
		srcAddr,
		desAddr,
		strconv.Itoa(s.id),
		5*time.Second,
		s.Logger,
	)
}

func NewSlave(cfg config.Node) *Slave {

	cfg.Logger.Println("Setting up new slave...")
	cfg.Logger.Println("PID:", os.Getpid())
	cfg.Logger.Println("CFG:", cfg)

	return &Slave{
		id:       os.Getpid(),
		Node:     cfg,
		Receiver: make(chan *comms.Payload, 10),
	}
}

func (s *Slave) ReceiveHandler() {

	for p := range s.Receiver {
		s.Logger.Println("Received:", p)
	}
}
