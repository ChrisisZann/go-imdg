package comms

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// Slave Connection support outbound communication only!
type NetworkWriter struct {
	addr     NodeAddr
	sendAddr NodeAddr
	conn     net.Conn

	id int

	send chan *Message

	logger *log.Logger
}

func NewSlaveConnection(src, dest NodeAddr, suid string, l *log.Logger) *NetworkWriter {
	i_suid, err := strconv.Atoi(suid)
	if err != nil {
		l.Fatalln("Failed to convert suid to int")
		return nil
	}

	if strings.Compare(dest.String(), "") == 0 {
		l.Println("error - Cannot create slave conn, Destination is not set :", suid)
	}

	newConn, err := net.Dial(dest.Network(), dest.String())
	if err != nil {
		return nil
	}

	return &NetworkWriter{
		addr:     src,
		sendAddr: dest,
		conn:     newConn,
		id:       i_suid,
		send:     make(chan *Message, 10),
		logger:   l,
	}
}

func (cb NetworkWriter) GetID() int {
	return cb.id
}

func (cb NetworkWriter) PrepareMsg(p *Payload) *Message {
	return &Message{
		source:      cb.addr,
		suid:        os.Getpid(),
		destination: cb.sendAddr,
		payload:     p,
	}
}

func (cb *NetworkWriter) SendPing() {
	p, err := NewPayload("ping", "cmd")
	if err != nil {
		cb.logger.Panicln("failed to create ping payload")
	}
	cb.send <- cb.PrepareMsg(p)
}

func (cb *NetworkWriter) SendPayload(p *Payload) {
	cb.send <- cb.PrepareMsg(p)
}

func (cb *NetworkWriter) SendMsg(msg *Message) {
	cb.send <- msg
}

func (cb *NetworkWriter) OpenSendChannel() {
	go cb.sendLoop()
}

func (cb *NetworkWriter) sendLoop() {

	for msg := range cb.send {

		fmt.Println("SENDING")

		if strings.Compare(cb.sendAddr.String(), "") == 0 {
			cb.logger.Fatal("Destination is not set")
		}

		conn, err := net.Dial("tcp", cb.sendAddr.String())
		if err != nil {
			cb.logger.Fatal("Failed to connect to "+cb.sendAddr.String(), "\n", err)
		}

		compiledMsg, err := msg.Compile()
		if err != nil {
			cb.logger.Fatal("Failed to compile message: ", err)
		}
		_, err = conn.Write([]byte(compiledMsg))
		if err != nil {
			cb.logger.Fatal("Failed to write")
		}
		conn.Close()
	}
}
