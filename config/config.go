package config

import (
	"encoding/json"
	"fmt"
	"go-imdg/comms"
	"log"
	"os"
	"strings"
	"sync"
)

type Node struct {
	Logger     *log.Logger
	nodeType   string
	lPort      string
	name       string
	masterConn string
}

func (n *Node) Start() {
	if strings.Compare(n.nodeType, "master") == 0 {
		fmt.Println("Starting new master...")
		master := comms.NewMaster("localhost", n.lPort)
		master.Start()

	} else if strings.Compare(n.nodeType, "worker") == 0 {
		fmt.Println("Starting new worker... ")
		worker := comms.NewWorker(n.masterConn, n.name, "localhost", n.lPort)

		fmt.Println("new worker:", worker)

		worker.Start()

		var message string
		for {
			fmt.Print("Enter message:")
			fmt.Scan(&message)
			// NewPayload(message, payloadType.def)

			worker.PrepareMsg("3:" + message)
			worker.SendMsg()
		}
	} else {

	}
}

type nodeCfg struct {
	logFile    *os.File
	nodeType   string
	port       string
	name       string
	masterConn string
}

func New(ncfg string) *Node {
	inputConfig, err := LoadConfig(ncfg)
	if err != nil {
		return nil
	}
	return GetInstance(inputConfig)

	// return GetInstance(func(s string) nodeCfg {
	// 	inputConfig, err := LoadConfig(s)
	// 	if err != nil {
	// 		return nil
	// 	}
	// 	return inputConfig
	// }(ncfg))

}

var instance *Node
var once sync.Once

func GetInstance(ncfg nodeCfg) *Node {
	once.Do(func() {
		instance = &Node{
			Logger:     log.New(ncfg.logFile, "", log.Ldate|log.Ltime|log.Lshortfile),
			nodeType:   ncfg.nodeType,
			lPort:      ncfg.port,
			name:       ncfg.name,
			masterConn: ncfg.masterConn,
		}
	})
	return instance
}

func LoadConfig(cfgFileName string) (nodeCfg, error) {
	var tempCfg nodeCfg
	cfgData, err := os.ReadFile(cfgFileName)
	if err != nil {
		return nodeCfg{}, err
	}

	var jsonMap map[string]interface{}
	json.Unmarshal(cfgData, &jsonMap)

	ln, ok := jsonMap["log_name"].(string)
	if !ok {
		log.Fatal("log_name is not a string", ok)
	}

	ld, ok := jsonMap["log_dir"].(string)
	if !ok {
		log.Fatal("log_file is not a string", ok)
	}

	err = os.MkdirAll(ld, 0755)
	if err != nil {
		log.Fatal("Failed to create log directory : ", err)
	}
	logFile, err := os.Create(ld + "/" + ln + ".log")
	if err != nil {
		log.Fatal("Failed to create log file", err)
	}
	tempCfg.logFile = logFile

	nt, ok := jsonMap["node_type"].(string)
	if !ok {
		log.Fatal("node_type is not a string", ok)
	}
	tempCfg.nodeType = nt

	lp, ok := jsonMap["listening_port"].(string)
	if !ok {
		log.Fatal("listening_port is not a string", ok)
	}
	tempCfg.port = lp

	nn, ok := jsonMap["node_name"].(string)
	if !ok {
		log.Fatal("log_file is not a string", ok)
	}
	tempCfg.name = nn

	if strings.Compare(nt, "worker") == 0 {
		mc, ok := jsonMap["master_conn"].(string)
		if !ok {
			log.Fatal("log_file is not a string", ok)
		}
		tempCfg.masterConn = mc
	} else {
		tempCfg.masterConn = "NA"
	}

	fmt.Println("CFG: ", tempCfg)

	return tempCfg, nil
}
