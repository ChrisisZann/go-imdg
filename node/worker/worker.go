package worker

import (
	"go-imdg/config"
	"go-imdg/node/comms"
)

type Worker struct {
	comms.CommsBox

	// embed Node struct
	config.NodeCfg

	internalFSM workerControl
}

func NewWorker(cfg *config.NodeCfg) *Worker {
	return &Worker{
		CommsBox:    *comms.NewCommsBox(cfg.MasterConn, cfg.Name, cfg.Hostname, cfg.LPort, cfg.Logger),
		NodeCfg:     *cfg,
		internalFSM: *NewWorkerControl(),
	}
}

func (s *Worker) Start() {
	// Tx Rx
	go s.RunComms()

	// Listen for incoming connections
	go s.Listen()

	// Listen for internal events
	go s.ListenInternalEvents()
	go s.fsmInternalProcessor()

	// Listen for external events
	go s.ListenConnEvents()
	go s.FsmConnProcessor()

	// Start the worker
	s.internalFSM.nxtState <- startup
}
