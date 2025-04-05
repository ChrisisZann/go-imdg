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

func (nr *NetworkReader) receiveDecoder(ctx context.Context, c chan<- *Message) {

	for {
		select {
		case <-ctx.Done():
			nr.logger.Println("ctx cancelled : stopping nr receiveDecoder", ctx.Err())
			return

		case msg := <-nr.receive:
			if msg.payload.ptype == network {
				nr.logger.Println("Decoder: received network related msg:", msg.payload.ReadData())

				// TODO handle network changes
				// -----------------------------------

				// -----------------------------------

			} else {
				// Send one level up
				nr.logger.Println("Decoder : received message:", msg.String())
				c <- msg
			}
		}
	}

}

func (nr *NetworkReader) Listen(ctx context.Context, receiveChannel chan *Message) {
	ln, err := net.Listen("tcp", nr.addr.String())
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	// Start the message listener internal decoder
	go nr.receiveDecoder(ctx, receiveChannel)

	// Log the listener start
	nr.logger.Println("Listening on ", nr.addr.String())

	for {

		select {
		case <-ctx.Done():
			// Context cancelled: stop listening
			nr.logger.Println("Context cancelled: stopping Listen", ctx.Err())
			ln.Close() //THIS IS IMPORTANT!!
			return
		default:
			conn, err := ln.Accept()
			if err != nil {
				nr.logger.Println("Error accepting connection:", err)
				return
			}

			// Handle the connection
			go nr.handleConnection(ctx, conn)
		}
	}
}

func (nr *NetworkReader) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	nr.logger.Println("starting new nr handler...")

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	for {

		select {
		case <-ctx.Done():
			nr.logger.Println("ctx cancelled : stopping nrw handleConnection", ctx.Err())
			return
		default:

			n, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					nr.logger.Println("Connection closed by master")
				} else {
					nr.logger.Println("Error reading:", err.Error())
				}
				return
			}
			data := buf[:n]
			msg, err := ParseMessage(string(data))
			if err != nil {
				nr.logger.Fatal(err)
			}
			nr.receive <- msg
		}
	}
}
