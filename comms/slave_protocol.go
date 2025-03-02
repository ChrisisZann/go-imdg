package comms

import "log"

type protocolMSG int

const (
	portRequest protocolMSG = iota
	dataRequest
	finACK
	done
	status
	success
	failure
)

type slaveState int

const (
	initialize slaveState = iota
	startListening
	runLoop
	waiting
	fatal
)

func (s *Slave) protocolFSM() {
	currentState := s.state
	var nextState slaveState
	for {
		select {
		case event := <-s.stateMSG:
			switch currentState {

			case initialize:
				switch event {
				case success:
					nextState = startListening
				case failure:
					nextState = fatal
				}

			case startListening:
				switch event {
				case success:
					nextState = runLoop
				case failure:
					nextState = fatal
				}

			case runLoop:
				switch event {
				case success:
					nextState = waiting
				case failure:
					nextState = fatal
				}

			case fatal:
				log.Fatal("Fatal : Process failed")

			default:
				log.Println("Unknown command")
			}
		}
		s.state = nextState
	}
}
