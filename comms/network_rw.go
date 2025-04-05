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

func NewNetworkRW(src, dest NodeAddr, suid string, hbInterval time.Duration, l *log.Logger) *NetworkRW {
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
		l.Println("error - ", err)
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

func (nrw NetworkRW) GetID() int {
	return nrw.id
}

func (nrw NetworkRW) GetAddr() NodeAddr {
	return nrw.addr
}

func (nrw NetworkRW) PrepareMsg(p *Payload) *Message {
	return &Message{
		source:      nrw.addr,
		suid:        os.Getpid(),
		destination: nrw.sendAddr,
		payload:     p,
	}
}

// func SendCMD
// func SendData
// func SendDef

func (nrw *NetworkRW) SendPayload(p *Payload) {
	nrw.send <- nrw.PrepareMsg(p)
}

func (nrw *NetworkRW) SendMsg(msg *Message) {
	nrw.send <- msg
}

func (nrw *NetworkRW) StartMasterConnectionLoop(c chan *Message) {

	go nrw.sendLoop()
	go nrw.sendHeartbeat()

	// nrw.listen()
}

func (nrw NetworkRW) sendHeartbeat() {
	for {
		p, err := NewPayload("alive", "cmd")
		if err != nil {
			nrw.logger.Panicln("Failed to create heartbeat payload")
		}
		nrw.send <- nrw.PrepareMsg(p)
		time.Sleep(5 * time.Second)
	}
}

func (nrw *NetworkRW) sendLoop() {

	for msg := range nrw.send {

		fmt.Println("SENDING")

		compiledMsg, err := msg.Compile()
		if err != nil {
			nrw.logger.Fatal("Failed to compile message:", err)
		}

		_, err = nrw.conn.Write([]byte(compiledMsg))
		if err != nil {
			nrw.logger.Fatal("Failed to write")
		}
	}
}

func (nrw *NetworkRW) receiveDecoder(ctx context.Context, c chan<- *Message) {

	for {
		select {
		case <-ctx.Done():
			nrw.logger.Println("ctx cancelled : stopping nrw receiveDecoder", ctx.Err())
			return
		case msg := <-nrw.receive:
			if msg.payload.ptype == network {
				nrw.logger.Println("Received network related Message:", msg.payload.ReadData())

				// TODO handle network changes
				// -----------------------------------

				// -----------------------------------

			} else {
				c <- msg
			}
		}
	}
}

func (nrw *NetworkRW) Listen(ctx context.Context, receiveChan chan *Message) {
	ln, err := net.Listen("tcp", nrw.addr.String())
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	go nrw.receiveDecoder(ctx, receiveChan)

	// Listen on Network
	nrw.logger.Println("Listening on ", nrw.addr.String())

	// Create a channel to signal that we should stop listening
	errChan := make(chan error, 1)

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			errChan <- err
			return
		}
		go nrw.handleConnection(ctx, conn)
	}()

	for {
		select {
		case <-ctx.Done():
			nrw.logger.Println("ctx cancelled : stopping nrw Listen", ctx.Err())
			ln.Close()
			nrw.logger.Println("Listener closed")

			return
		case err := <-errChan:
			// If Accept() fails due to context cancellation or other error, stop listening
			if err != nil && ctx.Err() == nil {
				// Unexpected error, panic
				panic(err)
			}
			// If error was due to context cancellation, exit the listener
			nrw.logger.Println("Listener stopped due to error:", err)
			return
		}
	}
}

func (nrw *NetworkRW) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	nrw.logger.Println("starting new network reader...")

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	for {

		select {
		case <-ctx.Done():
			nrw.logger.Println("ctx cancelled : stopping nrw handleConnection", ctx.Err())
			return
		default:

			n, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					nrw.logger.Println("Connection closed by master")
				} else {
					nrw.logger.Println("Error reading:", err.Error())
				}
				return
			}
			data := buf[:n]
			msg, err := ParseMessage(string(data))
			if err != nil {
				nrw.logger.Fatal(err)
			}
			nrw.receive <- msg

		}
	}
}
