package comms

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

type Slave struct {
	// master *Master
	conn  net.Conn
	mAddr string

	listenHost string
	listenPort string

	// Buffered channel of outbound messages.
	Send       chan []byte
	Receive    chan *message
	connStatus bool

	header []byte

	state slaveState

	stateMSG chan protocolMSG
}

func (s Slave) String() string {
	return s.mAddr + "|" + strconv.FormatBool(s.connStatus) + "|" + string(s.header)
}

func NewSlave(ma, name, lh, lp string) *Slave {
	return &Slave{
		mAddr:      ma,
		listenHost: lh,
		listenPort: lp,
		connStatus: false,
		header:     []byte(lh + ";" + lp + ";" + strconv.Itoa(os.Getpid()) + ";"),
		Send:       make(chan []byte, 256),
		Receive:    make(chan *message),
		stateMSG:   make(chan protocolMSG),
	}
}

func (s *Slave) Connect() {
	c, err := net.Dial("tcp", s.mAddr)
	if err != nil {
		log.Fatal("Failed to connect to "+s.mAddr, err)
	}
	s.conn = c
}

func (s *Slave) StartFSM() {
	go s.protocolFSM()
}

func (s *Slave) Run() {

	for {
		select {
		case message := <-s.Send:
			s.Connect()
			_, err := s.conn.Write(append(s.header, message...))
			if err != nil {
				log.Fatal("Failed to write")
			}
			s.conn.Close()
		case message := <-s.Receive:
			fmt.Println("===================================================")
			fmt.Printf("Decoded received message:\n%s\n", *message)
			fmt.Println("===================================================")
		}
	}
}

func (s *Slave) Listen() {
	ln, err := net.Listen("tcp", s.listenHost+":"+s.listenPort)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	s.stateMSG <- success

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go s.handleConnection(conn)
	}
}

func (s *Slave) handleConnection(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	msg, err := ParseMessage(string(buf))
	if err != nil {
		log.Fatal(err)
	}
	s.Receive <- msg
	conn.Close()
}
