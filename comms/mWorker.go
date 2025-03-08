package comms

import (
	"log"
	"net"
	"os"
	"strconv"
)

type mWorker struct {
	uid   int
	addr  nodeAddr
	maddr nodeAddr

	fsm *connFSM
}

func (s mWorker) String() string {
	return strconv.Itoa(s.uid) + "?" + s.addr.String()
}

func (ms *mWorker) PrepareMsg(p *Payload) *message {
	return &message{
		source:      ms.maddr,
		suid:        os.Getpid(),
		destination: ms.addr,
		payload:     p,
	}
}

func (s *mWorker) send(msg []byte) {
	c, err := net.Dial("tcp", s.addr.String())
	if err != nil {
		log.Fatal("Failed to connect to "+s.addr.String(), err)
	}
	defer c.Close()

	_, err = c.Write(msg)
	if err != nil {
		log.Fatal("Failed to write to "+s.addr.String(), err)
	}
}
