package node

import (
	"go-imdg/config"
	"log"
	"sync"
)

var instance *Node
var once sync.Once

type Node struct {
	Logger     *log.Logger
	NodeType   string
	LPort      string
	Name       string
	MasterConn string
}

func New(ncfg string) *Node {
	inputConfig, err := config.LoadConfig(ncfg)
	if err != nil {
		return nil
	}
	return GetInstance(inputConfig)
}

func GetInstance(ncfg config.NodeCfg) *Node {
	once.Do(func() {
		instance = &Node{
			Logger:     log.New(ncfg.LogFile, "", log.Ldate|log.Ltime|log.Lshortfile),
			NodeType:   ncfg.NodeType,
			LPort:      ncfg.Port,
			Name:       ncfg.Name,
			MasterConn: ncfg.MasterConn,
		}
	})
	return instance
}
