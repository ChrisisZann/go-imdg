package comms

import (
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

func NewNetworkWriter(src, dest NodeAddr, suid string, l *log.Logger) *NetworkWriter {
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

	nw := &NetworkWriter{
		addr:     src,
		sendAddr: dest,
		conn:     newConn,
		id:       i_suid,
		send:     make(chan *Message, 10),
		logger:   l,
	}

	return nw
}

func (nw NetworkWriter) GetID() int {
	return nw.id
}

func (nw NetworkWriter) PrepareMsg(p *Payload) *Message {
	return &Message{
		source:      nw.addr,
		suid:        os.Getpid(),
		destination: nw.sendAddr,
		payload:     p,
	}
}

func (nw *NetworkWriter) SendPing() {
	p, err := NewPayload("ping", "cmd")
	if err != nil {
		nw.logger.Panicln("failed to create ping payload")
	}
	nw.send <- nw.PrepareMsg(p)
}

func (nw *NetworkWriter) SendPayload(p *Payload) {
	nw.send <- nw.PrepareMsg(p)
}

func (nw *NetworkWriter) SendMsg(msg *Message) {
	nw.send <- msg
}

func (nw *NetworkWriter) OpenSendChannel() {
	go nw.sendLoop()
}

func (nw *NetworkWriter) sendLoop() {

	for msg := range nw.send {

		// fmt.Println("SENDING")

		if strings.Compare(nw.sendAddr.String(), "") == 0 {
			nw.logger.Fatal("Destination is not set")
		}

		conn, err := net.Dial("tcp", nw.sendAddr.String())
		if err != nil {
			nw.logger.Fatal("Failed to connect to "+nw.sendAddr.String(), "\n", err)
		}

		compiledMsg, err := msg.Compile()
		if err != nil {
			nw.logger.Fatal("Failed to compile message: ", err)
		}
		_, err = conn.Write([]byte(compiledMsg))
		if err != nil {
			nw.logger.Fatal("Failed to write")
		}
		conn.Close()
	}
}

func (nw *NetworkWriter) CloseConn() {

	err := nw.conn.Close()
	if err != nil {
		nw.logger.Println("error - Failed to close conn", err)
	}

}
