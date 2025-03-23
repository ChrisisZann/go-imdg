package comms

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type SlaveConnection struct {
	addr     NodeAddr
	sendAddr NodeAddr

	id     int
	header string

	send chan *Message

	logger *log.Logger
}

func NewSlaveConnection(src, dest NodeAddr, suid string, l *log.Logger) *SlaveConnection {
	i_suid, err := strconv.Atoi(suid)
	if err != nil {
		l.Fatalln("Failed to convert suid to int")
		return nil
	}

	return &SlaveConnection{
		addr:     src,
		sendAddr: dest,
		id:       i_suid,
		header:   CompileHeader(src.String(), suid, dest.String()),
		send:     make(chan *Message, 10),
		logger:   l,
	}
}

func (cb SlaveConnection) GetID() int {
	return cb.id
}

func (cb SlaveConnection) PrepareMsg(p *Payload) *Message {
	return &Message{
		source:      cb.addr,
		suid:        os.Getpid(),
		destination: cb.sendAddr,
		payload:     p,
	}
}

func (cb *SlaveConnection) SendPayload(p *Payload) {
	cb.send <- cb.PrepareMsg(p)
}

func (cb *SlaveConnection) SendMsg(msg *Message) {
	cb.send <- msg
}

func (cb *SlaveConnection) OpenSendChannel() {
	go cb.sendLoop()
}

func (cb *SlaveConnection) sendLoop() {

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
