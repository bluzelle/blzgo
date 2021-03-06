package main

import (
	"github.com/apex/log"
	bluzelle "github.com/bluzelle/blzgo"
)

func main() {
	bluzelle.SetupLogging()
	bluzelle.LoadEnv()

	ctx, err := bluzelle.NewTestClient()
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
