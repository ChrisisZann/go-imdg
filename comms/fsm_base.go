package comms

type connState int

const (
	notConnected connState = iota
	validateNewConn
	connected
)

func (cs connState) String() string {
	switch cs {
	case notConnected:
		return "notConnected"
	case validateNewConn:
		return "validateNewConn"
	case connected:
		return "connected"
	}
	return "error - bad connState"
}

type varFSM int

const (
	accept varFSM = iota
	failed
	open
	close
	wait
	send
)

type connControl struct {
	curState connState
	nxtState chan connState
	newEvent chan varFSM
}

func NewConnFsm() *connControl {
	return &connControl{
		curState: notConnected,
		nxtState: make(chan connState),
		newEvent: make(chan varFSM),
	}
}
