package comms

func (s *mWorker) Start(register, unregister chan<- *mWorker, directMsg chan<- *message) {
	go s.ListenEvents()
	go s.ListenFSM(register, unregister, directMsg)

}

// Decode FSM input/signals
func (s *mWorker) ListenEvents() {
	ListenConnEvents(s.logger, &s.fsm.curState, s.fsm.nxtState, s.fsm.NewEvent)
}

// Receives new state and trigger fuctionality
func (s *mWorker) ListenFSM(register, unregister chan<- *mWorker, directMsg chan<- *message) {

	for nxtState := range s.fsm.nxtState {
		s.logger.Println("Current:", s.fsm.curState)
		s.logger.Println("Received nxtState:", nxtState)
		newState := nxtState

		switch nxtState {

		case disconnected:

			unregister <- s

		case connecting:

			register <- s
			s.send([]byte(s.PrepareMsg(
				NewPayload(ParseVarFSM("accept").String(), cmd)).Compile()))

			s.logger.Println("New connection added")
			s.fsm.NewEvent <- accept

		case listening:

			s.logger.Println("worker waiting for requests...")
			// s.fsm.NewEvent <- wait

		default:

			s.logger.Println("ListenFSM : Bad State")
			newState = invalid_state

		}
		s.fsm.curState = newState
		s.logger.Println("Finished processing state")
		s.logger.Println("Current:", s.fsm.curState)

	}
}
