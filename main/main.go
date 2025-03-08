package main

import (
	"flag"
	"fmt"
	"go-imdg/comms"
	"go-imdg/node"
	"log"

	"strings"
)

var cfgFile = flag.String("c", "config.json", "configuration file")

func main() {

	fmt.Println("Node starting...")

	flag.Parse()

	testCfg := node.New(*cfgFile)

	if strings.Compare(testCfg.NodeType, "master") == 0 {
		fmt.Println("Starting new master...")
		master := comms.NewMaster("localhost", testCfg.LPort)
		master.Start()

	} else if strings.Compare(testCfg.NodeType, "worker") == 0 {
		fmt.Println("Starting new worker... ")
		worker := comms.NewWorker(testCfg.MasterConn, testCfg.Name, "localhost", testCfg.LPort)

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
		log.Panicln("Bad node type")
	}
}
