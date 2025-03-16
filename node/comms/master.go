package comms

import (
	"bytes"
	"fmt"
	"go-imdg/config"
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

	// embed Node struct
	config.NodeCfg

	// comms
	addr NodeAddr
}

func NewMaster(cfg *config.NodeCfg) *Master {
	cfg.Logger.Println(cfg.Hostname + ":" + cfg.LPort)
	return &Master{
		workers:    make(map[*mWorker]bool),
		broadcast:  make(chan string),
		register:   make(chan *mWorker),
		unregister: make(chan *mWorker),
		directMsg:  make(chan *message),
		addr:       NewNodeAddr("tcp", cfg.Hostname+":"+cfg.LPort),
		NodeCfg:    *cfg,
	}
}

func (m Master) findConnected(testUID int) *mWorker {
	for connection := range m.workers {
		if connection.uid == testUID {
			return connection
		}
	}
	m.Logger.Println("ERROR - Didnt find existing worker in pool")
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
			m.Logger.Println("Worker Added in pool: ", m.workers[worker])
			// worker.send("PORT")
			// m.workers[worker]

		case worker := <-m.unregister:
			m.Logger.Println("Deleting worker: ", worker)

			if _, ok := m.workers[worker]; ok {
				m.Logger.Println("i enter the if")
				delete(m.workers, worker)
				// close(worker.Send)
			}
			m.Logger.Println("Remaining Worker count :", len(m.workers))

		case message := <-m.directMsg:
			m.Logger.Println("received DM:", message)
			if existingWorker := m.findConnected(message.suid); existingWorker == nil {
				worker := &mWorker{
					uid:    message.suid,
					addr:   message.source,
					maddr:  m.addr,
					logger: m.Logger,
					fsm:    NewConnFsm(),
				}
				worker.Start(m.register, m.unregister, m.directMsg)
				worker.fsm.NewEvent <- open
				fmt.Println(NewPayload(VarFSM(0).String(), cmd))
				worker.send([]byte(worker.PrepareMsg(NewPayload(VarFSM(0).String(), cmd)).Compile()))

				// m.broadcast <- "NewWorkerAdded" + worker.String()

			} else {
				m.Logger.Println("Already listening")
				existingWorker.fsm.NewEvent <- wait
			}

		case message := <-m.broadcast:
			m.Logger.Println("Broadcasting message: ", message)
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

	m.Logger.Println("Listening on " + m.addr.String())
	m.Logger.Println("===================================================")

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
		m.Logger.Println("Error reading:", err.Error())
	}
	conn.Close()

	// DEBUGGING
	// fmt.Printf("Receiced message: len=%d ,msg=\"%s\"\n", reqLen, buf)

	finIdx := bytes.IndexByte(buf, 0)

	m.Logger.Printf("Received message: <%s>\n", string(buf[:finIdx]))
	// ------------------------------------------------
	// Process Message
	msg, err := ParseMessage(string(buf[:finIdx]))
	if err != nil {
		log.Fatal(err)
	}
	m.Logger.Printf("Decoded message: %s\n", msg)

	// Check destination
	if strings.Compare(msg.destination.String(), m.addr.String()) != 0 {
		m.Logger.Println("Error wrong destination:", msg.destination.String())
		m.Logger.Println("master addr is:", m.addr.String())
	}

	// Check
	// TODO: Move to FSM
	m.directMsg <- msg

}
