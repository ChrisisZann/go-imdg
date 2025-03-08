package main

import (
	"flag"
	"fmt"
	"go-imdg/config"
)

var cfgFile = flag.String("c", "config.json", "configuration file")

func main() {

	fmt.Println("Node starting...")

	flag.Parse()

	testCfg := config.New(*cfgFile)
	testCfg.Start()
}
