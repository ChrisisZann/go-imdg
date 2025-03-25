package main

import (
	"flag"
	"fmt"
	"go-imdg/comms"
	"go-imdg/config"
	"go-imdg/data"
	"go-imdg/node"
	"log"
	"strings"
)

func main() {
	cfg := flag.String("config", "bad.json", "provide config file")

	flag.Parse()

	fmt.Println(*cfg)

	n := config.New(*cfg)

	if strings.Compare(n.NodeType, "master") == 0 {
		m := node.NewMaster(*n)
		m.Start()

	} else if strings.Compare(n.NodeType, "slave") == 0 {
		s := node.NewSlave(*n)

		// Connection to master
		s.NewMasterConnection("localhost", "3333")
		fmt.Println("TESTING2:", s.MasterConnection)

		//
		go s.ReceiveHandler()
		s.StartMasterConnectionLoop(s.Receiver)

		go s.Listen()

		var Message string
		for {
			fmt.Print("Enter Message:")
			fmt.Scan(&Message)

			// s.MasterConnection.SendMsg(s.PrepareMsg(comms.NewPayload(Message, comms.PayloadType(0))))
			p, err := comms.NewPayload(Message, "cmd")
			if err != nil {
				log.Fatalln("Failed to create new payload")
			}
			s.MasterConnection.SendPayload(p)
		}

		// s.MasterConnection.Send <- s.PrepareMsg(comms.NewPayload("Hello", comms.PayloadType(0)))
	} else if strings.Compare(n.NodeType, "tester") == 0 {
		var s data.MemPage
		s.Init()

		s.Save([]byte("testing"))
		s.Save([]byte("testing1"))
		s.Save([]byte("testing2"))
		s.Save([]byte("testing3"))
		s.Save([]byte("testing4"))
		s.Save([]byte("testing5"))
		s.Save([]byte("testing6"))
		s.Save([]byte("testing7"))
		s.Save([]byte("testing8"))
		s.Save([]byte("testing9"))
		s.Save([]byte("testing10"))
		s.Save([]byte("testing11"))
		s.Save([]byte("testing12"))
		s.Save([]byte("testing13"))
		s.Save([]byte("testing14"))

		for i := 0; i < len(s.Page); i++ {
			n.Logger.Println("Reading:", string(s.Read(i)))
		}

		n.Logger.Println("Page:", s)

	} else {
		log.Panicln("Bad node type")
	}
}
