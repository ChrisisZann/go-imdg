package comms

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Master Connection support 2-way communication
type NetworkRW struct {
	addr     NodeAddr
	sendAddr NodeAddr
	conn     net.Conn

	id int

	send    chan *Message
	receive chan *Message

	heartbeatInterval time.Duration

	logger *log.Logger
}

func NewMasterConnection(src, dest NodeAddr, suid string, hbInterval time.Duration, l *log.Logger) *NetworkRW {
	i_suid, err := strconv.Atoi(suid)
	if err != nil {
		l.Fatalln("Failed to convert suid to int")
		return nil
	}

	if strings.Compare(dest.String(), "") == 0 {
		l.Println("error - Cannot create master conn, Destination is not set :", suid)
	}

	newConn, err := net.Dial(dest.Network(), dest.String())
	if err != nil {
		return nil
	}

	return &NetworkRW{
		addr:              src,
		sendAddr:          dest,
		conn:              newConn,
		id:                i_suid,
		send:              make(chan *Message, 10),
		receive:           make(chan *Message, 10),
		heartbeatInterval: hbInterval,
		logger:            l,
	}
}

func (mc NetworkRW) GetID() int {
	return mc.id
}

func (mc NetworkRW) PrepareMsg(p *Payload) *Message {
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

func (mc *NetworkRW) SendPayload(p *Payload) {
	mc.send <- mc.PrepareMsg(p)
}

func (mc *NetworkRW) SendMsg(msg *Message) {
	mc.send <- msg
}

func (mc *NetworkRW) StartMasterConnectionLoop(c chan *Payload) {

	go mc.sendLoop()
	go mc.receiveDecoder(c)
	go mc.sendHeartbeat()

	// mc.listen()
}

func (mc NetworkRW) sendHeartbeat() {
	for {
		p, err := NewPayload("alive", "cmd")
		if err != nil {
			mc.logger.Panicln("Failed to create heartbeat payload")
		}
		mc.send <- mc.PrepareMsg(p)
		time.Sleep(5 * time.Second)
	}
}

func (mc *NetworkRW) sendLoop() {

	for msg := range mc.send {

		fmt.Println("SENDING")

		compiledMsg, err := msg.Compile()
		if err != nil {
			mc.logger.Fatal("Failed to compile message:", err)
		}

		_, err = mc.conn.Write([]byte(compiledMsg))
		if err != nil {
			mc.logger.Fatal("Failed to write")
		}
	}
}

func (mc *NetworkRW) receiveDecoder(c chan<- *Payload) {
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

func (mc *NetworkRW) Listen() {
	ln, err := net.Listen("tcp", mc.addr.String())
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen on Network
	mc.logger.Println("Listening on ", mc.addr.String())
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go mc.handleConnection(ctx, conn)
	}
}

func (mc *NetworkRW) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	mc.logger.Println("starting new network reader...")

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	for {

		select {
		case <-ctx.Done():
			mc.logger.Println("context cancelled this handler", ctx.Err())
			mc.logger.Println("stopping handler...")
			return
		default:

			n, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					mc.logger.Println("Connection closed by master")
				} else {
					mc.logger.Println("Error reading:", err.Error())
				}
				return
			}
			data := buf[:n]
			msg, err := ParseMessage(string(data))
			if err != nil {
				mc.logger.Fatal(err)
			}
			mc.receive <- msg

		}
	}
}
