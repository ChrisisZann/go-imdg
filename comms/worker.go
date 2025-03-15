package comms

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

type Worker struct {
	mAddr nodeAddr
	addr  nodeAddr

	header string

	Send    chan *message
	Receive chan *message

	msg *message

	fsm *connControl
}

func NewWorker(master, name, h, p string) *Worker {
	return &Worker{
		mAddr:   NewNodeAddr("tcp", master),
		addr:    NewNodeAddr("tcp", h+":"+p),
		header:  h + ";" + p + ";" + strconv.Itoa(os.Getpid()) + ";",
		Send:    make(chan *message),
		Receive: make(chan *message),
	}
}

func (s *Worker) PrepareMsg(inputS string) {
	s.msg = &message{
		source:      s.addr,
		suid:        os.Getpid(),
		destination: s.mAddr,
		payload:     ParsePayload(inputS),
	}
}

func (s *Worker) SendMsg() {
	s.Send <- s.msg
}

func (s *Worker) Start() {
	go s.RunComms()
	go s.Listen()
}

func (s *Worker) RunComms() {

	for {
		select {
		case message := <-s.Send:
			conn, err := net.Dial("tcp", s.mAddr.String())
			if err != nil {
				log.Fatal("Failed to connect to "+s.mAddr.String(), err)
			}

			_, err = conn.Write([]byte(message.Compile()))
			if err != nil {
				log.Fatal("Failed to write")
			}
			conn.Close()

		case message := <-s.Receive:
			fmt.Println("===================================================")
			fmt.Printf("Decoded received message:\n%s\n", *message)
			fmt.Println("===================================================")
			// s.stateMSG <- new_message
		}
	}
}

func (s *Worker) Listen() {
	ln, err := net.Listen("tcp", s.addr.String())
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go s.handleConnection(conn)
	}
}

func (s *Worker) handleConnection(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

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
