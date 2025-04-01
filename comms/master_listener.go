package comms

import (
	"log"
	"net"
	"time"
)

const DEFAULT_LISTENER_HANDLER_TIMEOUT = 10 * time.Second

type MasterListener struct {
	addr NodeAddr

	receive chan *Message

	logger *log.Logger

	handlerTimeout time.Duration
}

func NewMasterListener(src NodeAddr, l *log.Logger, t time.Duration) *MasterListener {
	return &MasterListener{
		addr:           src,
		receive:        make(chan *Message, 10),
		logger:         l,
		handlerTimeout: t,
	}
}

func (ml *MasterListener) receiveDecoder(c chan<- *Message) {
	for msg := range ml.receive {

		if msg.payload.ptype == network {
			ml.logger.Println("Received network related Message:", msg.payload.ReadData())

			// TODO handle network changes
			// -----------------------------------

			// -----------------------------------

		} else {
			// Send one level up
			c <- msg
		}
	}
}

func (ml *MasterListener) Listen(receiveChannel chan *Message) {
	ln, err := net.Listen("tcp", ml.addr.String())
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	// Start message Listener decoder
	go ml.receiveDecoder(receiveChannel)

	// Listen on network
	ml.logger.Println("Listening on ", ml.addr.String())
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go ml.handleConnection(conn)
	}
}

func (ml *MasterListener) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	conn.SetDeadline(time.Now().Add(ml.handlerTimeout))

	_, err := conn.Read(buf)
	if err != nil {
		ml.logger.Println("Error reading:", err.Error())
	}
	msg, err := ParseMessage(string(buf))
	if err != nil {
		ml.logger.Fatal(err)
	}
	ml.receive <- msg
}
