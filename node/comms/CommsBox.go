package comms

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type CommsBox struct {
	mAddr NodeAddr
	addr  NodeAddr

	header string

	Send    chan *message
	Receive chan *message

	// msg *message

	ConnControl *connControl

	logger *log.Logger
}

func NewCommsBox(master, name, h, p string, lgr *log.Logger) *CommsBox {
	return &CommsBox{
		mAddr:       NewNodeAddr("tcp", master),
		addr:        NewNodeAddr("tcp", h+":"+p),
		header:      h + ";" + p + ";" + strconv.Itoa(os.Getpid()) + ";",
		Send:        make(chan *message),
		Receive:     make(chan *message),
		ConnControl: NewConnFsm(),
		logger:      lgr,
	}
}

func (cb *CommsBox) PrepareMsg(p *Payload) *message {
	return &message{
		source:      cb.addr,
		suid:        os.Getpid(),
		destination: cb.mAddr,
		payload:     p,
	}
}

func (cb *CommsBox) SendMsg(msg *message) {
	cb.Send <- msg
}

func (cb *CommsBox) RunComms() {

	for {
		select {
		case message := <-cb.Send:
			conn, err := net.Dial("tcp", cb.mAddr.String())
			if err != nil {
				log.Fatal("Failed to connect to "+cb.mAddr.String(), err)
			}

			_, err = conn.Write([]byte(message.Compile()))
			if err != nil {
				log.Fatal("Failed to write")
			}
			conn.Close()

		case message := <-cb.Receive:
			cb.logger.Println("===================================================")
			//Dont print zeros in the message
			// finIdx := bytes.IndexByte(message, 0)
			trimmed := strings.Trim(message.payload.data, "\x00")
			cb.logger.Printf("Decoded received message: <%s>\n", trimmed)
			// cb.logger.Printf("Decoded received payload: <%s>\n", message.payload)
			// cb.logger.Printf("Decoded received payload: <%s>\n", message.payload.ptype)
			// cb.logger.Printf("Decoded received payload: <%s>\n", message.payload.data)
			// cb.logger.Printf("Decoded received payload: <%s>\n", ParseVarFSM(message.payload.data))
			// cb.logger.Printf("Decoded received payload: <%s>\n", ParseVarFSM(strings.Trim(message.payload.data, "\x00")))
			cb.logger.Println("===================================================")

			if message.payload.ptype == cmd {
				cb.ConnControl.NewEvent <- ParseVarFSM(message.payload.data)
			}
		}
	}
}

func (cb *CommsBox) Listen() {
	ln, err := net.Listen("tcp", cb.addr.String())
	if err != nil {
		panic(err)
	}
	defer ln.Close()

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

	_, err := conn.Read(buf)
	if err != nil {
		cb.logger.Println("Error reading:", err.Error())
	}
	msg, err := ParseMessage(string(buf))
	if err != nil {
		log.Fatal(err)
	}
	cb.Receive <- msg
	conn.Close()
}
