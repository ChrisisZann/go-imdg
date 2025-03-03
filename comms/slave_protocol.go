package comms

// import (
// 	"fmt"
// 	"log"
// )

// type protocolMSG int

// const (
// 	portRequest protocolMSG = iota
// 	dataRequest
// 	finACK
// 	done
// 	status
// 	success
// 	failure
// 	new_message
// )

// type slaveState int

// const (
// 	initialize slaveState = iota
// 	startListening
// 	runLoop
// 	waiting
// 	fatal
// )

// func (s *Slave) protocolFSM() {
// 	// currentState := s.curState
// 	for {
// 		select {
// 		case event := <-s.stateMSG:

// 			switch s.curState {

// 			case initialize:
// 				switch event {
// 				case success:
// 					s.nxtState = startListening
// 				case failure:
// 					s.nxtState = fatal
// 				}

// 			case startListening:
// 				switch event {
// 				case success:
// 					s.nxtState = runLoop
// 				case failure:
// 					s.nxtState = fatal
// 				}

// 			case runLoop:
// 				switch event {
// 				case success:
// 					s.nxtState = waiting
// 				case failure:
// 					s.nxtState = fatal
// 				case new_message:
// 				}

// 			case fatal:
// 				log.Fatal("Fatal : Process failed")

// 			default:
// 				log.Println("Unknown state")
// 			}
// 		}
// 		s.curState = s.nxtState
// 		fmt.Println("New state", s.curState)
// 	}
// }
