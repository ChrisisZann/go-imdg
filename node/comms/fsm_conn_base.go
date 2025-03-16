package comms

import (
	"strings"
)

type connState int

const (
	disconnected connState = iota
	connecting
	listening
	invalid_state
)

func (cs connState) String() string {
	switch cs {
	case disconnected:
		return "disconnected"
	case connecting:
		return "connecting"
	case listening:
		return "listening"
	case invalid_state:
		return "invalid_state"
	}
	return "error - bad connState"
}

type VarFSM int

const (
	accept VarFSM = iota
	fatal
	open
	close
	wait
	send
)

func (vf VarFSM) String() string {
	switch vf {
	case accept:
		return "accept"
	case fatal:
		return "fatal"
	case open:
		return "open"
	case close:
		return "close"
	case wait:
		return "wait"
	case send:
		return "send"
	}
	return "error - bad VarFSM"
}

func ParseVarFSM(s string) VarFSM {

	switch strings.Trim(s, "\x00") {
	case "accept":
		return accept
	case "fatal":
		return fatal
	case "open":
		return open
	case "close":
		return close
	case "wait":
		return wait
	case "send":
		return send
	}
	return -1
}

type connControl struct {
	curState connState
	nxtState chan connState
	NewEvent chan VarFSM
}

func NewConnFsm() *connControl {
	return &connControl{
		curState: disconnected,
		nxtState: make(chan connState),
		NewEvent: make(chan VarFSM),
	}
}
