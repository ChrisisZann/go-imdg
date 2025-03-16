package worker

type workerState int

const (
	stopped workerState = iota
	startup
	waiting
	process
	shutdown
	fatal
)

func (ws workerState) String() string {
	switch ws {
	case stopped:
		return "stopped"
	case startup:
		return "startup"
	case waiting:
		return "waiting"
	case shutdown:
		return "shutdown"
	case fatal:
		return "fatal"
	}
	return "error - bad workerState"
}

type varWorkerFSM int

const (
	success varWorkerFSM = iota
	failure
	request
	wait
	start
	stop
)

type workerControl struct {
	curState workerState
	nxtState chan workerState
	newEvent chan varWorkerFSM
}

func (wc workerControl) String() string {
	return wc.curState.String()
}

// NewWorkerControl creates a new workerControl struct
func NewWorkerControl() *workerControl {
	return &workerControl{
		curState: stopped,
		nxtState: make(chan workerState),
		newEvent: make(chan varWorkerFSM),
	}
}
