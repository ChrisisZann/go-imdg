package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var instance *NodeCfg
var once sync.Once

type NodeCfg struct {
	LogFile    *os.File
	Logger     *log.Logger
	NodeType   string
	Hostname   string
	LPort      string
	Name       string
	MasterConn string
}

func New(ncfg string) *NodeCfg {
	inputConfig, err := LoadConfig(ncfg)
	if err != nil {
		return nil
	}
	return GetInstance(inputConfig)
}

func GetInstance(ncfg NodeCfg) *NodeCfg {
	once.Do(func() {
		instance = &ncfg
	})
	return instance
}

func readString(input interface{}) string {
	nt, ok := input.(string)
	if !ok {
		log.Fatal("node_type is not a string", ok)
		return ""
	}
	return nt
}

func LoadConfig(cfgFileName string) (NodeCfg, error) {
	var tempCfg NodeCfg
	cfgData, err := os.ReadFile(cfgFileName)
	if err != nil {
		return NodeCfg{}, err
	}

	var jsonMap map[string]interface{}
	json.Unmarshal(cfgData, &jsonMap)

	ln := readString(jsonMap["log_name"])
	ld := readString(jsonMap["log_dir"])

	err = os.MkdirAll(ld, 0755)
	if err != nil {
		log.Fatal("Failed to create log directory : ", err)
	}
	logFile, err := os.Create(ld + "/" + ln + ".log")
	if err != nil {
		log.Fatal("Failed to create log file", err)
	}
	tempCfg.LogFile = logFile

	tempCfg.NodeType = readString(jsonMap["node_type"])
	tempCfg.Hostname = readString(jsonMap["hostname"])
	tempCfg.LPort = readString(jsonMap["listening_port"])
	tempCfg.Name = readString(jsonMap["node_name"])
	tempCfg.Logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)

	if strings.Compare(tempCfg.NodeType, "worker") == 0 {
		mc, ok := jsonMap["master_conn"].(string)
		if !ok {
			log.Fatal("log_file is not a string", ok)
		}
		tempCfg.MasterConn = mc
	} else {
		tempCfg.MasterConn = "NA"
	}

	fmt.Println("CFG: ", tempCfg)

	return tempCfg, nil
}
