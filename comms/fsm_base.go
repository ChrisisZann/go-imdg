package comms

type connState int

const (
	disconnected connState = iota
	connecting
	listening
)

func (cs connState) String() string {
	switch cs {
	case disconnected:
		return "disconnected"
	case connecting:
		return "connecting"
	case listening:
		return "listening"
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
		curState: disconnected,
		nxtState: make(chan connState),
		newEvent: make(chan varFSM),
	}
}
