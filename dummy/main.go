package main

import (
	"flag"
	dummy "ping-ping/dummy/src"
)

// this file defines the main routine of Dummy, which takes input arguments from the command line

func main() {
	configFile := flag.String("config", "dedis-config.yaml", "configuration file")
	name := flag.Int64("name", 1, "name of the replica")
	debugOn := flag.Bool("debugOn", false, "true / false")
	debugLevel := flag.Int("debugLevel", 1, "debug level")

	flag.Parse()

	replicas, err := dummy.ReadYAML(*configFile, int(*name))
	if err != nil {
		panic(err.Error())
	}

	proxyInstance := dummy.NewProxy(*name, replicas, *debugOn, *debugLevel)

	proxyInstance.NetworkInit()
	proxyInstance.Run()

}
