package main

import (
	"flag"
	"fmt"
	"go-imdg/comms"
)

var cfgFile = flag.String("c", "config.json", "configuration file")

func main() {

	fmt.Println("Node starting...")

	var nodeType = flag.String("type", "Slave", "select node type: master or Slave")
	flag.Parse()

	switch *nodeType {
	case "master":
		fmt.Println("Starting new master...")
		master := comms.NewMaster("localhost", "3333")
		master.Start()

	case "Slave":
		fmt.Println("Starting new slave...")
		// slave := comms.NewSlave("localhost:3333", "s1", "localhost", "3334")
		slave := comms.NewSlave("localhost:3333", "s2", "localhost", "3335")

		slave.Start()

		var message string
		for {
			fmt.Print("Enter message:")
			fmt.Scan(&message)
			// NewPayload(message, payloadType.def)

			slave.PrepareMsg("3:" + message)
			slave.SendMsg()
		}
	}
}
