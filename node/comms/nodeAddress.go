package comms

// NodeAddr implements net.Addr interface
type NodeAddr struct {
	network string
	addr    string
}

func (m NodeAddr) Network() string {
	return m.network
}

func (m NodeAddr) String() string {
	return m.addr
}

// NodeAddr implements net.Addr interface
func NewNodeAddr(n, a string) NodeAddr {
	return NodeAddr{
		network: n,
		addr:    a,
	}
}
