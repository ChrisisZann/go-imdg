package node

import (
	"context"
	"go-imdg/comms"
	"go-imdg/config"
	"os"
	"sync"
	"time"
)

type NodeBase struct {
	id int

	config.Node
	comms.NetworkReader
	comms.NetworkWriter

	Receiver chan *comms.Message

	ctx    context.Context
	cancel func()
	wg     sync.WaitGroup
}

func NewNodeBase(cfg config.Node) *NodeBase {

	newAddr, err := comms.NewNodeAddr("tcp", cfg.Hostname+":"+cfg.LPort)
	if err != nil {
		cfg.Logger.Fatal("Failed to create NodeAddr for node:", err)
	}

	m, cf := context.WithCancel(context.Background())

	newNetR := comms.NewNetworkReader(
		newAddr,
		cfg.Logger,
		10*time.Second,
	)
	if newNetR == nil {
		cfg.Logger.Fatal("error - failed to create new NewNetworkReader")
	}

	return &NodeBase{
		id:            os.Getpid(),
		Node:          cfg,
		NetworkReader: *newNetR,
		Receiver:      make(chan *comms.Message, 10),
		ctx:           m,
		cancel:        cf,
	}
}

func (n *NodeBase) Start() {
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		n.NetworkReader.Listen(n.ctx, n.Receiver)
	}()

	n.wg.Wait()
}
