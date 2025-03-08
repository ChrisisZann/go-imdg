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

func LoadConfig(cfgFileName string) (NodeCfg, error) {
	var tempCfg NodeCfg
	cfgData, err := os.ReadFile(cfgFileName)
	if err != nil {
		return NodeCfg{}, err
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
	tempCfg.LogFile = logFile

	nt, ok := jsonMap["node_type"].(string)
	if !ok {
		log.Fatal("node_type is not a string", ok)
	}
	tempCfg.NodeType = nt

	lp, ok := jsonMap["listening_port"].(string)
	if !ok {
		log.Fatal("listening_port is not a string", ok)
	}
	tempCfg.Port = lp

	nn, ok := jsonMap["node_name"].(string)
	if !ok {
		log.Fatal("log_file is not a string", ok)
	}
	tempCfg.Name = nn

	if strings.Compare(nt, "worker") == 0 {
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
