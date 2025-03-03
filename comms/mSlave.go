package comms

type mSlave struct {
	uuid int
	addr string
}

// func (s *mSlave) send(msg []byte) {
// 	c, err := net.Dial("tcp", s.addr)
// 	if err != nil {
// 		log.Fatal("Failed to connect to "+s.addr, err)
// 	}
// 	defer c.Close()

// 	_, err = c.Write(msg)
// 	if err != nil {
// 		log.Fatal("Failed to write to "+s.addr, err)
// 	}
// }
