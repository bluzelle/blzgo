package main

import (
	"github.com/apex/log"
	"github.com/vbstreetz/blzgo/examples/util"
)

func main() {
	util.SetupLogging()
	util.LoadEnv()

	ctx, err := util.NewClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Infof("getting bluzelle version")

	if v, err := ctx.Version(); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("%s", v)
	}
}
