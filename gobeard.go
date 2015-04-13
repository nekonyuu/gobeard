package main

import (
	"flag"

	"github.com/apognu/gobeard/api"
	"github.com/apognu/gobeard/client"
	"github.com/apognu/gobeard/util"
)

var (
	daemon     *bool   = flag.Bool("d", false, "launch in daemon mode")
	configPath *string = flag.String("config", "/etc/gobeard.yaml", "configuration file to use")
)

func main() {
	flag.Parse()
	util.SetConfig(*configPath)

	if *daemon {
		go api.NewWebUi()
		NewDaemon()
	} else {
		client.NewClient()
	}
}
