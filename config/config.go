package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type NodeCfg struct {
	LogFile    *os.File
	NodeType   string
	Port       string
	Name       string
	MasterConn string
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
	tempCfg.Port = readString(jsonMap["listening_port"])
	tempCfg.Name = readString(jsonMap["node_name"])

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
