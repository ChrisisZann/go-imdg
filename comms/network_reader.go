package comms

import (
	"context"
	"io"
	"log"
	"net"
	"time"
)

const DEFAULT_LISTENER_HANDLER_TIMEOUT = 10 * time.Second

type NetworkReader struct {
	addr NodeAddr

	receive chan *Message

	logger *log.Logger

	handlerTimeout time.Duration
}

func NewNetworkReader(src NodeAddr, l *log.Logger, t time.Duration) *NetworkReader {
	return &NetworkReader{
		addr:           src,
		receive:        make(chan *Message, 10),
		logger:         l,
		handlerTimeout: t,
	}
}

func (nr NetworkReader) GetAddr() NodeAddr {
	return nr.addr
}

func (ml *NetworkReader) receiveDecoder(c chan<- *Message) {
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

func (ml *NetworkReader) Listen(ctx context.Context, receiveChannel chan *Message) {
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

		go ml.handleConnection(ctx, conn)
	}
}

func (ml *NetworkReader) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	ml.logger.Println("starting new handler...")

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			ml.logger.Println("context timeout", ctx.Err())
			ml.logger.Println("stopping handler...")
			return
		case <-ticker.C:
			ml.logger.Println("iTS ALIVE!!!!")

		default:

			n, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					ml.logger.Println("Connection closed by client")
				} else {
					ml.logger.Println("Error reading:", err.Error())
				}
				return
			}
			data := buf[:n]

			ml.logger.Println("received in open conn:", conn.LocalAddr())

			msg, err := ParseMessage(string(data))
			if err != nil {
				ml.logger.Println("error - ", err)
			}
			ml.receive <- msg
		}
	}
}
