package comms

import (
	"fmt"
	"log"
	"net"
)

type Master struct {
	broadcast  chan []byte
	register   chan *mSlave
	unregister chan *mSlave

	// comms
	addr nodeAddr
}

func NewMaster(h, p string) *Master {
	return &Master{
		// slaves:     make(map[*mSlave]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *mSlave),
		unregister: make(chan *mSlave),
		addr:       NewNodeAddr("tcp", h+":"+p),
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

	// ===========================================================================================
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Printf("Receiced message: len=%d ,msg=\"%s\"\n", reqLen, buf)

	// ===========================================================================================
	// Process Message
	msg, err := ParseMessage(string(buf))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decoded message:\t%s\n", msg)

	conn.Close()
	fmt.Println("===================================================")

}
