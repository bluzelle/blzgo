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

	log.Infof("getting number of keys...")

	if v, err := ctx.Count(); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("number of keys (%d)", v)
	}
}
