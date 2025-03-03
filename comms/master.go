package comms

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Master struct {
	slaves     map[*mSlave]bool
	broadcast  chan string
	register   chan *mSlave
	unregister chan *mSlave
	directMsg  chan *message

	// comms
	addr nodeAddr
}

func NewMaster(h, p string) *Master {
	return &Master{
		slaves:     make(map[*mSlave]bool),
		broadcast:  make(chan string),
		register:   make(chan *mSlave),
		unregister: make(chan *mSlave),
		directMsg:  make(chan *message),
		addr:       NewNodeAddr("tcp", h+":"+p),
	}
}

func (m Master) checkConnected(testUID int) bool {
	for connection := range m.slaves {
		if connection.uid == testUID {
			return true
		}
	}
	return false
}

func (m Master) findConnected(testUID int) *mSlave {
	for connection := range m.slaves {
		if connection.uid == testUID {
			return connection
		}
	}
	fmt.Println("Didnt find slave connection")
	return nil
}

func (m *Master) PrepareMsg(p *Payload, ms *mSlave) *message {
	return &message{
		source:      m.addr,
		suid:        os.Getpid(),
		destination: ms.addr,
		payload:     p,
	}
}

func (m *Master) Start() {
	go m.RunComms()
	m.Listen()
}

func (m *Master) RunComms() {
	for {
		select {

		case slave := <-m.register:

			m.slaves[slave] = true
			fmt.Println("Saved Slave: ", m.slaves[slave])
			// slave.send("PORT")
			// m.slaves[slave]

		case slave := <-m.unregister:
			fmt.Println("Deleting slave: ", slave)
			fmt.Printf("Slaves count before:%d\n", len(m.slaves))

			if _, ok := m.slaves[slave]; ok {
				fmt.Println("i enter the if")
				delete(m.slaves, slave)
				// close(slave.Send)
			}
			fmt.Printf("Slaves count after:%d\n", len(m.slaves))

		case message := <-m.directMsg:
			fmt.Println("received DM:", message)

		case message := <-m.broadcast:
			for s := range m.slaves {
				s.send([]byte(s.PrepareMsg(NewPayload(message, cmd)).Compile()))
			}
		}
	}
}

func (m *Master) Listen() {

	ln, err := net.Listen(m.addr.Network(), m.addr.String())
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("Listening on " + m.addr.String())
	fmt.Println("===================================================")

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go m.handleConnection(conn)
	}
}

func (m *Master) handleConnection(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	// ------------------------------------------------
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	conn.Close()

	// DEBUGGING
	// fmt.Printf("Receiced message: len=%d ,msg=\"%s\"\n", reqLen, buf)

	// ------------------------------------------------
	// Process Message
	msg, err := ParseMessage(string(buf))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decoded message:\t%s\n", msg)

	// Check destination
	if strings.Compare(msg.destination.String(), m.addr.String()) != 0 {
		fmt.Println("Error wrong destination:", msg.destination.String())
	}

	// Check
	if oldSlave := m.findConnected(msg.suid); oldSlave == nil {
		slave := &mSlave{
			uid:   msg.suid,
			addr:  msg.source,
			maddr: m.addr,
			fsm:   NewConnFsm(),
		}
		slave.Start(m.register, m.unregister, m.directMsg)
		slave.fsm.newEvent <- open
		m.broadcast <- "New Slave Added "
		// m.broadcast <- "NewSlaveAdded" + slave.String()

	} else {
		fmt.Println("Already connected")
		oldSlave.fsm.newEvent <- wait
	}

	fmt.Println("===================================================")
}
