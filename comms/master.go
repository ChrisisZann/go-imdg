package comms

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Master struct {
	workers    map[*mWorker]bool
	broadcast  chan string
	register   chan *mWorker
	unregister chan *mWorker
	directMsg  chan *message

	// comms
	addr nodeAddr
}

func NewMaster(h, p string) *Master {
	return &Master{
		workers:    make(map[*mWorker]bool),
		broadcast:  make(chan string),
		register:   make(chan *mWorker),
		unregister: make(chan *mWorker),
		directMsg:  make(chan *message),
		addr:       NewNodeAddr("tcp", h+":"+p),
	}
}

func (m Master) checkConnected(testUID int) bool {
	for connection := range m.workers {
		if connection.uid == testUID {
			return true
		}
	}
	return false
}

func (m Master) findConnected(testUID int) *mWorker {
	for connection := range m.workers {
		if connection.uid == testUID {
			return connection
		}
	}
	fmt.Println("Didnt find worker connection")
	return nil
}

func (m *Master) PrepareMsg(p *Payload, ms *mWorker) *message {
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

		case worker := <-m.register:

			m.workers[worker] = true
			fmt.Println("Saved Worker: ", m.workers[worker])
			// worker.send("PORT")
			// m.workers[worker]

		case worker := <-m.unregister:
			fmt.Println("Deleting worker: ", worker)
			fmt.Printf("Workers count before:%d\n", len(m.workers))

			if _, ok := m.workers[worker]; ok {
				fmt.Println("i enter the if")
				delete(m.workers, worker)
				// close(worker.Send)
			}
			fmt.Printf("Workers count after:%d\n", len(m.workers))

		case message := <-m.directMsg:
			fmt.Println("received DM:", message)

		case message := <-m.broadcast:
			for s := range m.workers {
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
	if oldWorker := m.findConnected(msg.suid); oldWorker == nil {
		worker := &mWorker{
			uid:   msg.suid,
			addr:  msg.source,
			maddr: m.addr,
			fsm:   NewConnFsm(),
		}
		worker.Start(m.register, m.unregister, m.directMsg)
		worker.fsm.newEvent <- open
		m.broadcast <- "New Worker Added "
		// m.broadcast <- "NewWorkerAdded" + worker.String()

	} else {
		fmt.Println("Already connected")
		oldWorker.fsm.newEvent <- wait
	}

	fmt.Println("===================================================")
}
