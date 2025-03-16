package main

import (
	"flag"
	"fmt"
	"go-imdg/config"
	"go-imdg/node/comms"
	"go-imdg/node/data"
	"go-imdg/node/worker"
	"log"
	"strings"
)

var cfgFile = flag.String("c", "config.json", "configuration file")

func main() {

	fmt.Println("Node starting...")

	flag.Parse()

	testCfg := config.New(*cfgFile)

	if strings.Compare(testCfg.NodeType, "master") == 0 {

		testCfg.Logger.Println("Starting new master...")

		master := comms.NewMaster(testCfg)

		master.Start()

	} else if strings.Compare(testCfg.NodeType, "worker") == 0 {
		testCfg.Logger.Println("Starting new worker... ")

		// worker := worker.NewWorker(testCfg.MasterConn, testCfg.Name, "localhost", testCfg.LPort)

		worker := worker.NewWorker(testCfg)
		testCfg.Logger.Println("new worker:", worker)

		worker.Start()

		var message string
		for {
			fmt.Print("Enter message:")
			fmt.Scan(&message)

			worker.SendMsg(worker.PrepareMsg(comms.NewPayload(message, comms.PayloadType(0))))
		}
	} else if strings.Compare(testCfg.NodeType, "tester") == 0 {
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
			testCfg.Logger.Println("Reading:", string(s.Read(i)))
		}

		testCfg.Logger.Println("Page:", s)

	} else {
		log.Panicln("Bad node type")
	}
}
