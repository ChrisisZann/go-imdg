package main

import (
	"flag"
	"fmt"
	"go-imdg/comms"
)

var cfgFile = flag.String("c", "config.json", "configuration file")

func main() {

	fmt.Println("Node starting...")

	var nodeType = flag.String("type", "slave", "select node type: master or slave")
	flag.Parse()

	switch *nodeType {
	case "master":
		master := comms.NewMaster("localhost", "3333")

		go master.Run()
		master.Listen()

	case "slave":
		s := comms.NewSlave("localhost:3333", "s2", "localhost", "3335")

		go s.Run()
		go s.Listen()

		var message string
		for {
			fmt.Print("Enter message:")
			fmt.Scan(&message)
			s.Send <- []byte(message)
		}
	}
}
