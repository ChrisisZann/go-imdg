package comms

import "fmt"

func (s *mWorker) Start(register, unregister chan<- *mWorker, directMsg chan<- *message) {
	go s.ListenEvents()
	go s.ListenFSM(register, unregister, directMsg)
}

// Decode FSM input/signals
func (s *mWorker) ListenEvents() {

	for {
		select {
		case event := <-s.fsm.newEvent:
			fmt.Println("DEBUG - received new event : ", event)

			switch s.fsm.curState {
			case notConnected:
				switch event {
				case accept:
					s.fsm.nxtState <- connected
				case failed:
					s.fsm.nxtState <- notConnected
				case open:
					s.fsm.nxtState <- validateNewConn
				case close:
					s.fsm.nxtState <- notConnected
				case wait:
					s.fsm.nxtState <- s.fsm.curState
				default:
					fmt.Println("ERROR - Bad event!")
				}
			case validateNewConn:
				switch event {
				case accept:
					s.fsm.nxtState <- connected
				case failed:
					s.fsm.nxtState <- notConnected
				case open:
					s.fsm.nxtState <- validateNewConn
				case close:
					s.fsm.nxtState <- notConnected
				case wait:
					s.fsm.nxtState <- s.fsm.curState
				default:
					fmt.Println("ERROR - Bad event!")
				}
			case connected:
				switch event {
				case accept:
					s.fsm.nxtState <- connected
				case failed:
					s.fsm.nxtState <- notConnected
				case open:
					s.fsm.nxtState <- validateNewConn
				case close:
					s.fsm.nxtState <- notConnected
				case wait:
					fmt.Println("received wait - staying in same state!")
					s.fsm.nxtState <- s.fsm.curState

				default:
					fmt.Println("ERROR - Bad event!")
				}
			default:
				fmt.Println("ListenEvents : Bad State")
			}
		}
	}
}

// Receives new state and trigger fuctionality
func (s *mWorker) ListenFSM(register, unregister chan<- *mWorker, directMsg chan<- *message) {

	for {
		select {
		case nxtState := <-s.fsm.nxtState:
			fmt.Println("DEBUG - received next state : ", nxtState)

			switch nxtState {

			case notConnected:
				fmt.Println("received new state : ", nxtState)
				unregister <- s
				s.fsm.curState = nxtState

			case validateNewConn:
				fmt.Println("received new state : ", nxtState)
				register <- s
				s.send([]byte(s.PrepareMsg(NewPayload("RD", cmd)).Compile()))
				fmt.Println("New connection added")
				s.fsm.curState = nxtState
				s.fsm.newEvent <- accept

			case connected:
				fmt.Println("received new state : ", nxtState)
				fmt.Println("waiting for requests")
				s.fsm.curState = nxtState

			default:
				fmt.Println("ListenFSM : Bad State")

			}
		}
	}
}
