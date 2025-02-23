package comms

import "net"

type Slave struct {
	master *Master
	conn   *net.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func Connect(m *Master) {
	c, err := net.Dial("tcp", "golang.org:80")
	if err != nil {
		// handle error
	}

	slave := &Slave{
		master: m,
		conn:   &c,
		send:   make(chan []byte, 256),
	}
	slave.master.register <- slave
}
