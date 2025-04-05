package node

import (
	"context"
	"fmt"
	"go-imdg/comms"
	"go-imdg/config"
	"go-imdg/data" // Import the data package
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Slave struct {
	id int

	config.Node
	comms.NetworkRW

	Receiver  chan *comms.Message
	DataStore *data.MemPage // Use MemPage as the data store

	// Context params
	ctx    context.Context
	cancel func()
	wg     sync.WaitGroup
}

func (s *Slave) NewNetworkRW(dest string, destPort string) {

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
		s.TxLogger,
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

	m, cf := context.WithCancel(context.Background())

	return &Slave{
		id:        os.Getpid(),
		Node:      cfg,
		Receiver:  make(chan *comms.Message, 10),
		DataStore: memPage, // Assign MemPage to DataStore,
		ctx:       m,
		cancel:    cf,
	}
}

// Start master node go-routines
func (s *Slave) Start() {

	// s.initMasterCommands()
	s.NewNetworkRW("localhost", "3333")
	s.StartMasterConnectionLoop(s.Receiver)

	s.wg.Add(3)
	go func() {
		defer s.wg.Done()
		s.ReceiveHandler()
	}()
	go func() {
		defer s.wg.Done()
		s.userInput()
	}()
	go func() {
		defer s.wg.Done()
		s.Listen(s.ctx, s.Receiver)
	}()

	s.wg.Wait()
}

// Stops gracefully master node
func (s *Slave) Stop() {
	s.Logger.Println("Stopping Slave...")

	// Cancel the context to signal all goroutines to stop
	s.cancel()

	// Wait for all goroutines to finish
	s.wg.Wait()

	// Close the Receiver channel to unblock ReceiveHandler
	close(s.Receiver)

	s.Logger.Println("Slave successfully shut down")
}

func (s *Slave) ReceiveHandler() {

	for {
		select {
		case <-s.ctx.Done():
			s.Logger.Println("ctx cancelled : stopping ReceiveHandler", s.ctx.Err())
			//Close resources?

			return
		case msg := <-s.Receiver:
			// for p := range s.Receiver {
			s.Logger.Println("Received:", msg.String())
			str_ptype := msg.ReadPayloadType()
			str_savetype := comms.StringToPayloadType("save").String()

			// fmt.Println("str_ptype:", str_ptype)
			// fmt.Println("str_savetype:", str_savetype)
			// fmt.Println("msg.ReadSenderID()", msg.ReadSenderID())

			if msg.ReadSenderID() == 0 {
				//Internal message handling
				if strings.Compare(msg.ReadPayloadData(), "stop") == 0 {
					s.Logger.Println("Received stop command")
					go s.Stop()
				}
			} else if strings.Compare(str_ptype, str_savetype) == 0 {
				// Example: Handle a "save" command to store data in MemPage
				data := []byte(msg.ReadPayloadData())
				err := s.DataStore.Save(data)
				if err != nil {
					s.Logger.Println("Error saving data:", err)
				} else {
					s.Logger.Println("Data saved successfully")
				}
			}

			str_readtype := comms.StringToPayloadType("read").String()

			// Example: Handle a "read" command to retrieve data from MemPage
			if strings.Compare(str_ptype, str_readtype) == 0 {
				pos, err := strconv.Atoi(msg.ReadPayloadData())
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
}

// Receive direct messages from stdin
func (s *Slave) userInput() {
	inputChan := make(chan string)
	go func() {
		var userInput string
		for {
			fmt.Print("Enter command:")
			fmt.Scan(&userInput)
			inputChan <- userInput

			// TBD : if i really need this check, why i need the goroutine??
			if strings.Compare(userInput, "stop") == 0 {
				break
			}
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			s.Logger.Println("ctx cancelled : stopping userInput", s.ctx.Err())
			close(inputChan)

			return
		case userInput := <-inputChan:
			p, err := comms.NewPayload(userInput, "cmd")
			if err != nil {
				s.Logger.Println("error - cant create payload from user")
				s.Logger.Println("Closing user input routine...")
				return
			}

			msg := comms.NewMessage(
				comms.NodeAddr{},
				0,
				s.GetAddr(),
				p,
			)
			s.Receiver <- msg
		}
	}
}
