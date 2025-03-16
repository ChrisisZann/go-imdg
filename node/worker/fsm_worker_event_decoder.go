package worker

import "fmt"

func (w *Worker) ListenInternalEvents() {
	for event := range w.internalFSM.newEvent {
		switch w.internalFSM.curState {
		case stopped:
			switch event {
			case start:
				w.internalFSM.nxtState <- startup
			case stop:
				w.internalFSM.nxtState <- shutdown
			default:
				fmt.Println("ERROR - Bad event:", event)
			}
		case startup:
			switch event {
			case success:
				w.internalFSM.nxtState <- waiting
			case failure:
				w.internalFSM.nxtState <- fatal
			default:
				fmt.Println("ERROR - Bad event:", event)
			}
		case waiting:
			switch event {
			case request:
				w.internalFSM.nxtState <- process
			case failure:
				w.internalFSM.nxtState <- fatal
			case wait:
				w.internalFSM.nxtState <- w.internalFSM.curState
			default:
				fmt.Println("ERROR - Bad event:", event)
			}
		case process:
			switch event {
			case success:
				w.internalFSM.nxtState <- waiting
			case failure:
				w.internalFSM.nxtState <- fatal
			default:
				fmt.Println("ERROR - Bad event:", event)
			}
		case shutdown:
			switch event {
			case success:
				w.internalFSM.nxtState <- w.internalFSM.curState
			case failure:
				w.internalFSM.nxtState <- fatal
			default:
				fmt.Println("ERROR - Bad event:", event)
			}
		}
	}
}
