package comms

type nodeAddr struct {
	network string
	addr    string
}

func (m nodeAddr) Network() string {
	return m.network
}
func (m nodeAddr) String() string {
	return m.addr
}

func NewNodeAddr(n, a string) nodeAddr {
	return nodeAddr{
		network: n,
		addr:    a,
	}
}
