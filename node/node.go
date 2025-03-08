package node

import (
	"go-imdg/config"
	"log"
	"sync"
)

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

	// return GetInstance(func(s string) NodeCfg {
	// 	inputConfig, err := LoadConfig(s)
	// 	if err != nil {
	// 		return nil
	// 	}
	// 	return inputConfig
	// }(ncfg))

}

var instance *Node
var once sync.Once

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
