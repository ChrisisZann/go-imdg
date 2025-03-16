package comms

import (
	"log"
)

func ListenConnEvents(l *log.Logger, curState *connState, nxtState chan<- connState, newEvent <-chan VarFSM) {

	for event := range newEvent {
		l.Printf("DEBUG - curState:%s <- %s\n", curState, event)

		switch *curState {
		case disconnected:
			switch event {
			case open:
				nxtState <- connecting
			case wait:
				nxtState <- *curState
			default:
				l.Println("ERROR - Bad event:", event)
				nxtState <- invalid_state
			}
		case connecting:
			switch event {
			case accept:
				nxtState <- listening
			case fatal:
				nxtState <- disconnected
			case wait:
				nxtState <- *curState
			default:
				l.Println("ERROR - Bad event:", event)
				nxtState <- invalid_state
			}
		case listening:
			switch event {
			case fatal:
				nxtState <- disconnected
			case close:
				nxtState <- disconnected
			case wait:
				log.Println("received wait - staying in same state!")
				nxtState <- *curState

			default:
				l.Println("ERROR - Bad event:", event)
				nxtState <- invalid_state
			}
		default:
			l.Println("ListenEvents : Bad State : ", curState)
			nxtState <- invalid_state
		}
	}
}
