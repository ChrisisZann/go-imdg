package comms

import "fmt"

func (w *Worker) ListenEvents() {

	for event := range w.fsm.newEvent {
		fmt.Println("DEBUG - received new event : ", event)

		switch w.fsm.curState {
		case disconnected:
			switch event {
			case accept:
				w.fsm.nxtState <- listening
			case failed:
				w.fsm.nxtState <- disconnected
			case open:
				w.fsm.nxtState <- connecting
			case close:
				w.fsm.nxtState <- disconnected
			case wait:
				w.fsm.nxtState <- w.fsm.curState
			default:
				fmt.Println("ERROR - Bad event:", event)
			}
		case connecting:
			switch event {
			case accept:
				w.fsm.nxtState <- listening
			case failed:
				w.fsm.nxtState <- disconnected
			case open:
				w.fsm.nxtState <- connecting
			case close:
				w.fsm.nxtState <- disconnected
			case wait:
				w.fsm.nxtState <- w.fsm.curState
			default:
				fmt.Println("ERROR - Bad event:", event)
			}
		case listening:
			switch event {
			case accept:
				w.fsm.nxtState <- listening
			case failed:
				w.fsm.nxtState <- disconnected
			case open:
				w.fsm.nxtState <- connecting
			case close:
				w.fsm.nxtState <- disconnected
			case wait:
				fmt.Println("received wait - staying in same state!")
				w.fsm.nxtState <- w.fsm.curState

			default:
				fmt.Println("ERROR - Bad event:", event)
			}
		default:
			fmt.Println("ListenEvents : Bad State : ", w.fsm.curState)
		}
	}
}
