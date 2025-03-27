package comms

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Master Connection support 2-way communication
type MasterConnection struct {
	addr     NodeAddr
	sendAddr NodeAddr

	id     int
	header string

	send    chan *Message
	receive chan *Message

	heartbeatInterval time.Duration

	logger *log.Logger
}

func NewMasterConnection(src, dest NodeAddr, suid string, hbInterval time.Duration, l *log.Logger) *MasterConnection {
	i_suid, err := strconv.Atoi(suid)
	if err != nil {
		l.Fatalln("Failed to convert suid to int")
		return nil
	}

	return &MasterConnection{
		addr:              src,
		sendAddr:          dest,
		id:                i_suid,
		header:            CompileHeader(src.String(), suid, dest.String()),
		send:              make(chan *Message, 10),
		receive:           make(chan *Message, 10),
		heartbeatInterval: hbInterval,
		logger:            l,
	}
}

func (mc MasterConnection) GetID() int {
	return mc.id
}

func (mc MasterConnection) PrepareMsg(p *Payload) *Message {
	return &Message{
		source:      mc.addr,
		suid:        os.Getpid(),
		destination: mc.sendAddr,
		payload:     p,
	}
}

// func SendCMD
// func SendData
// func SendDef

func (mc *MasterConnection) SendPayload(p *Payload) {
	mc.send <- mc.PrepareMsg(p)
}

func (mc *MasterConnection) SendMsg(msg *Message) {
	mc.send <- msg
}

func (mc *MasterConnection) StartMasterConnectionLoop(c chan *Payload) {

	go mc.sendLoop()
	go mc.receiveDecoder(c)
	go mc.sendHeartbeat()

	// mc.listen()
}

func (mc MasterConnection) sendHeartbeat() {
	for {
		p, err := NewPayload("alive", "cmd")
		if err != nil {
			mc.logger.Panicln("Failed to create heartbeat payload")
		}
		mc.send <- mc.PrepareMsg(p)
		time.Sleep(5 * time.Second)
	}
}

func (mc *MasterConnection) sendLoop() {

	for msg := range mc.send {

		fmt.Println("SENDING")

		if strings.Compare(mc.sendAddr.String(), "") == 0 {
			mc.logger.Fatal("Destination is not set")
		}

		conn, err := net.Dial("tcp", mc.sendAddr.String())
		if err != nil {
			mc.logger.Fatal("Failed to connect to "+mc.sendAddr.String(), "\n", err)
		}

		compiledMsg, err := msg.Compile()
		if err != nil {
			mc.logger.Fatal("Failed to compile message:", err)
		}
		_, err = conn.Write([]byte(compiledMsg))
		if err != nil {
			mc.logger.Fatal("Failed to write")
		}
		conn.Close()
	}
}

func (mc *MasterConnection) receiveDecoder(c chan<- *Payload) {
	for msg := range mc.receive {

		if msg.payload.ptype == network {
			mc.logger.Println("Received network related Message:", msg.payload.ReadData())

			// TODO handle network changes
			// -----------------------------------

			// -----------------------------------

		} else {
			c <- msg.payload
		}
	}
}

func (mc *MasterConnection) Listen() {
	ln, err := net.Listen("tcp", mc.addr.String())
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	mc.logger.Println("Listening on ", mc.addr.String())
	for {

		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go mc.handleConnection(conn)
	}
}

func (mc *MasterConnection) handleConnection(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	//should make loop? to read multiple messages?

	_, err := conn.Read(buf)
	if err != nil {
		mc.logger.Println("Error reading:", err.Error())
	}
	msg, err := ParseMessage(string(buf))
	if err != nil {
		mc.logger.Fatal(err)
	}
	mc.receive <- msg
	conn.Close()
}
