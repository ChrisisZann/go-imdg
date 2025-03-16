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

		switch nxtState {

		case disconnected:

			unregister <- s
			s.fsm.curState = nxtState

		case connecting:

			register <- s
			s.send([]byte(s.PrepareMsg(NewPayload("RD", cmd)).Compile()))
			s.logger.Println("New connection added")
			s.fsm.curState = nxtState
			s.fsm.NewEvent <- accept
			// s.logger.Println(s.fsm.curState, "...")

		case listening:
			// s.logger.Println(nxtState, "...")
			s.logger.Println("worker waiting for requests...")

			s.fsm.curState = nxtState

		default:
			s.logger.Println("ListenFSM : Bad State")

		}

		s.logger.Println("Finished processing state")
		s.logger.Println("Current:", s.fsm.curState)

	}
}
