package comms

import "fmt"

func (w *Worker) ListenEvents() {

	for {
		select {
		case event := <-w.fsm.newEvent:
			fmt.Println("DEBUG - received new event : ", event)

			switch w.fsm.curState {
			case notConnected:
				switch event {
				case accept:
					w.fsm.nxtState <- connected
				case failed:
					w.fsm.nxtState <- notConnected
				case open:
					w.fsm.nxtState <- validateNewConn
				case close:
					w.fsm.nxtState <- notConnected
				case wait:
					w.fsm.nxtState <- w.fsm.curState
				default:
					fmt.Println("ERROR - Bad event!")
				}
			case validateNewConn:
				switch event {
				case accept:
					w.fsm.nxtState <- connected
				case failed:
					w.fsm.nxtState <- notConnected
				case open:
					w.fsm.nxtState <- validateNewConn
				case close:
					w.fsm.nxtState <- notConnected
				case wait:
					w.fsm.nxtState <- w.fsm.curState
				default:
					fmt.Println("ERROR - Bad event!")
				}
			case connected:
				switch event {
				case accept:
					w.fsm.nxtState <- connected
				case failed:
					w.fsm.nxtState <- notConnected
				case open:
					w.fsm.nxtState <- validateNewConn
				case close:
					w.fsm.nxtState <- notConnected
				case wait:
					fmt.Println("received wait - staying in same state!")
					w.fsm.nxtState <- w.fsm.curState

				default:
					fmt.Println("ERROR - Bad event!")
				}
			default:
				fmt.Println("ListenEvents : Bad State")
			}
		}
	}
}
