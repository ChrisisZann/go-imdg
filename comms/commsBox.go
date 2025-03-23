package comms

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type MasterConnection struct {
	addr     NodeAddr
	sendAddr NodeAddr

	id     int
	header string

	send    chan *Message
	receive chan *Message

	logger *log.Logger
}

func NewMasterConnection(src, dest NodeAddr, suid string, l *log.Logger) *MasterConnection {
	i_suid, err := strconv.Atoi(suid)
	if err != nil {
		l.Fatalln("Failed to convert suid to int")
		return nil
	}

	return &MasterConnection{
		addr:     src,
		sendAddr: dest,
		id:       i_suid,
		header:   CompileHeader(src.String(), suid, dest.String()),
		send:     make(chan *Message, 10),
		receive:  make(chan *Message, 10),
		logger:   l,
	}
}

func (cb MasterConnection) GetID() int {
	return cb.id
}

func (cb MasterConnection) PrepareMsg(p *Payload) *Message {
	return &Message{
		source:      cb.addr,
		suid:        os.Getpid(),
		destination: cb.sendAddr,
		payload:     p,
	}
}

// func SendCMD
// func SendData
// func SendDef

func (cb *MasterConnection) SendPayload(p *Payload) {
	cb.send <- cb.PrepareMsg(p)
}

func (cb *MasterConnection) SendMsg(msg *Message) {
	cb.send <- msg
}

func (cb *MasterConnection) StartMasterConnectionLoop(c chan *Payload) {

	go cb.sendLoop()
	go cb.receiveLoop(c)
	// cb.listen()
}

func (cb *MasterConnection) sendLoop() {

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

func (cb *MasterConnection) receiveLoop(c chan<- *Payload) {
	for msg := range cb.receive {

		if msg.payload.ptype == network {
			cb.logger.Println("Received network related Message:", msg.payload.ReadData())

			// TODO handle network changes
			// -----------------------------------

			// -----------------------------------

		} else {
			c <- msg.payload
		}
	}
}

func (cb *MasterConnection) Listen() {
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

func (cb *MasterConnection) handleConnection(conn net.Conn) {
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
