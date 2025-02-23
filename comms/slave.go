package comms

import "github.com/gorilla/websocket"

type Slave struct {
	master Master
	conn   *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}
