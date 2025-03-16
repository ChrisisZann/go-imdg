package comms

func (cb *CommsBox) ListenConnEvents() {

	ListenConnEvents(cb.logger, &cb.ConnControl.curState, cb.ConnControl.nxtState, cb.ConnControl.NewEvent)
}

func (cb *CommsBox) FsmConnProcessor() {

	for nxtState := range cb.ConnControl.nxtState {

		cb.logger.Println("current state 	  : ", cb.ConnControl.curState)
		cb.logger.Println("received new state : ", nxtState)
		newState := nxtState

		switch nxtState {
		case disconnected:
			// Disconnect from master
		case connecting:
			// Connect to master

			// Send registration message
			cb.Send <- cb.PrepareMsg(
				NewPayload(VarFSM(2).String(), PayloadType(0)))

		case listening:
			// Wait for requests
		default:
			cb.logger.Println("Invalid state")
			newState = connState(fatal)
		}

		cb.ConnControl.curState = newState
	}
}
