package comms

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type CommsBox struct {
	addr     NodeAddr
	sendAddr NodeAddr

	header string

	send    chan *message
	receive chan *message

	logger *log.Logger
}

func NewCommsBox(src, dest NodeAddr, suid string, l *log.Logger) *CommsBox {
	return &CommsBox{
		addr:     src,
		sendAddr: dest,
		header:   CompileHeader(src.String(), suid, dest.String()),
		send:     make(chan *message, 10),
		receive:  make(chan *message, 10),
		logger:   l,
	}
}

func (cb CommsBox) PrepareMsg(p *Payload) *message {
	return &message{
		source:      cb.addr,
		suid:        os.Getpid(),
		destination: cb.sendAddr,
		payload:     p,
	}
}

// func SendCMD
// func SendData
// func SendDef

func (cb *CommsBox) SendPayload(p *Payload) {
	cb.send <- cb.PrepareMsg(p)
}

func (cb *CommsBox) SendMsg(msg *message) {
	cb.send <- msg
}

func (cb *CommsBox) StartCommsBoxLoop(c chan *Payload) {

	go cb.sendLoop()
	go cb.receiveLoop(c)
	// cb.listen()
}

func (cb *CommsBox) sendLoop() {

	for msg := range cb.send {

		fmt.Println("SENDING")

		if strings.Compare(cb.sendAddr.String(), "") == 0 {
			cb.logger.Fatal("Destination is not set")
		}

		conn, err := net.Dial("tcp", cb.sendAddr.String())
		if err != nil {
			cb.logger.Fatal("Failed to connect to "+cb.sendAddr.String(), "\n", err)
		}

		_, err = conn.Write([]byte(msg.Compile()))
		if err != nil {
			cb.logger.Fatal("Failed to write")
		}
		conn.Close()
	}
}

func (cb *CommsBox) receiveLoop(c chan<- *Payload) {
	for msg := range cb.receive {

		if msg.payload.ptype == network {
			cb.logger.Println("Received network related message:", msg.payload.ReadData())

			// TODO handle network changes
			// -----------------------------------

			// -----------------------------------

		} else {
			c <- msg.payload
		}
	}
}

func (cb *CommsBox) Listen() {
	ln, err := net.Listen("tcp", cb.addr.String())
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	cb.logger.Println("Listening on ", cb.addr.String())
	for {

		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go cb.handleConnection(conn)
	}
}

func (cb *CommsBox) handleConnection(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	//should make loop? to read multiple messages?

	_, err := conn.Read(buf)
	if err != nil {
		cb.logger.Println("Error reading:", err.Error())
	}
	msg, err := ParseMessage(string(buf))
	if err != nil {
		cb.logger.Fatal(err)
	}
	cb.receive <- msg
	conn.Close()
}
