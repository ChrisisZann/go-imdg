package node

import (
	"go-imdg/comms"
	"go-imdg/config"
	"go-imdg/data" // Import the data package
	"os"
	"strconv"
	"strings"
	"time"
)

type Slave struct {
	id int

	config.Node
	comms.NetworkRW

	Receiver  chan *comms.Payload
	DataStore *data.MemPage // Use MemPage as the data store
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

	newNetRW := comms.NewNetworkRW(
		srcAddr,
		desAddr,
		strconv.Itoa(s.id),
		5*time.Second,
		s.Logger,
	)
	if newNetRW == nil {
		s.Logger.Fatal("error - failed to create new NewNetworkRW")
	}

	s.NetworkRW = *newNetRW
}

func NewSlave(cfg config.Node) *Slave {

	cfg.Logger.Println("Setting up new slave...")
	cfg.Logger.Println("PID:", os.Getpid())
	cfg.Logger.Println("CFG:", cfg)

	// Initialize MemPage
	memPage := &data.MemPage{}
	memPage.Init()

	return &Slave{
		id:        os.Getpid(),
		Node:      cfg,
		Receiver:  make(chan *comms.Payload, 10),
		DataStore: memPage, // Assign MemPage to DataStore
	}
}

func (s *Slave) ReceiveHandler() {

	for p := range s.Receiver {
		s.Logger.Println("Received:", p)

		// Example: Handle a "save" command to store data in MemPage
		str := p.ReadType().String()
		cmp_str := comms.StringToPayloadType("save").String()

		if strings.Compare(str, cmp_str) == 0 {
			data := []byte(p.ReadData())
			err := s.DataStore.Save(data)
			if err != nil {
				s.Logger.Println("Error saving data:", err)
			} else {
				s.Logger.Println("Data saved successfully")
			}
		}

		cmp_str = comms.StringToPayloadType("read").String()

		// Example: Handle a "read" command to retrieve data from MemPage
		if strings.Compare(str, cmp_str) == 0 {
			pos, err := strconv.Atoi(p.ReadData())
			if err != nil {
				s.Logger.Println("Invalid position:", err)
				continue
			}
			data, err := s.DataStore.Read(pos)
			if err != nil {
				s.Logger.Println("Error reading data:", err)
				continue
			}
			s.Logger.Printf("Read data at position %d: %s\n", pos, string(data))
		}
	}
}
