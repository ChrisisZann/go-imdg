package comms

import (
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

type masterAddr struct {
	network string
	addr    string
}

func (m masterAddr) Network() string {
	return m.network
}
func (m masterAddr) String() string {
	return m.addr
}

type Master struct {
	slaves     map[*mSlave]bool
	broadcast  chan []byte
	register   chan *mSlave
	unregister chan *mSlave

	// comms
	addr masterAddr
}

func NewMaster(h, p string) *Master {
	return &Master{
		slaves:     make(map[*mSlave]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *mSlave),
		unregister: make(chan *mSlave),
		addr:       masterAddr{network: "tcp", addr: h + p},
	}
}

func (m Master) CreateMessage(msg string) message {
	return message{
		source:     m.addr,
		sourceUUID: 666,
		content:    msg,
		respPort:   CONN_PORT,
	}
}

func (m *Master) alreadyConnected(s mSlave) bool {
	for connections := range m.slaves {
		// if connections.addr == s.addr && connections.uuid == s.uuid {
		if connections.uuid == s.uuid {
			return true
		}
	}
	return false
}

func (m *Master) Run() {
	for {
		select {

		case slave := <-m.register:
			m.slaves[slave] = true
			// slave.send("PORT")
			// m.slaves[slave]

		case slave := <-m.unregister:
			if _, ok := m.slaves[slave]; ok {
				delete(m.slaves, slave)
				// close(slave.Send)
			}

			// case message := <-m.broadcast:
			// 	for slave := range m.slaves {
			// 		select {
			// 		case slave.Send <- message:
			// 		default:
			// 			//mSlave disconnected
			// 			close(slave.Send)
			// 			delete(m.slaves, slave)
			// 		}
			// 	}
		}
	}
}

func (m *Master) Listen() {
	ln, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
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

	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	conn.Close()

	fmt.Printf("Receiced message: len=%d ,msg=\"%s\"\n", reqLen, buf)
	msg, err := ParseMessage(string(buf))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decoded message:\n%s\n", msg)

	sAddr := strings.Split(conn.RemoteAddr().String(), ":")[0]

	slave := &mSlave{
		uuid: msg.sourceUUID,
		addr: sAddr + ":" + msg.respPort,
	}
	log.Println(slave)

	if !m.alreadyConnected(*slave) {

		m.register <- slave
		fmt.Println("New connection added")

		slave.send(m.CreateMessage("OK").CompileMessage())
	} else {
		fmt.Println("Already connected")
		slave.send(m.CreateMessage("READY").CompileMessage())
	}
	fmt.Println("===================================================")
}
