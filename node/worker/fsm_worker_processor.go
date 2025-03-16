package worker

import (
	"go-imdg/node/comms"
)

func (w *Worker) fsmInternalProcessor() {

	for nxtState := range w.internalFSM.nxtState {
		w.Logger.Println("Current:", w.internalFSM.curState)
		w.Logger.Println("Received nxtState:", nxtState)
		newState := nxtState

		switch nxtState {

		case startup:
			// Connect to master
			// Send registration message
			w.Logger.Println("sending open message")
			w.ConnControl.NewEvent <- comms.ParseVarFSM("open")

		case waiting:
			// Wait for requests
		case process:
			// Process requests
		case shutdown:
			// Disconnect from master
		case fatal:
			// Log error
		default:
			w.Logger.Println("Invalid state")
			newState = fatal
		}

		w.internalFSM.curState = newState
	}
}
