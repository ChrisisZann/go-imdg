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
				nr.logger.Println("Received network related Message:", msg.payload.ReadData())

				// TODO handle network changes
				// -----------------------------------

				// -----------------------------------

			} else {
				// Send one level up
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
		// Create a channel to signal that we should stop listening
		errChan := make(chan error, 1)

		// Attempt to accept a connection in a separate goroutine
		go func() {
			conn, err := ln.Accept()
			if err != nil {
				errChan <- err
				return
			}

			// Handle the connection
			go nr.handleConnection(ctx, conn)
		}()

		select {
		case <-ctx.Done():
			// Context cancelled: stop listening
			nr.logger.Println("Context cancelled: stopping Listen", ctx.Err())
			ln.Close() //THIS IS IMPORTANT!!
			return
		case err := <-errChan:
			// If Accept() fails due to context cancellation or other error, stop listening
			if err != nil && ctx.Err() == nil {
				// Unexpected error, panic
				panic(err)
			}
			// If error was due to context cancellation, exit the listener
			nr.logger.Println("Listener stopped due to error:", err)
			return
		}
	}
}

func (nr *NetworkReader) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	nr.logger.Println("starting new nr handler...")

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			nr.logger.Println("ctx cancelled : stopping nr handleConnection", ctx.Err())
			return

		case <-ticker.C:
			nr.logger.Println("iTS ALIVE!!!!")

		default:

			n, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					nr.logger.Println("Connection closed by client")
				} else {
					nr.logger.Println("Error reading:", err.Error())
				}
				return
			}
			data := buf[:n]

			nr.logger.Println("received in open conn:", conn.LocalAddr())

			msg, err := ParseMessage(string(data))
			if err != nil {
				nr.logger.Println("error - ", err)
			}
			nr.receive <- msg
		}
	}
}
