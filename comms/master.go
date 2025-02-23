package comms

type Master struct {
	slaves     map[*Slave]bool
	broadcast  chan []byte
	register   chan *Slave
	unregister chan *Slave
}

func NewMaster() *Master {
	return &Master{
		slaves:     make(map[*Slave]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Slave),
		unregister: make(chan *Slave),
	}
}

func (m *Master) Run() {
	for {
		select {

		case slave := <-m.register:
			m.slaves[slave] = true

		case slave := <-m.unregister:
			if _, ok := m.slaves[slave]; ok {
				delete(m.slaves, slave)
				close(slave.send)
			}

		case message := <-m.broadcast:
			for slave := range m.slaves {
				select {
				case slave.send <- message:
				default:
					//Slave disconnected
					close(slave.send)
					delete(m.slaves, slave)
				}
			}
		}
	}
}
