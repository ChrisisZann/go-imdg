package comms

import "fmt"

func (s *mWorker) Start(register, unregister chan<- *mWorker, directMsg chan<- *message) {
	go s.ListenEvents()
	go s.ListenFSM(register, unregister, directMsg)
}

// Decode FSM input/signals
func (s *mWorker) ListenEvents() {

	for event := range s.fsm.newEvent {
		switch s.fsm.curState {
		case disconnected:
			switch event {
			case accept:
				s.fsm.nxtState <- listening
			case failed:
				s.fsm.nxtState <- disconnected
			case open:
				s.fsm.nxtState <- connecting
			case close:
				s.fsm.nxtState <- disconnected
			case wait:
				s.fsm.nxtState <- s.fsm.curState
			default:
				fmt.Println("ERROR - Bad event:", event)
			}
		case connecting:
			switch event {
			case accept:
				s.fsm.nxtState <- listening
			case failed:
				s.fsm.nxtState <- disconnected
			case open:
				s.fsm.nxtState <- connecting
			case close:
				s.fsm.nxtState <- disconnected
			case wait:
				s.fsm.nxtState <- s.fsm.curState
			default:
				fmt.Println("ERROR - Bad event:", event)
			}
		case listening:
			switch event {
			case accept:
				s.fsm.nxtState <- listening
			case failed:
				s.fsm.nxtState <- disconnected
			case open:
				s.fsm.nxtState <- connecting
			case close:
				s.fsm.nxtState <- disconnected
			case wait:
				fmt.Println("received wait - staying in same state!")
				s.fsm.nxtState <- s.fsm.curState

			default:
				fmt.Println("ERROR - Bad event:", event)
			}
		default:
			fmt.Println("ListenEvents : Bad State : ", s.fsm.newEvent)
		}
	}
}

// Receives new state and trigger fuctionality
func (s *mWorker) ListenFSM(register, unregister chan<- *mWorker, directMsg chan<- *message) {

	for nxtState := range s.fsm.nxtState {
		fmt.Println("DEBUG - received next state : ", nxtState)

		switch nxtState {

		case disconnected:
			fmt.Println("received new state : ", nxtState)
			unregister <- s
			s.fsm.curState = nxtState

		case connecting:
			fmt.Println("received new state : ", nxtState)
			register <- s
			s.send([]byte(s.PrepareMsg(NewPayload("RD", cmd)).Compile()))
			fmt.Println("New connection added")
			s.fsm.curState = nxtState
			s.fsm.newEvent <- accept

		case listening:
			fmt.Println("received new state : ", nxtState)
			fmt.Println("waiting for requests")
			s.fsm.curState = nxtState

		default:
			fmt.Println("ListenFSM : Bad State")

		}
	}
}
